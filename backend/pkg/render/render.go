// Package render evaluates message templates. Templates use Go template
// syntax with {{.var}} placeholders; the email channel is rendered with
// html/template so user data cannot inject markup.
package render

import (
	"bytes"
	"errors"
	"fmt"
	htmltemplate "html/template"
	"sort"
	"strings"
	texttemplate "text/template"
	"text/template/parse"

	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// ErrMissingVars is returned when data lacks variables the template uses.
var ErrMissingVars = errors.New("render: missing variables")

// MaxOutput caps rendered size to protect providers from runaway templates.
const MaxOutput = 256 * 1024

// Funcs are the helpers available in templates.
var Funcs = map[string]any{
	"upper": strings.ToUpper,
	"lower": strings.ToLower,
	"trim":  strings.TrimSpace,
	"default": func(def any, v any) any {
		if v == nil || v == "" {
			return def
		}
		return v
	},
}

// Result is a rendered template.
type Result struct {
	Subject string
	Body    string
}

// Validate parses subject and body and returns the variables they use,
// sorted. It reports syntax errors without rendering.
func Validate(subject, body string) ([]string, error) {
	seen := map[string]struct{}{}
	for name, src := range map[string]string{"subject": subject, "body": body} {
		if src == "" {
			continue
		}
		tr, err := texttemplate.New(name).Funcs(Funcs).Parse(src)
		if err != nil {
			return nil, fmt.Errorf("%w: %s: %v", domain.ErrValidation, name, cleanErr(err))
		}
		collectVars(tr.Tree.Root, seen)
	}
	vars := make([]string, 0, len(seen))
	for v := range seen {
		vars = append(vars, v)
	}
	sort.Strings(vars)
	return vars, nil
}

// Render evaluates subject and body with data for the given channel.
// Missing variables are reported together in ErrMissingVars.
func Render(channel domain.Channel, subject, body string, data map[string]any) (Result, error) {
	vars, err := Validate(subject, body)
	if err != nil {
		return Result{}, err
	}
	var missing []string
	for _, v := range vars {
		if _, ok := data[v]; !ok {
			missing = append(missing, v)
		}
	}
	if len(missing) > 0 {
		return Result{}, fmt.Errorf("%w: %s", ErrMissingVars, strings.Join(missing, ", "))
	}
	if data == nil {
		data = map[string]any{}
	}

	exec := textExec
	if channel == domain.ChannelEmail {
		exec = htmlExec
	}
	out := Result{}
	if subject != "" {
		// Subjects are always plain text, even for email.
		if out.Subject, err = textExec(subject, data); err != nil {
			return Result{}, err
		}
	}
	if out.Body, err = exec(body, data); err != nil {
		return Result{}, err
	}
	if len(out.Body) > MaxOutput {
		return Result{}, fmt.Errorf("%w: rendered body exceeds %d bytes", domain.ErrValidation, MaxOutput)
	}
	return out, nil
}

func textExec(src string, data map[string]any) (string, error) {
	t, err := texttemplate.New("t").Funcs(Funcs).Option("missingkey=error").Parse(src)
	if err != nil {
		return "", fmt.Errorf("%w: %v", domain.ErrValidation, cleanErr(err))
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("%w: %v", domain.ErrValidation, cleanErr(err))
	}
	return buf.String(), nil
}

func htmlExec(src string, data map[string]any) (string, error) {
	t, err := htmltemplate.New("t").Funcs(Funcs).Option("missingkey=error").Parse(src)
	if err != nil {
		return "", fmt.Errorf("%w: %v", domain.ErrValidation, cleanErr(err))
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("%w: %v", domain.ErrValidation, cleanErr(err))
	}
	return buf.String(), nil
}

// collectVars walks the parse tree and records top-level .field names.
func collectVars(node parse.Node, out map[string]struct{}) {
	switch n := node.(type) {
	case *parse.ListNode:
		if n == nil {
			return
		}
		for _, c := range n.Nodes {
			collectVars(c, out)
		}
	case *parse.ActionNode:
		collectVars(n.Pipe, out)
	case *parse.PipeNode:
		if n == nil {
			return
		}
		for _, c := range n.Cmds {
			collectVars(c, out)
		}
	case *parse.CommandNode:
		for _, a := range n.Args {
			collectVars(a, out)
		}
	case *parse.FieldNode:
		if len(n.Ident) > 0 {
			out[n.Ident[0]] = struct{}{}
		}
	case *parse.IfNode:
		collectBranch(&n.BranchNode, out)
	case *parse.RangeNode:
		collectBranch(&n.BranchNode, out)
	case *parse.WithNode:
		collectBranch(&n.BranchNode, out)
	}
}

func collectBranch(b *parse.BranchNode, out map[string]struct{}) {
	collectVars(b.Pipe, out)
	collectVars(b.List, out)
	if b.ElseList != nil {
		collectVars(b.ElseList, out)
	}
}

// cleanErr strips the "template: t:1:" prefix Go adds.
func cleanErr(err error) string {
	s := err.Error()
	if i := strings.Index(s, ": "); i >= 0 && strings.HasPrefix(s, "template:") {
		s = s[i+2:]
		if j := strings.Index(s, ": "); j >= 0 && strings.Contains(s[:j], ":") {
			s = s[j+2:]
		}
	}
	return s
}
