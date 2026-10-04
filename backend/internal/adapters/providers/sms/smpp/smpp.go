// Package smpp is a real SMPP 3.4 client: bind_transceiver, submit_sm with
// UCS2 for Cyrillic/Turkmen text, long-message concatenation via UDH,
// deliver_sm delivery receipts, enquire_link keep-alive and automatic
// rebinding. One Session is kept per provider for the worker's lifetime.
package smpp

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/linxGnu/gosmpp"
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/pdu"
	"github.com/rs/zerolog"

	"github.com/Esca6585dev/habarchy/backend/internal/ports"
)

// Config is the decrypted credentials JSON of an smpp provider. It mirrors
// providers.SMPPConfig.
type Config struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	SystemID       string `json:"system_id"`
	Password       string `json:"password"`
	SystemType     string `json:"system_type"`
	SourceAddr     string `json:"source_addr"`
	SourceTON      int    `json:"source_ton"`
	SourceNPI      int    `json:"source_npi"`
	DestTON        int    `json:"dest_ton"`
	DestNPI        int    `json:"dest_npi"`
	EnquireLinkSec int    `json:"enquire_link_sec"`
	UseTLS         bool   `json:"use_tls"`
	RequestDLR     bool   `json:"request_dlr"`
	// SubmitTimeoutSec bounds waiting for submit_sm_resp (default 20).
	SubmitTimeoutSec int `json:"submit_timeout_sec"`
}

// Receipt is a parsed deliver_sm delivery report.
type Receipt struct {
	MessageID string
	Status    string // DELIVRD, UNDELIV, EXPIRED, REJECTD, ACCEPTD, ...
	Delivered bool
	Final     bool
	Text      string
}

// ReceiptHandler receives delivery reports.
type ReceiptHandler func(ctx context.Context, r Receipt)

// Session is a bound SMPP transceiver implementing ports.SMSProvider.
type Session struct {
	cfg     Config
	log     zerolog.Logger
	onDLR   ReceiptHandler
	session *gosmpp.Session

	mu      sync.Mutex
	pending map[int32]chan *pdu.SubmitSMResp
	bound   bool
	closed  bool
}

// Dial binds to the SMSC. It returns once the bind succeeded or failed.
func Dial(ctx context.Context, cfg Config, log zerolog.Logger, onDLR ReceiptHandler) (*Session, error) {
	if cfg.Host == "" || cfg.Port == 0 || cfg.SystemID == "" {
		return nil, errors.New("smpp: host, port and system_id are required")
	}
	if cfg.EnquireLinkSec <= 0 {
		cfg.EnquireLinkSec = 60
	}
	if cfg.SubmitTimeoutSec <= 0 {
		cfg.SubmitTimeoutSec = 20
	}
	s := &Session{cfg: cfg, log: log, onDLR: onDLR, pending: map[int32]chan *pdu.SubmitSMResp{}}

	auth := gosmpp.Auth{SMSC: net.JoinHostPort(cfg.Host, fmt.Sprint(cfg.Port)), SystemID: cfg.SystemID, Password: cfg.Password, SystemType: cfg.SystemType}
	dialer := gosmpp.NonTLSDialer
	if cfg.UseTLS {
		dialer = func(addr string) (net.Conn, error) {
			d := &tls.Dialer{NetDialer: &net.Dialer{Timeout: 10 * time.Second}, Config: &tls.Config{MinVersion: tls.VersionTLS12}}
			return d.DialContext(ctx, "tcp", addr)
		}
	}
	settings := gosmpp.Settings{
		ReadTimeout:  time.Duration(cfg.EnquireLinkSec)*time.Second + 30*time.Second,
		WriteTimeout: 10 * time.Second,
		EnquireLink:  time.Duration(cfg.EnquireLinkSec) * time.Second,
		OnPDU:        s.onPDU,
		OnSubmitError: func(p pdu.PDU, err error) {
			log.Error().Err(err).Int32("seq", p.GetSequenceNumber()).Msg("smpp submit error")
			s.resolve(p.GetSequenceNumber(), nil)
		},
		OnReceivingError: func(err error) { log.Warn().Err(err).Msg("smpp receive error") },
		OnRebindingError: func(err error) { log.Error().Err(err).Msg("smpp rebind failed") },
		OnRebind:         func() { log.Info().Msg("smpp rebound") },
		OnClosed: func(state gosmpp.State) {
			log.Warn().Str("state", state.String()).Msg("smpp connection closed")
			s.mu.Lock()
			s.bound = state == gosmpp.ExplicitClosing || !s.closed && false
			s.mu.Unlock()
		},
	}
	sess, err := gosmpp.NewSession(gosmpp.TRXConnector(dialer, auth), settings, 5*time.Second)
	if err != nil {
		return nil, &ports.ProviderError{Code: "smpp_bind", Message: err.Error(), Retryable: true}
	}
	s.session = sess
	s.mu.Lock()
	s.bound = true
	s.mu.Unlock()
	return s, nil
}

