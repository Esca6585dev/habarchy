package contactimport

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
)

// CardDAVClient reads address books over CardDAV (RFC 6352): iCloud with an
// Apple ID + app-specific password, Nextcloud, Fastmail, Zimbra…
type CardDAVClient struct {
	HTTP *http.Client
}

// ICloudURL is Apple's CardDAV entry point.
const ICloudURL = "https://contacts.icloud.com/"

// NewCardDAVClient creates a client (nil → default HTTP client).
func NewCardDAVClient(h *http.Client) *CardDAVClient {
	if h == nil {
		h = &http.Client{Timeout: 60 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	}
	return &CardDAVClient{HTTP: h}
}

// Fetch discovers the user's address books and returns every contact.
func (c *CardDAVClient) Fetch(ctx context.Context, serverURL, username, password string) ([]contacts.Input, error) {
	serverURL = strings.TrimSpace(serverURL)
	if serverURL == "" {
		serverURL = ICloudURL
	}
	if !strings.HasPrefix(serverURL, "https://") && !strings.HasPrefix(serverURL, "http://localhost") && !strings.HasPrefix(serverURL, "http://127.0.0.1") {
		return nil, errors.New("carddav: server_url must use https")
	}
	base, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}
	auth := func(req *http.Request) { req.SetBasicAuth(username, password) }

	// 1. current-user-principal
	principal, err := c.propfind(ctx, base, "0", `<d:propfind xmlns:d="DAV:"><d:prop><d:current-user-principal/></d:prop></d:propfind>`, auth)
	if err != nil {
		return nil, err
	}
	principalHref := principal.firstHref("current-user-principal")
	if principalHref == "" {
		return nil, errors.New("carddav: login failed or server has no principal (check Apple ID and app-specific password)")
	}
	principalURL := resolve(base, principalHref)

	// 2. addressbook-home-set
	home, err := c.propfind(ctx, principalURL, "0", `<d:propfind xmlns:d="DAV:" xmlns:card="urn:ietf:params:xml:ns:carddav"><d:prop><card:addressbook-home-set/></d:prop></d:propfind>`, auth)
	if err != nil {
		return nil, err
	}
	homeHref := home.firstHref("addressbook-home-set")
	if homeHref == "" {
		return nil, errors.New("carddav: no address book home")
	}
	homeURL := resolve(principalURL, homeHref)

	// 3. address book collections
	list, err := c.propfind(ctx, homeURL, "1", `<d:propfind xmlns:d="DAV:"><d:prop><d:resourcetype/><d:displayname/></d:prop></d:propfind>`, auth)
	if err != nil {
		return nil, err
	}
	var books []*url.URL
	for _, r := range list.Responses {
		if r.isAddressbook() {
			books = append(books, resolve(homeURL, r.Href))
		}
	}
	if len(books) == 0 {
		return nil, errors.New("carddav: no address books found")
	}

	// 4. all vCards of every book
	var out []contacts.Input
	for _, book := range books {
		ms, err := c.report(ctx, book, auth)
		if err != nil {
			return nil, err
		}
		for _, r := range ms.Responses {
			data := r.addressData()
			if data == "" {
				continue
			}
			ins, _, err := ParseVCard(strings.NewReader(data))
			if err != nil {
				continue
			}
			out = append(out, ins...)
			if len(out) >= MaxRows {
				return out, nil
			}
		}
	}
	if len(out) == 0 {
		return nil, errEmpty
	}
	return out, nil
}

type multistatus struct {
	Responses []davResponse `xml:"response"`
}

type davResponse struct {
	Href     string `xml:"href"`
	Propstat []struct {
		Prop struct {
			CurrentUserPrincipal struct {
				Href string `xml:"href"`
			} `xml:"current-user-principal"`
			AddressbookHomeSet struct {
				Href string `xml:"href"`
			} `xml:"addressbook-home-set"`
			ResourceType struct {
				Inner []xml.Name `xml:",any"`
			} `xml:"resourcetype"`
			AddressData string `xml:"address-data"`
		} `xml:"prop"`
	} `xml:"propstat"`
}

func (m *multistatus) firstHref(prop string) string {
	for _, r := range m.Responses {
		for _, ps := range r.Propstat {
			switch prop {
			case "current-user-principal":
				if ps.Prop.CurrentUserPrincipal.Href != "" {
					return ps.Prop.CurrentUserPrincipal.Href
				}
			case "addressbook-home-set":
				if ps.Prop.AddressbookHomeSet.Href != "" {
					return ps.Prop.AddressbookHomeSet.Href
				}
			}
		}
	}
	return ""
}

func (r davResponse) isAddressbook() bool {
	for _, ps := range r.Propstat {
		for _, n := range ps.Prop.ResourceType.Inner {
			if n.Local == "addressbook" {
				return true
			}
		}
	}
	return false
}

func (r davResponse) addressData() string {
	for _, ps := range r.Propstat {
		if ps.Prop.AddressData != "" {
			return ps.Prop.AddressData
		}
	}
	return ""
}

func (c *CardDAVClient) propfind(ctx context.Context, u *url.URL, depth, body string, auth func(*http.Request)) (*multistatus, error) {
	return c.dav(ctx, "PROPFIND", u, depth, body, auth)
}

func (c *CardDAVClient) report(ctx context.Context, u *url.URL, auth func(*http.Request)) (*multistatus, error) {
	body := `<card:addressbook-query xmlns:d="DAV:" xmlns:card="urn:ietf:params:xml:ns:carddav"><d:prop><d:getetag/><card:address-data/></d:prop></card:addressbook-query>`
	return c.dav(ctx, "REPORT", u, "1", body, auth)
}

// dav performs one WebDAV request, following redirects manually so the
// Authorization header survives host changes (iCloud → pXX-contacts.icloud.com).
func (c *CardDAVClient) dav(ctx context.Context, method string, u *url.URL, depth, body string, auth func(*http.Request)) (*multistatus, error) {
	cur := u
	for hop := 0; hop < 6; hop++ {
		req, err := http.NewRequestWithContext(ctx, method, cur.String(), strings.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/xml; charset=utf-8")
		req.Header.Set("Depth", depth)
		req.Header.Set("User-Agent", "Habarchy/1.0 CardDAV")
		auth(req)
		res, err := c.HTTP.Do(req)
		if err != nil {
			return nil, fmt.Errorf("carddav: %w", err)
		}
		raw, _ := io.ReadAll(io.LimitReader(res.Body, 50<<20))
		_ = res.Body.Close()
		switch {
		case res.StatusCode >= 300 && res.StatusCode < 400:
			loc := res.Header.Get("Location")
			if loc == "" {
				return nil, fmt.Errorf("carddav: redirect without location")
			}
			cur = resolve(cur, loc)
			continue
		case res.StatusCode == 401 || res.StatusCode == 403:
			return nil, errors.New("carddav: authentication failed (iCloud needs an app-specific password from appleid.apple.com)")
		case res.StatusCode == 207 || res.StatusCode == 200:
			var ms multistatus
			if err := xml.Unmarshal(raw, &ms); err != nil {
				return nil, fmt.Errorf("carddav: bad response: %w", err)
			}
			return &ms, nil
		default:
			return nil, fmt.Errorf("carddav: %s %s → %d", method, cur.Path, res.StatusCode)
		}
	}
	return nil, errors.New("carddav: too many redirects")
}

func resolve(base *url.URL, href string) *url.URL {
	ref, err := url.Parse(strings.TrimSpace(href))
	if err != nil {
		return base
	}
	return base.ResolveReference(ref)
}
