// Command smpptest binds to an SMPP SMSC and optionally sends one SMS, so a
// 0796-style operator gateway can be verified from a machine that can reach
// it (the API/worker do the same over the smpp adapter). Credentials are read
// from the environment — the same SMPP_* names an operator hands out — so
// nothing secret has to be typed on the command line or committed:
//
//	SMPP_IP / SMPP_HOST      SMSC host            (119.235.115.204)
//	SMPP_PORT                SMSC port            (15019)
//	SMPP_USERNAME / SMPP_SYSTEM_ID   bind system_id   (Tstandart)
//	SMPP_PASSWORD            bind password
//	SMPP_FROM / SMPP_SOURCE_ADDR     sender / short code (0796)
//
// Usage:
//
//	go run ./cmd/smpptest                         # bind only: proves login works
//	go run ./cmd/smpptest -to 65656585 -text "Salam"   # bind + send one SMS
//	go run ./cmd/smpptest -to +99365656585 -dlr        # also ask for a delivery receipt
//
// It exits non-zero on failure so it is usable in a script.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/providers/sms/smpp"
	"github.com/Esca6585dev/habarchy/backend/internal/ports"
	"github.com/Esca6585dev/habarchy/backend/pkg/phone"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	to := flag.String("to", "", "recipient phone number; bare Turkmen numbers (65656585) are read as +993…")
	text := flag.String("text", "Habarçy SMPP synagy", "message text when -to is set (Cyrillic/Turkmen is sent as UCS2)")
	useTLS := flag.Bool("tls", envBool("SMPP_TLS"), "connect over TLS")
	dlr := flag.Bool("dlr", false, "request a delivery receipt (registered_delivery=1) and wait for it")
	wait := flag.Duration("wait", 20*time.Second, "how long to wait for the delivery receipt with -dlr")
	flag.Parse()

	port, err := strconv.Atoi(firstEnv("SMPP_PORT"))
	if err != nil {
		return fmt.Errorf("SMPP_PORT must be a number: %w", err)
	}
	cfg := smpp.Config{
		Host:       firstEnv("SMPP_HOST", "SMPP_IP"),
		Port:       port,
		SystemID:   firstEnv("SMPP_SYSTEM_ID", "SMPP_USERNAME"),
		Password:   firstEnv("SMPP_PASSWORD"),
		SystemType: firstEnv("SMPP_SYSTEM_TYPE"),
		SourceAddr: firstEnv("SMPP_SOURCE_ADDR", "SMPP_FROM"),
		UseTLS:     *useTLS,
		RequestDLR: *dlr,
	}
	if cfg.Host == "" || cfg.SystemID == "" {
		return fmt.Errorf("set SMPP_IP (or SMPP_HOST), SMPP_PORT, SMPP_USERNAME (or SMPP_SYSTEM_ID) and SMPP_PASSWORD")
	}

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen}).With().Timestamp().Logger()

	// Collect delivery receipts (deliver_sm) so -dlr can report one.
	receipts := make(chan smpp.Receipt, 4)
	onDLR := func(_ context.Context, r smpp.Receipt) {
		select {
		case receipts <- r:
		default:
		}
	}

	ctx := context.Background()
	fmt.Printf("binding to %s:%d as %q (TLS=%v)…\n", cfg.Host, cfg.Port, cfg.SystemID, cfg.UseTLS)
	sess, err := smpp.Dial(ctx, cfg, logger, onDLR)
	if err != nil {
		return fmt.Errorf("bind failed: %w", err)
	}
	defer func() { _ = sess.Close() }()
	fmt.Println("✓ bind OK (login accepted)")

	if *to == "" {
		fmt.Println("no -to given, nothing sent. Login works.")
		return nil
	}

	dst, err := phone.Normalize(*to, "")
	if err != nil {
		return fmt.Errorf("invalid -to %q: %w", *to, err)
	}
	sender := cfg.SourceAddr
	fmt.Printf("sending %q → %s (from %q)…\n", *text, dst, sender)
	res, err := sess.Send(ctx, ports.SMSMessage{To: dst, Text: *text, Sender: sender})
	if err != nil {
		return fmt.Errorf("submit failed: %w", err)
	}
	fmt.Printf("✓ submitted: provider_message_id=%s raw=%v\n", res.ProviderMessageID, res.Raw)

	if *dlr {
		fmt.Printf("waiting %s for a delivery receipt…\n", *wait)
		select {
		case r := <-receipts:
			fmt.Printf("✓ receipt: status=%s delivered=%v id=%s\n", r.Status, r.Delivered, r.MessageID)
		case <-time.After(*wait):
			fmt.Println("… no receipt within the wait window (the SMS may still arrive; some SMSCs are slow)")
		}
	}
	return nil
}

// firstEnv returns the first non-empty value among the named variables.
func firstEnv(names ...string) string {
	for _, n := range names {
		if v := strings.TrimSpace(os.Getenv(n)); v != "" {
			return v
		}
	}
	return ""
}

func envBool(name string) bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(name)))
	return v == "1" || v == "true" || v == "yes"
}
