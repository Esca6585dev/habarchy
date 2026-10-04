// Package contactimport turns files (CSV, Excel, Word, vCard, text) and
// external address books (Google Contacts, CardDAV / iCloud) into contact
// inputs, and merges them into a project.
package contactimport

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"path"
	"regexp"
	"strings"
	"unicode"

	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// MaxRows caps one import.
const MaxRows = 5000

// RowError explains why one row was skipped.
type RowError struct {
	Row    int    `json:"row"`
	Reason string `json:"reason"`
	Line   string `json:"line,omitempty"`
}

// Parse dispatches on the file extension.
func Parse(filename string, r io.Reader) ([]contacts.Input, []RowError, error) {
	ext := strings.ToLower(path.Ext(filename))
	switch ext {
	case ".csv", ".tsv", ".txt":
		return ParseText(r)
	case ".xlsx", ".xlsm":
		return ParseXLSX(r)
	case ".vcf", ".vcard":
		return ParseVCard(r)
	case ".docx":
		return ParseDOCX(r)
	}
	return nil, nil, fmt.Errorf("unsupported file type %q (use .xlsx, .csv, .vcf, .docx or .txt)", ext)
}

var (
	phoneRe   = regexp.MustCompile(`^\+?[\d\s()./-]{6,}$`)
	emailRe   = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	slackRe   = regexp.MustCompile(`^[UC][A-Z0-9]{6,}$`)
	chatIDRe  = regexp.MustCompile(`^-?\d{5,}$`)
	digitsRe  = regexp.MustCompile(`\D`)
	headerMap = map[string]string{
		"name": "name", "ady": "name", "ad": "name", "имя": "name", "фио": "name", "full name": "name", "fullname": "name", "contact": "name", "customer": "name", "müşderi": "name", "musderi": "name",
		"phone": "phone", "telefon": "phone", "tel": "phone", "mobile": "phone", "mobil": "phone", "телефон": "phone", "nomer": "phone", "number": "phone", "belgi": "phone", "msisdn": "phone", "gsm": "phone",
		"email": "email", "e-mail": "email", "mail": "email", "почта": "email", "poçta": "email",
		"whatsapp": "whatsapp", "wa": "whatsapp",
		"telegram": "telegram", "telegram_chat_id": "telegram", "chat_id": "telegram", "tg": "telegram",
		"slack": "slack", "slack_id": "slack",
		"external_id": "external_id", "id": "external_id", "user_id": "external_id", "customer_id": "external_id",
		"tags": "tags", "tag": "tags", "group": "tags", "topar": "tags", "группа": "tags",
		"locale": "locale", "lang": "locale", "language": "locale", "dil": "locale", "язык": "locale",
	}
)

// ParseText reads CSV / TSV / semicolon separated or free lines.
func ParseText(r io.Reader) ([]contacts.Input, []RowError, error) {
	data, err := io.ReadAll(io.LimitReader(r, 20<<20))
	if err != nil {
		return nil, nil, err
	}
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))
	sep := detectSeparator(data)
	var rows [][]string
	if sep != 0 {
		rd := csv.NewReader(bytes.NewReader(data))
		rd.Comma = sep
		rd.FieldsPerRecord = -1
		rd.LazyQuotes = true
		rd.TrimLeadingSpace = true
		rows, err = rd.ReadAll()
		if err != nil {
			return nil, nil, fmt.Errorf("csv: %w", err)
		}
	} else {
		sc := bufio.NewScanner(bytes.NewReader(data))
		sc.Buffer(make([]byte, 1<<20), 1<<20)
		for sc.Scan() {
			rows = append(rows, []string{sc.Text()})
		}
	}
	return FromRows(rows), nil, nil
}

func detectSeparator(data []byte) rune {
	head := data
	if len(head) > 4096 {
		head = head[:4096]
	}
	counts := map[rune]int{}
	for _, r := range string(head) {
		switch r {
		case ',', ';', '\t', '|':
			counts[r]++
		}
	}
	best, n := rune(0), 1
	for _, r := range []rune{'\t', ';', ',', '|'} {
		if counts[r] > n {
			best, n = r, counts[r]
		}
	}
	return best
}

// FromRows handles both tabular data (header row with known names) and
// free rows where every cell is classified by shape.
func FromRows(rows [][]string) []contacts.Input {
	if len(rows) == 0 {
		return nil
	}
	cols := headerColumns(rows[0])
	start := 0
	if cols != nil {
		start = 1
	}
	var out []contacts.Input
	for i := start; i < len(rows) && len(out) < MaxRows; i++ {
		var in contacts.Input
		var ok bool
		if cols != nil {
			in, ok = fromHeaderRow(cols, rows[i])
		} else {
			in, ok = FromCells(rows[i])
		}
		if ok {
			out = append(out, in)
		}
	}
	return out
}