// Close unbinds.
func (s *Session) Close() error {
	s.mu.Lock()
	s.closed = true
	s.mu.Unlock()
	if s.session != nil {
		return s.session.Close()
	}
	return nil
}

// Send implements ports.SMSProvider. Long messages are split with UDH and
// every part is submitted; the first part's SMSC id is the message id.
func (s *Session) Send(ctx context.Context, msg ports.SMSMessage) (*ports.SendResult, error) {
	if s.session == nil {
		return nil, &ports.ProviderError{Code: "smpp_unbound", Message: "session not bound", Retryable: true}
	}
	submit := pdu.NewSubmitSM().(*pdu.SubmitSM)
	src := msg.Sender
	if src == "" {
		src = s.cfg.SourceAddr
	}
	srcAddr, err := pdu.NewAddressWithTonNpiAddr(byte(s.cfg.SourceTON), byte(s.cfg.SourceNPI), src) //nolint:gosec // 0-255
	if err != nil {
		return nil, &ports.ProviderError{Code: "invalid_sender", Message: err.Error()}
	}
	dst := strings.TrimPrefix(msg.To, "+")
	dstTON, dstNPI := s.cfg.DestTON, s.cfg.DestNPI
	if dstTON == 0 && dstNPI == 0 {
		dstTON, dstNPI = 1, 1 // international, ISDN
	}
	dstAddr, err := pdu.NewAddressWithTonNpiAddr(byte(dstTON), byte(dstNPI), dst) //nolint:gosec // 0-255
	if err != nil {
		return nil, &ports.ProviderError{Code: "invalid_recipient", Message: err.Error()}
	}
	submit.SourceAddr, submit.DestAddr = srcAddr, dstAddr
	if s.cfg.RequestDLR {
		submit.RegisteredDelivery = 1
	}
	enc := chooseEncoding(msg.Text)
	if err := submit.Message.SetLongMessageWithEnc(msg.Text, enc); err != nil {
		return nil, &ports.ProviderError{Code: "encoding_error", Message: err.Error()}
	}
	parts, err := submit.Split()
	if err != nil {
		return nil, &ports.ProviderError{Code: "split_error", Message: err.Error()}
	}

	timeout := time.Duration(s.cfg.SubmitTimeoutSec) * time.Second
	var firstID string
	ids := make([]string, 0, len(parts))
	for i, part := range parts {
		resp, err := s.submitAndWait(ctx, part, timeout)
		if err != nil {
			return nil, err
		}
		ids = append(ids, resp.MessageID)
		if i == 0 {
			firstID = resp.MessageID
		}
	}
	return &ports.SendResult{ProviderMessageID: firstID, Raw: map[string]any{"parts": len(parts), "message_ids": ids, "encoding": encodingName(enc)}}, nil
}

func (s *Session) submitAndWait(ctx context.Context, p *pdu.SubmitSM, timeout time.Duration) (*pdu.SubmitSMResp, *ports.ProviderError) {
	ch := make(chan *pdu.SubmitSMResp, 1)
	seq := p.GetSequenceNumber()
	s.mu.Lock()
	s.pending[seq] = ch
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		delete(s.pending, seq)
		s.mu.Unlock()
	}()
	if err := s.session.Transceiver().Submit(p); err != nil {
		return nil, &ports.ProviderError{Code: "smpp_submit", Message: err.Error(), Retryable: true}
	}
	select {
	case resp := <-ch:
		if resp == nil {
			return nil, &ports.ProviderError{Code: "smpp_submit", Message: "submit failed", Retryable: true}
		}
		if !resp.IsOk() {
			status := resp.GetHeader().CommandStatus
			return nil, classifyStatus(status)
		}
		return resp, nil
	case <-time.After(timeout):
		return nil, &ports.ProviderError{Code: "smpp_timeout", Message: "no submit_sm_resp within timeout", Retryable: true}
	case <-ctx.Done():
		return nil, &ports.ProviderError{Code: "cancelled", Message: ctx.Err().Error(), Retryable: true}
	}
}

func (s *Session) resolve(seq int32, resp *pdu.SubmitSMResp) {
	s.mu.Lock()
	ch, ok := s.pending[seq]
	s.mu.Unlock()
	if ok {
		select {
		case ch <- resp:
		default:
		}
	}
}

