package config

import (
	"strings"
	"testing"
)

func setValid(t *testing.T) {
	t.Helper()
	t.Setenv("HABARCHY_DATABASE_URL", "postgres://u:p@localhost:5432/habarchy")
	t.Setenv("HABARCHY_JWT_SECRET", strings.Repeat("s", 32))
	t.Setenv("HABARCHY_MASTER_KEY", strings.Repeat("a", 64))
}

func TestLoadDefaults(t *testing.T) {
	setValid(t)
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.HTTP.Addr != ":8080" || c.Env != "development" || c.Queue.MaxRetry != 3 {
		t.Fatalf("unexpected defaults: %+v", c)
	}
	if c.Redis.URL != "redis://localhost:6379/0" {
		t.Fatalf("redis default: %s", c.Redis.URL)
	}
	if len(c.HTTP.CORSOrigins) != 1 {
		t.Fatalf("cors default: %v", c.HTTP.CORSOrigins)
	}
}

func TestLoadRequiresDatabase(t *testing.T) {
	setValid(t)
	t.Setenv("HABARCHY_DATABASE_URL", "")
	if _, err := Load(); err == nil {
		t.Fatal("expected error for missing DATABASE_URL")
	}
}

func TestValidate(t *testing.T) {
	setValid(t)
	t.Setenv("HABARCHY_JWT_SECRET", "short")
	t.Setenv("HABARCHY_OTP_LENGTH", "2")
	t.Setenv("HABARCHY_ENV", "staging")
	_, err := Load()
	if err == nil {
		t.Fatal("expected validation errors")
	}
	for _, want := range []string{"JWT_SECRET", "OTP_LENGTH", "ENV must be"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q should mention %s", err, want)
		}
	}
}
