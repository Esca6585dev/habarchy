// Command import-smssender copies contacts from the sms-sender-app backend
// (https://atom.com.tm/elitepanel/sms, the "other repository") into a Habarçy
// tenant. It reads the source over its HTTP API with the Integration-Key
// header, normalizes each phone to E.164, and upserts it as a contact in the
// chosen tenant. Idempotent: a contact's external_id is "smssender:<phone>",
// so re-running updates instead of duplicating, and numbers stay inside the
// one tenant (never merged with another's).
//
// Usage:
//
//	go run ./cmd/import-smssender -project demo -key "<Integration-Key>"
//	go run ./cmd/import-smssender -project demo -key "$SMSSENDER_KEY" -dry-run
//
// -source defaults to the sms-sender backend; override for a staging host.
// DB connection comes from the usual HABARCHY_* env (same .env as the API).
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
	"github.com/Esca6585dev/habarchy/backend/internal/config"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/pkg/phone"
)

// srcContact mirrors the sms-sender backend contact shape.
type srcContact struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Phone string `json:"phone_number"`
	Sent  bool   `json:"is_sent"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	slug := flag.String("project", "", "target Habarçy tenant slug to import into (required)")
	source := flag.String("source", envOr("SMSSENDER_URL", "https://atom.com.tm/elitepanel/sms"), "sms-sender backend base URL")
	key := flag.String("key", os.Getenv("SMSSENDER_KEY"), "Integration-Key header for the source API (required)")
	tag := flag.String("tag", "smssender", "tag added to every imported contact")
	dryRun := flag.Bool("dry-run", false, "fetch and report only; write nothing")
	limit := flag.Int("limit", 0, "stop after importing this many (0 = all)")
	timeout := flag.Duration("timeout", 60*time.Second, "HTTP timeout for the source API")
	flag.Parse()

	if strings.TrimSpace(*slug) == "" {
		return errors.New("set -project <tenant slug> (the Habarçy tenant to import into)")
	}
	if strings.TrimSpace(*key) == "" {
		return errors.New("set -key or SMSSENDER_KEY (the source Integration-Key)")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	ctx := context.Background()
	db, err := postgres.Connect(ctx, cfg.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	project, err := db.Queries.GetProjectBySlug(ctx, strings.TrimSpace(*slug))
	if err != nil {
		return fmt.Errorf("tenant %q not found (create it or run -seed): %w", *slug, err)
	}
	svc := contacts.New(db)

	base := strings.TrimRight(strings.TrimSpace(*source), "/")
	rows, err := fetchAll(ctx, base, *key, *timeout)
	if err != nil {
		return err
	}
	fmt.Printf("fetched %d contacts from %s → tenant %q\n", len(rows), base, project.Slug)

	var created, updated, invalid, blank int
	for i := range rows {
		if *limit > 0 && created+updated >= *limit {
			break
		}
		raw := strings.TrimSpace(rows[i].Phone)
		if raw == "" {
			blank++
			continue
		}
		if *dryRun {
			if _, err := phone.Normalize(raw, ""); err != nil {
				invalid++
			}
			continue
		}
		_, isNew, err := svc.Upsert(ctx, project.ID, contacts.Input{
			ExternalID: "smssender:" + raw,
			Name:       strings.TrimSpace(rows[i].Title),
			Phone:      raw,
			Tags:       []string{*tag},
			Attributes: map[string]any{"source": "smssender", "is_sent": rows[i].Sent, "source_id": rows[i].ID},
		})
		if err != nil {
			if domain.CodeOf(err) == domain.CodeValidationFailed {
				invalid++
				continue
			}
			return fmt.Errorf("contact %q: %w", raw, err)
		}
		if isNew {
			created++
		} else {
			updated++
		}
	}
	fmt.Printf("done: created=%d updated=%d invalid=%d blank=%d dry_run=%v\n", created, updated, invalid, blank, *dryRun)
	if invalid > 0 {
		fmt.Printf("note: %d numbers could not be normalized to a valid phone and were skipped\n", invalid)
	}
	return nil
}

// fetchAll pulls every contact from GET /contacts/all/ (a JSON array). A
// paginated object ({"results": [...]}) is also accepted.
func fetchAll(ctx context.Context, base, key string, timeout time.Duration) ([]srcContact, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+"/contacts/all/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Integration-Key", key)
	req.Header.Set("Accept", "application/json")
	resp, err := (&http.Client{Timeout: timeout}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("source request failed (is the host reachable from here?): %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<20))
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("source returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	// Try a bare array first, then a {"results": [...]} envelope.
	var arr []srcContact
	if err := json.Unmarshal(body, &arr); err == nil {
		return arr, nil
	}
	var env struct {
		Results []srcContact `json:"results"`
	}
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("unexpected source response: %w", err)
	}
	return env.Results, nil
}

func envOr(name, def string) string {
	if v := strings.TrimSpace(os.Getenv(name)); v != "" {
		return v
	}
	return def
}