func headerColumns(row []string) map[int]string {
	cols := map[int]string{}
	hits := 0
	for i, h := range row {
		key := strings.ToLower(strings.TrimSpace(h))
		if f, ok := headerMap[key]; ok {
			cols[i] = f
			hits++
		}
	}
	if hits == 0 || (hits == 1 && len(row) > 1 && looksLikeData(row)) {
		return nil
	}
	return cols
}

func looksLikeData(row []string) bool {
	for _, c := range row {
		c = strings.TrimSpace(c)
		if emailRe.MatchString(c) || (phoneRe.MatchString(c) && len(digitsRe.ReplaceAllString(c, "")) >= 6) {
			return true
		}
	}
	return false
}

func fromHeaderRow(cols map[int]string, row []string) (contacts.Input, bool) {
	var in contacts.Input
	for i, cell := range row {
		v := strings.TrimSpace(cell)
		if v == "" {
			continue
		}
		switch cols[i] {
		case "name":
			in.Name = v
		case "phone":
			in.Phone = cleanPhone(v)
		case "email":
			in.Email = v
		case "whatsapp":
			in.WhatsApp = cleanPhone(v)
		case "telegram":
			in.TelegramChatID = v
		case "slack":
			in.SlackID = v
		case "external_id":
			in.ExternalID = v
		case "locale":
			in.Locale = contactsLocale(v)
		case "tags":
			for _, t := range strings.FieldsFunc(v, func(r rune) bool { return r == ',' || r == ';' || r == '|' }) {
				if t = strings.TrimSpace(t); t != "" {
					in.Tags = append(in.Tags, t)
				}
			}
		}
	}
	return in, hasAddress(in)
}

// FromCells classifies free cells: phones, e-mails, Slack / Telegram ids,
// the rest becomes the name.
func FromCells(cells []string) (contacts.Input, bool) {
	var in contacts.Input
	var nameParts []string
	for _, cell := range cells {
		for _, v := range splitLoose(cell) {
			if classify(&in, v) {
				continue
			}
			// Space separated ("maral@example.tm Maral"): classify word by word.
			words := strings.Fields(v)
			if len(words) > 1 {
				for _, w := range words {
					if !classify(&in, w) {
						nameParts = append(nameParts, w)
					}
				}
				continue
			}
			nameParts = append(nameParts, v)
		}
	}
	in.Name = strings.Join(nameParts, " ")
	return in, hasReach(in)
}

// classify assigns v to the matching field; false when it is plain text.
func classify(in *contacts.Input, v string) bool {
	switch {
	case emailRe.MatchString(v):
		if in.Email == "" {
			in.Email = strings.ToLower(v)
		}
	case slackRe.MatchString(v):
		in.SlackID = v
	case isPhone(v):
		if in.Phone == "" {
			in.Phone = cleanPhone(v)
		} else if in.WhatsApp == "" {
			in.WhatsApp = cleanPhone(v)
		}
	case chatIDRe.MatchString(v):
		in.TelegramChatID = v
	default:
		return false
	}
	return true
}

// hasReach reports whether the contact can be messaged on some channel.
func hasReach(in contacts.Input) bool {
	return in.Phone != "" || in.Email != "" || in.WhatsApp != "" || in.TelegramChatID != "" || in.SlackID != ""
}

// FromLine parses one free-text line ("Aman Amanow, +99365123456, mail@…").
func FromLine(line string) (contacts.Input, bool) {
	return FromCells([]string{line})
}

func splitLoose(cell string) []string {
	var out []string
	for _, p := range strings.FieldsFunc(cell, func(r rune) bool { return r == ',' || r == ';' || r == '\t' || r == '|' }) {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func isPhone(v string) bool {
	if !phoneRe.MatchString(v) {
		return false
	}
	d := len(digitsRe.ReplaceAllString(v, ""))
	return d >= 6 && d <= 15
}

// cleanPhone keeps + and digits and fixes Excel's float rendering (9.9365E+10).
func cleanPhone(v string) string {
	v = strings.TrimSpace(v)
	if strings.ContainsAny(v, "Ee") && !strings.ContainsAny(v, "+-( ") {
		var f float64
		if _, err := fmt.Sscanf(v, "%g", &f); err == nil && f > 0 {
			return fmt.Sprintf("%.0f", f)
		}
	}
	plus := strings.HasPrefix(v, "+")
	d := digitsRe.ReplaceAllString(v, "")
	if plus {
		return "+" + d
	}
	return d
}

func hasAddress(in contacts.Input) bool {
	return in.Phone != "" || in.Email != "" || in.WhatsApp != "" || in.TelegramChatID != "" || in.SlackID != "" || in.ExternalID != ""
}

func contactsLocale(v string) domain.Locale {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "ru", "rus", "русский":
		return "ru"
	case "en", "eng", "english":
		return "en"
	case "tk", "tm", "tkm", "türkmen", "turkmen":
		return "tk"
	}
	return ""
}

// errEmpty is returned when a file yields nothing.
var errEmpty = errors.New("no contacts found in the file")

// titleCase is used for vCard N parts.
func titleCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
