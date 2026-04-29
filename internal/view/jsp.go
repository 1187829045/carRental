package view

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RenderJSPFile reads an existing .jsp view file and renders it as plain HTML by:
// - stripping JSP directives (`<%@ ... %>`) and JSP comments (`<%-- ... --%>`)
// - doing a minimal `${...}` substitution against ctx (supports dot-notation like `user.realname`)
//
// This keeps the original JSP files largely "as-is" while moving execution to Go.
func RenderJSPFile(viewRoot, viewPath string, ctx map[string]any) ([]byte, error) {
	full := filepath.Join(viewRoot, filepath.FromSlash(viewPath)+".jsp")
	raw, err := os.ReadFile(full)
	if err != nil {
		return nil, fmt.Errorf("read view %q: %w", full, err)
	}

	s := string(raw)
	s = stripJSPDirectives(s)
	s = stripJSPComments(s)
	s = substituteEL(s, ctx)
	return []byte(s), nil
}

func stripJSPDirectives(s string) string {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, ln := range lines {
		t := strings.TrimSpace(ln)
		if strings.HasPrefix(t, "<%@") && strings.Contains(t, "%>") {
			continue
		}
		out = append(out, ln)
	}
	return strings.Join(out, "\n")
}

func stripJSPComments(s string) string {
	// Remove all occurrences of <%-- ... --%> (can span lines).
	for {
		start := strings.Index(s, "<%--")
		if start < 0 {
			return s
		}
		end := strings.Index(s[start+4:], "--%>")
		if end < 0 {
			// Unclosed comment; drop the rest.
			return s[:start]
		}
		end = start + 4 + end + len("--%>")
		s = s[:start] + s[end:]
	}
}

func substituteEL(s string, ctx map[string]any) string {
	var b bytes.Buffer
	for {
		i := strings.Index(s, "${")
		if i < 0 {
			b.WriteString(s)
			break
		}
		b.WriteString(s[:i])
		s = s[i+2:]
		j := strings.IndexByte(s, '}')
		if j < 0 {
			// No closing brace; keep the rest as-is.
			b.WriteString("${")
			b.WriteString(s)
			break
		}
		expr := strings.TrimSpace(s[:j])
		s = s[j+1:]

		if v, ok := resolveExpr(ctx, expr); ok {
			b.WriteString(v)
		}
	}
	return b.String()
}

func resolveExpr(ctx map[string]any, expr string) (string, bool) {
	if expr == "" {
		return "", false
	}
	parts := strings.Split(expr, ".")
	var cur any = ctx
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return "", false
		}
		m, ok := cur.(map[string]any)
		if !ok {
			return "", false
		}
		cur, ok = m[p]
		if !ok {
			return "", false
		}
	}
	switch v := cur.(type) {
	case nil:
		return "", true
	case string:
		return v, true
	case fmt.Stringer:
		return v.String(), true
	default:
		return fmt.Sprintf("%v", v), true
	}
}

