package smpp

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/pdu"
)

func TestChooseEncoding(t *testing.T) {
	// ä, ö, ü exist in the GSM 7-bit alphabet; ň, ş, ý, ž and Cyrillic do not.
	for _, s := range []string{"Hello 123", "täze", "söz"} {
		if chooseEncoding(s) != data.GSM7BIT {
			t.Fatalf("%q should be GSM7", s)
		}
	}
	for _, s := range []string{"Siziň koduňyz", "Код: 1234", "ýaz", "žurnal"} {
		if chooseEncoding(s) != data.UCS2 {
			t.Fatalf("%q should be UCS2", s)
		}
	}
}

func TestParseReceiptText(t *testing.T) {
	d := pdu.NewDeliverSM().(*pdu.DeliverSM)
	d.EsmClass = data.SM_SMSC_DLV_RCPT_TYPE
	_ = d.Message.SetMessageWithEncoding("id:ABC123 sub:001 dlvrd:001 submit date:2410041200 done date:2410041201 stat:DELIVRD err:000 text:Kod", data.GSM7BIT)
	r := parseReceipt(d)
	if r.MessageID != "ABC123" || r.Status != "DELIVRD" || !r.Delivered || !r.Final {
		t.Fatalf("%+v", r)
	}
	_ = d.Message.SetMessageWithEncoding("id:XYZ stat:EXPIRED err:001", data.GSM7BIT)
	r = parseReceipt(d)
	if r.MessageID != "XYZ" || r.Delivered || !r.Final {
		t.Fatalf("%+v", r)
	}
}

func TestParseReceiptTLV(t *testing.T) {
	d := pdu.NewDeliverSM().(*pdu.DeliverSM)
	d.RegisterOptionalParam(pdu.Field{Tag: pdu.TagReceiptedMessageID, Data: []byte("TLV-77\x00")})
	d.RegisterOptionalParam(pdu.Field{Tag: pdu.TagMessageStateOption, Data: []byte{5}})
	r := parseReceipt(d)
	if r.MessageID != "TLV-77" || r.Status != "UNDELIV" || r.Delivered || !r.Final {
		t.Fatalf("%+v", r)
	}
}

func TestClassifyStatus(t *testing.T) {
	if classifyStatus(data.ESME_RINVDSTADR).Retryable || classifyStatus(data.ESME_RINVDSTADR).Code != "invalid_recipient" {
		t.Fatal("invalid dest must be permanent")
	}
	if !classifyStatus(data.ESME_RTHROTTLED).Retryable {
		t.Fatal("throttled must be retryable")
	}
}

func TestLongMessageSplit(t *testing.T) {
	submit := pdu.NewSubmitSM().(*pdu.SubmitSM)
	long := ""
	for i := 0; i < 20; i++ {
		long += "Siziň koduňyz 1234. "
	}
	if err := submit.Message.SetLongMessageWithEnc(long, data.UCS2); err != nil {
		t.Fatal(err)
	}
	parts, err := submit.Split()
	if err != nil || len(parts) < 2 {
		t.Fatalf("expected multipart, got %d (%v)", len(parts), err)
	}
	for _, p := range parts {
		if p.EsmClass&data.SM_UDH_GSM == 0 {
			t.Fatal("UDHI must be set on concatenated parts")
		}
	}
}
