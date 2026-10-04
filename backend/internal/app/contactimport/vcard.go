package contactimport

import (
	"bufio"
	"io"
	"mime/quotedprintable"
	"strings"

	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
)

// ParseVCard reads vCard 2.1 / 3.0 / 4.0 files (Google, Apple, Outlook
// exports): FN / N, TEL, EMAIL, X-ABLabel-free. Multiple cards per file.
func ParseVCard(r io.Reader) ([]contacts.Input, []RowError, error) {
	cards, err := readVCards(r)
	if err != nil {
		return nil, nil, err
	}
	var out []contacts.Input
	for _, c := range cards {
		if in, ok := c.toInput(); ok {
			out = append(out, in)
		}
		if len(out) >= MaxRows {
			break
		}
	}
	if len(out) == 0 {
		return nil, nil, errEmpty
	}
	return out, nil, nil
}

type vcard struct {
	props []vprop
}

type vprop struct {
	name   string
	params map[string]string
	value  string
}

func (c vcard) first(name string) string {
	for _, p := range c.props {
		if p.name == name {
			return p.value
		}
	}
	return ""
}

func (c vcard) all(name string) []string {
	var out []string
	for _, p := range c.props {
		if p.name == name && strings.TrimSpace(p.value) != "" {
			out = append(out, strings.TrimSpace(p.value))
		}
	}
	return out
}

func (c vcard) toInput() (contacts.Input, bool) {
	var in contacts.Input
	in.Name = strings.TrimSpace(c.first("FN"))
	if in.Name == "" {
		parts := strings.Split(c.first("N"), ";")
		var name []string
		for _, i := range []int{3, 1, 2, 0} { // prefix, given, additional, family
			if i < len(parts) && strings.TrimSpace(parts[i]) != "" {
				name = append(name, titleCase(parts[i]))
			}
		}
		in.Name = strings.Join(name, " ")
	}
	for _, tel := range c.all("TEL") {
		tel = strings.TrimPrefix(tel, "tel:")
		if !isPhone(tel) {
			continue
		}
		if in.Phone == "" {
			in.Phone = cleanPhone(tel)
		} else if in.WhatsApp == "" && cleanPhone(tel) != in.Phone {
			in.WhatsApp = cleanPhone(tel)
		}
	}
	for _, mail := range c.all("EMAIL") {
		mail = strings.TrimPrefix(strings.ToLower(mail), "mailto:")
		if emailRe.MatchString(mail) {
			in.Email = mail
			break
		}
	}
	if v := c.first("X-TELEGRAM"); v != "" {
		in.TelegramChatID = v
	}
	if v := c.first("UID"); v != "" && len(v) <= 128 {
		in.ExternalID = "vcard:" + v
	}
	for _, cat := range c.all("CATEGORIES") {
		for _, t := range strings.Split(cat, ",") {
			if t = strings.TrimSpace(t); t != "" {
				in.Tags = append(in.Tags, t)
			}
		}
	}
	return in, hasReach(in)
}

func readVCards(r io.Reader) ([]vcard, error) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	var lines []string
	for sc.Scan() {
		line := strings.TrimRight(sc.Text(), "\r")
		if (strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t")) && len(lines) > 0 {
			lines[len(lines)-1] += line[1:] // unfold
			continue
		}
		lines = append(lines, line)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	var cards []vcard
	var cur *vcard
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		upper := strings.ToUpper(strings.TrimSpace(line))
		switch {
		case upper == "BEGIN:VCARD":
			cur = &vcard{}
			continue
		case upper == "END:VCARD":
			if cur != nil {
				cards = append(cards, *cur)
			}
			cur = nil
			continue
		case cur == nil || line == "":
			continue
		}
		p, ok := parseProp(line)
		if !ok {
			continue
		}
		// vCard 2.1 quoted-printable soft line breaks end with "=".
		for p.params["ENCODING"] == "QUOTED-PRINTABLE" && strings.HasSuffix(p.value, "=") && i+1 < len(lines) {
			i++
			p.value = strings.TrimSuffix(p.value, "=") + lines[i]
		}
		if p.params["ENCODING"] == "QUOTED-PRINTABLE" {
			if b, err := io.ReadAll(quotedprintable.NewReader(strings.NewReader(p.value))); err == nil {
				p.value = string(b)
			}
		}
		cur.props = append(cur.props, p)
	}
	return cards, nil
}

func parseProp(line string) (vprop, bool) {
	colon := strings.Index(line, ":")
	if colon <= 0 {
		return vprop{}, false
	}
	head, value := line[:colon], line[colon+1:]
	segs := strings.Split(head, ";")
	name := strings.ToUpper(segs[0])
	if dot := strings.LastIndex(name, "."); dot >= 0 { // item1.TEL (Apple groups)
		name = name[dot+1:]
	}
	params := map[string]string{}
	for _, s := range segs[1:] {
		k, v, found := strings.Cut(s, "=")
		if !found {
			// vCard 2.1 bare params like "CELL" or "QUOTED-PRINTABLE".
			if strings.EqualFold(s, "QUOTED-PRINTABLE") {
				params["ENCODING"] = "QUOTED-PRINTABLE"
			} else {
				params["TYPE"] = strings.ToUpper(s)
			}
			continue
		}
		params[strings.ToUpper(k)] = strings.ToUpper(strings.Trim(v, "\""))
	}
	value = strings.NewReplacer(`\,`, ",", `\;`, ";", `\n`, "\n", `\N`, "\n").Replace(value)
	return vprop{name: name, params: params, value: value}, true
}
