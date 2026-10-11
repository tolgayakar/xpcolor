package printer

import (
	"io"
	"regexp"
	"strings"

	"github.com/tolgayakar/xpcolor/internal/color"
	"github.com/tolgayakar/xpcolor/internal/config"
)

type ValidatePrinter struct {
	Theme *config.Theme
}

var validateTotal = regexp.MustCompile(`^Total (\d+) resources: (\d+) missing schemas, (\d+) success cases, (\d+) failure cases$`)

func (p *ValidatePrinter) Print(r io.Reader, w io.Writer) {
	eachLine(r, w, p.line)
}

func (p *ValidatePrinter) line(line string) string {
	markers := []struct {
		prefix string
		color  color.Color
	}{
		{"[✓]", p.Theme.Status.True},
		{"[!]", p.Theme.Status.Unknown},
		{"[x]", p.Theme.Status.False},
	}
	for _, m := range markers {
		if rest, ok := strings.CutPrefix(line, m.prefix); ok {
			return m.color.Render(m.prefix) + rest
		}
	}

	if m := validateTotal.FindStringSubmatch(line); m != nil {
		count := func(n string, c color.Color) string {
			if n == "0" {
				return n
			}
			return c.Render(n)
		}
		return "Total " + p.Theme.Data.Number.Render(m[1]) + " resources: " +
			count(m[2], p.Theme.Status.Unknown) + " missing schemas, " +
			count(m[3], p.Theme.Status.True) + " success cases, " +
			count(m[4], p.Theme.Status.False) + " failure cases"
	}

	if strings.HasPrefix(line, "schemas does not exist, downloading:") || strings.HasPrefix(line, "updating cached schemas:") {
		return p.Theme.Base.Muted.Render(line)
	}
	return line
}
