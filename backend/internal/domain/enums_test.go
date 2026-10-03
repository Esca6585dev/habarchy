package domain

import (
	"errors"
	"fmt"
	"testing"
)

func TestParseChannel(t *testing.T) {
	for _, in := range []string{"sms", "email", "push", "telegram", "auto"} {
		if _, err := ParseChannel(in); err != nil {
			t.Fatalf("ParseChannel(%q) unexpected error: %v", in, err)
		}
	}
	if _, err := ParseChannel("fax"); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
	if ChannelAuto.Valid() {
		t.Fatal("auto must not be a deliverable channel")
	}
}

func TestMessageStatusTransitions(t *testing.T) {
	cases := []struct {
		from, to MessageStatus
		ok       bool
	}{
		{StatusQueued, StatusProcessing, true},
		{StatusQueued, StatusCancelled, true},
		{StatusQueued, StatusDelivered, false},
		{StatusProcessing, StatusSent, true},
		{StatusProcessing, StatusQueued, true},
		{StatusSent, StatusDelivered, true},
		{StatusSent, StatusQueued, false},
		{StatusDelivered, StatusFailed, false},
		{StatusFailed, StatusQueued, false},
		{StatusCancelled, StatusProcessing, false},
	}
	for _, c := range cases {
		if got := c.from.CanTransitionTo(c.to); got != c.ok {
			t.Errorf("%s -> %s: got %v want %v", c.from, c.to, got, c.ok)
		}
	}
	for _, s := range []MessageStatus{StatusDelivered, StatusFailed, StatusCancelled} {
		if !s.Terminal() {
			t.Errorf("%s should be terminal", s)
		}
	}
}

func TestProviderTypeChannel(t *testing.T) {
	want := map[ProviderType]Channel{
		ProviderHTTPSMS: ChannelSMS, ProviderSMPP: ChannelSMS, ProviderSMTP: ChannelEmail,
		ProviderFCM: ChannelPush, ProviderTelegramBot: ChannelTelegram,
	}
	for pt, ch := range want {
		if pt.Channel() != ch || !pt.Valid() {
			t.Errorf("%s: channel %s, want %s", pt, pt.Channel(), ch)
		}
	}
	if ProviderType("carrier_pigeon").Valid() {
		t.Error("unknown provider type must be invalid")
	}
}

func TestRoleAtLeast(t *testing.T) {
	if !RoleOwner.AtLeast(RoleViewer) || !RoleAdmin.AtLeast(RoleAdmin) {
		t.Error("higher roles must satisfy lower requirements")
	}
	if RoleViewer.AtLeast(RoleDeveloper) || Role("ghost").AtLeast(RoleViewer) {
		t.Error("lower or unknown roles must not satisfy higher requirements")
	}
}

func TestDomainErrorMatching(t *testing.T) {
	wrapped := fmt.Errorf("sending: %w", ErrQuotaExceeded.WithMessage("daily quota of 100 reached"))
	if !errors.Is(wrapped, ErrQuotaExceeded) {
		t.Fatal("wrapped derived error must match sentinel")
	}
	if CodeOf(wrapped) != CodeQuotaExceeded {
		t.Fatalf("CodeOf = %s", CodeOf(wrapped))
	}
	if CodeOf(errors.New("boom")) != CodeInternal {
		t.Fatal("unknown errors must map to internal_error")
	}
	if CodeQuotaExceeded.HTTPStatus() != 429 || CodeNotFound.HTTPStatus() != 404 || CodeInternal.HTTPStatus() != 500 {
		t.Fatal("unexpected HTTP status mapping")
	}
	d := ErrValidation.WithDetails(map[string]any{"field": "to"})
	if d.Details["field"] != "to" || ErrValidation.Details != nil {
		t.Fatal("WithDetails must not mutate the sentinel")
	}
}
