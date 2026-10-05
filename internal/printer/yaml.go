package printer

import (
	"io"
	"strconv"
	"strings"

	"github.com/tolgayakar/xpcolor/internal/color"
	"github.com/tolgayakar/xpcolor/internal/config"
)

type YAMLPrinter struct {
	Theme *config.Theme

	blockIndent int
	inQuote     byte
}

func (p *YAMLPrinter) Print(r io.Reader, w io.Writer) {
	p.blockIndent = -1
	eachLine(r, w, p.line)
}

func (p *YAMLPrinter) line(line string) string {
	trimmed := strings.TrimLeft(line, " ")
	indent := line[:len(line)-len(trimmed)]

	if p.blockIndent >= 0 {
		if trimmed == "" || len(indent) > p.blockIndent {
			return indent + p.Theme.Data.String.Render(trimmed)
		}
		p.blockIndent = -1
	}
	if p.inQuote != 0 {
		if closesQuote(trimmed, p.inQuote) {
			p.inQuote = 0
		}
		return indent + p.Theme.Data.String.Render(trimmed)
	}

	switch {
	case trimmed == "---" || trimmed == "..." || strings.HasPrefix(trimmed, "#"):
		return indent + p.Theme.Base.Muted.Render(trimmed)
	case trimmed == "":
		return line
	}

	var b strings.Builder
	b.WriteString(indent)
	col := len(indent)
	for strings.HasPrefix(trimmed, "- ") || trimmed == "-" {
		b.WriteString(trimmed[:min(2, len(trimmed))])
		trimmed = trimmed[min(2, len(trimmed)):]
		col += 2
	}

	key, value, isKey := cutKey(trimmed)
	if !isKey {
		b.WriteString(p.value(trimmed, col))
		return b.String()
	}
	b.WriteString(p.Theme.Data.Key.Render(key))
	b.WriteString(":")
	if value != "" {
		spaces := len(value) - len(strings.TrimLeft(value, " "))
		b.WriteString(value[:spaces])
		b.WriteString(p.value(value[spaces:], col))
	}
	return b.String()
}

func (p *YAMLPrinter) value(v string, col int) string {
	if v == "" {
		return v
	}
	if isBlockScalar(v) {
		p.blockIndent = col
		return v
	}
	if q := v[0]; (q == '"' || q == '\'') && !closesQuote(v[1:], q) {
		p.inQuote = q
		return p.Theme.Data.String.Render(v)
	}
	switch v {
	case "{}", "[]", "|", ">":
		return v
	}
	return valueColor(v, p.Theme).Render(v)
}

// cutKey splits "key: value" or "key:"; keys may be quoted.
func cutKey(s string) (key, value string, ok bool) {
	start := 0
	if s != "" && (s[0] == '"' || s[0] == '\'') {
		end := strings.IndexByte(s[1:], s[0])
		if end < 0 {
			return "", "", false
		}
		start = end + 2
	}
	i := strings.Index(s[start:], ": ")
	if i >= 0 {
		return s[:start+i], s[start+i+1:], true
	}
	if strings.HasSuffix(s, ":") && !strings.ContainsAny(s[start:len(s)-1], " ") {
		return s[:len(s)-1], "", true
	}
	return "", "", false
}

func isBlockScalar(v string) bool {
	if v[0] != '|' && v[0] != '>' {
		return false
	}
	return strings.Trim(v[1:], "+-0123456789") == ""
}

func closesQuote(s string, q byte) bool {
	if q == '\'' {
		return strings.HasSuffix(s, "'") && !strings.HasSuffix(s, "''")
	}
	return strings.HasSuffix(s, `"`) && !strings.HasSuffix(s, `\"`)
}

func valueColor(v string, theme *config.Theme) color.Color {
	switch strings.Trim(v, `"'`) {
	case "true", "True":
		return theme.Data.True
	case "false", "False":
		return theme.Data.False
	}
	switch v {
	case "null", "~":
		return theme.Data.Null
	}
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return theme.Data.Number
	}
	return theme.Data.String
}
