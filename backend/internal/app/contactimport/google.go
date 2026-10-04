package contactimport

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
)

// GoogleConnector imports Google Contacts through OAuth 2.0 + People API.
// Tokens are used once for the import and never stored.
type GoogleConnector struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AuthURL      string // default https://accounts.google.com/o/oauth2/v2/auth
	TokenURL     string // default https://oauth2.googleapis.com/token
	PeopleURL    string // default https://people.googleapis.com/v1/people/me/connections
	HTTP         *http.Client
	stateSecret  []byte
}

// NewGoogleConnector configures the connector; stateSecret signs the OAuth
// state (use the JWT secret).
func NewGoogleConnector(clientID, clientSecret, redirectURL string, stateSecret []byte) *GoogleConnector {
	return &GoogleConnector{
		ClientID: clientID, ClientSecret: clientSecret, RedirectURL: redirectURL,
		AuthURL: "https://accounts.google.com/o/oauth2/v2/auth", TokenURL: "https://oauth2.googleapis.com/token",
		PeopleURL: "https://people.googleapis.com/v1/people/me/connections", HTTP: &http.Client{Timeout: 30 * time.Second}, stateSecret: stateSecret,
	}
}

// Configured reports whether client credentials are present.
func (g *GoogleConnector) Configured() bool {
	return g != nil && g.ClientID != "" && g.ClientSecret != ""
}

// State is what the OAuth round-trip carries (signed, short-lived).
type State struct {
	ProjectID uuid.UUID  `json:"p"`
	UserID    uuid.UUID  `json:"u"`
	GroupID   *uuid.UUID `json:"g,omitempty"`
	ReturnTo  string     `json:"r,omitempty"`
	Exp       int64      `json:"e"`
}

// SignState encodes and signs a state.
func (g *GoogleConnector) SignState(st State) string {
	if st.Exp == 0 {
		st.Exp = time.Now().Add(15 * time.Minute).Unix()
	}
	b, _ := json.Marshal(st)
	payload := base64.RawURLEncoding.EncodeToString(b)
	return payload + "." + g.mac(payload)
}

// ParseState verifies and decodes a state.
func (g *GoogleConnector) ParseState(s string) (*State, error) {
	payload, sig, ok := strings.Cut(s, ".")
	if !ok || !hmac.Equal([]byte(sig), []byte(g.mac(payload))) {
		return nil, errors.New("invalid state")
	}
	b, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, errors.New("invalid state")
	}
	var st State
	if err := json.Unmarshal(b, &st); err != nil || time.Now().Unix() > st.Exp {
		return nil, errors.New("state expired")
	}
	return &st, nil
}

func (g *GoogleConnector) mac(payload string) string {
	m := hmac.New(sha256.New, g.stateSecret)
	m.Write([]byte(payload))
	return hex.EncodeToString(m.Sum(nil))
}

// AuthorizationURL builds the consent URL (read-only contacts scope).
func (g *GoogleConnector) AuthorizationURL(state string) string {
	q := url.Values{
		"client_id": {g.ClientID}, "redirect_uri": {g.RedirectURL}, "response_type": {"code"},
		"scope": {"https://www.googleapis.com/auth/contacts.readonly"}, "access_type": {"online"}, "prompt": {"select_account"}, "state": {state},
	}
	return g.AuthURL + "?" + q.Encode()
}

// Exchange trades the code for an access token.
func (g *GoogleConnector) Exchange(ctx context.Context, code string) (string, error) {
	form := url.Values{"code": {code}, "client_id": {g.ClientID}, "client_secret": {g.ClientSecret}, "redirect_uri": {g.RedirectURL}, "grant_type": {"authorization_code"}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := g.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = res.Body.Close() }()
	var out struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
		Desc        string `json:"error_description"`
	}
	_ = json.NewDecoder(io.LimitReader(res.Body, 1<<20)).Decode(&out)
	if res.StatusCode != 200 || out.AccessToken == "" {
		return "", fmt.Errorf("google token: %s %s", out.Error, out.Desc)
	}
	return out.AccessToken, nil
}

// FetchContacts lists the user's connections (names, phones, e-mails).
func (g *GoogleConnector) FetchContacts(ctx context.Context, accessToken string) ([]contacts.Input, error) {
	var out []contacts.Input
	pageToken := ""
	for page := 0; page < 20; page++ {
		q := url.Values{"personFields": {"names,phoneNumbers,emailAddresses"}, "pageSize": {"1000"}}
		if pageToken != "" {
			q.Set("pageToken", pageToken)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.PeopleURL+"?"+q.Encode(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		res, err := g.HTTP.Do(req)
		if err != nil {
			return nil, err
		}
		var body struct {
			Connections []struct {
				ResourceName string `json:"resourceName"`
				Names        []struct {
					DisplayName string `json:"displayName"`
				} `json:"names"`
				PhoneNumbers []struct {
					Value         string `json:"value"`
					CanonicalForm string `json:"canonicalForm"`
					FormattedType string `json:"formattedType"`
					Type          string `json:"type"`
				} `json:"phoneNumbers"`
				EmailAddresses []struct {
					Value string `json:"value"`
				} `json:"emailAddresses"`
			} `json:"connections"`
			NextPageToken string `json:"nextPageToken"`
			Error         struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		err = json.NewDecoder(io.LimitReader(res.Body, 50<<20)).Decode(&body)
		_ = res.Body.Close()
		if err != nil {
			return nil, err
		}
		if res.StatusCode != 200 {
			return nil, fmt.Errorf("google people api: %d %s", res.StatusCode, body.Error.Message)
		}
		for _, p := range body.Connections {
			var in contacts.Input
			if len(p.Names) > 0 {
				in.Name = p.Names[0].DisplayName
			}
			for _, ph := range p.PhoneNumbers {
				v := ph.CanonicalForm
				if v == "" {
					v = ph.Value
				}
				if !isPhone(v) {
					continue
				}
				if in.Phone == "" {
					in.Phone = cleanPhone(v)
				} else if in.WhatsApp == "" {
					in.WhatsApp = cleanPhone(v)
				}
			}
			if len(p.EmailAddresses) > 0 {
				in.Email = strings.ToLower(p.EmailAddresses[0].Value)
			}
			if !hasReach(in) {
				continue
			}
			if p.ResourceName != "" {
				in.ExternalID = "google:" + strings.TrimPrefix(p.ResourceName, "people/")
			}
			out = append(out, in)
			if len(out) >= MaxRows {
				return out, nil
			}
		}
		if body.NextPageToken == "" {
			break
		}
		pageToken = body.NextPageToken
	}
	return out, nil
}
