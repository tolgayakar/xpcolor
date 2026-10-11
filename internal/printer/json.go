package printer

import (
	"io"
	"strings"

	"github.com/tolgayakar/xpcolor/internal/config"
)

type JSONPrinter struct {
	Theme *config.Theme
}

func (p *JSONPrinter) Print(r io.Reader, w io.Writer) {
	eachLine(r, w, p.line)
}

func (p *JSONPrinter) line(line string) string {
	var b strings.Builder
	for i := 0; i < len(line); {
		c := line[i]
		switch {
		case c == '"':
			end := stringEnd(line, i)
			s := line[i:end]
			if strings.HasPrefix(strings.TrimLeft(line[end:], " "), ":") {
				b.WriteString(p.Theme.Data.Key.Render(s))
			} else {
				b.WriteString(valueColor(s, p.Theme).Render(s))
			}
			i = end
		case strings.ContainsRune("{}[]:, \t", rune(c)):
			b.WriteByte(c)
			i++
		default:
			end := i
			for end < len(line) && !strings.ContainsRune("{}[]:, \t\"", rune(line[end])) {
				end++
			}
			word := line[i:end]
			b.WriteString(valueColor(word, p.Theme).Render(word))
			i = end
		}
	}
	return b.String()
}

// stringEnd returns the index just past the closing quote of the string at i.
func stringEnd(s string, i int) int {
	for j := i + 1; j < len(s); j++ {
		switch s[j] {
		case '\\':
			j++
		case '"':
			return j + 1
		}
	}
	return len(s)
}