// onPDU handles everything the SMSC sends. deliver_sm is already
// acknowledged by gosmpp (responded == true).
func (s *Session) onPDU(p pdu.PDU, _ bool) {
	switch x := p.(type) {
	case *pdu.SubmitSMResp:
		s.resolve(x.GetSequenceNumber(), x)
	case *pdu.DeliverSM:
		if x.EsmClass&data.SM_SMSC_DLV_RCPT_TYPE == 0 && len(x.OptionalParameters) == 0 {
			// Mobile-originated message, not a receipt; ignore for now.
			return
		}
		r := parseReceipt(x)
		if r.MessageID == "" || s.onDLR == nil {
			return
		}
		go s.onDLR(context.Background(), r)
	}
}

var (
	reID   = regexp.MustCompile(`(?i)\bid:([^\s]+)`)
	reStat = regexp.MustCompile(`(?i)\bstat:([A-Z]+)`)
)

// parseReceipt reads the receipted_message_id / message_state TLVs and
// falls back to the standard "id:... stat:..." text format.
func parseReceipt(d *pdu.DeliverSM) Receipt {
	r := Receipt{}
	text, _ := d.Message.GetMessage()
	r.Text = text
	if f, ok := d.OptionalParameters[pdu.TagReceiptedMessageID]; ok {
		r.MessageID = strings.TrimRight(string(f.Data), "\x00")
	}
	if f, ok := d.OptionalParameters[pdu.TagMessageStateOption]; ok && len(f.Data) == 1 {
		r.Status = messageState(f.Data[0])
	}
	if r.MessageID == "" {
		if m := reID.FindStringSubmatch(text); len(m) == 2 {
			r.MessageID = m[1]
		}
	}
	if r.Status == "" {
		if m := reStat.FindStringSubmatch(text); len(m) == 2 {
			r.Status = strings.ToUpper(m[1])
		}
	}
	switch r.Status {
	case "DELIVRD":
		r.Delivered, r.Final = true, true
	case "UNDELIV", "EXPIRED", "REJECTD", "DELETED", "UNKNOWN":
		r.Final = true
	}
	return r
}

func messageState(b byte) string {
	switch b {
	case 1:
		return "ENROUTE"
	case 2:
		return "DELIVRD"
	case 3:
		return "EXPIRED"
	case 4:
		return "DELETED"
	case 5:
		return "UNDELIV"
	case 6:
		return "ACCEPTD"
	case 7:
		return "UNKNOWN"
	case 8:
		return "REJECTD"
	}
	return ""
}

// chooseEncoding uses GSM 7-bit when every rune fits, UCS2 otherwise
// (Turkmen letters such as ä, ň, ş, ý, ž and all Cyrillic need UCS2).
func chooseEncoding(text string) data.Encoding {
	if _, err := data.GSM7BIT.Encode(text); err == nil && isGSM7(text) {
		return data.GSM7BIT
	}
	return data.UCS2
}

const gsm7 = "@£$¥èéùìòÇ\nØø\rÅåΔ_ΦΓΛΩΠΨΣΘΞÆæßÉ !\"#¤%&'()*+,-./0123456789:;<=>?¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüà^{}\\[~]|€"

func isGSM7(text string) bool {
	for _, r := range text {
		if !strings.ContainsRune(gsm7, r) {
			return false
		}
	}
	return true
}

func encodingName(e data.Encoding) string {
	if e == data.UCS2 {
		return "UCS2"
	}
	return "GSM7"
}

// classifyStatus maps SMPP command_status values to retryable / permanent.
func classifyStatus(status data.CommandStatusType) *ports.ProviderError {
	msg := fmt.Sprintf("submit_sm_resp status 0x%08X", uint32(status))
	switch status {
	case data.ESME_RINVDSTADR, data.ESME_RINVDSTTON, data.ESME_RINVDSTNPI:
		return &ports.ProviderError{Code: "invalid_recipient", Message: msg}
	case data.ESME_RINVSRCADR, data.ESME_RINVSRCTON, data.ESME_RINVSRCNPI:
		return &ports.ProviderError{Code: "invalid_sender", Message: msg}
	case data.ESME_RINVMSGLEN, data.ESME_RINVESMCLASS:
		return &ports.ProviderError{Code: "invalid_message", Message: msg}
	case data.ESME_RTHROTTLED, data.ESME_RMSGQFUL, data.ESME_RSYSERR:
		return &ports.ProviderError{Code: "provider_unavailable", Message: msg, Retryable: true}
	}
	return &ports.ProviderError{Code: "smpp_error", Message: msg, Retryable: false}
}

// ParseConfig decodes credentials JSON.
func ParseConfig(raw json.RawMessage) (Config, error) {
	var cfg Config
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
