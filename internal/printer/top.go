package printer

import (
	"io"
	"strings"
	"unicode"

	"github.com/tolgayakar/xpcolor/internal/config"
)

type TopPrinter struct {
	Theme *config.Theme
}

func (p *TopPrinter) Print(r io.Reader, w io.Writer) {
	var cols []column
	eachLine(r, w, func(line string) string {
		if strings.HasPrefix(line, "TYPE ") && strings.Contains(line, "CPU(cores)") {
			cols = parseHeader(line)
			return p.Theme.Table.Header.Render(line)
		}
		if cols != nil {
			return p.row(line, cols)
		}
		if strings.HasPrefix(line, "No Crossplane pods found") {
			return p.Theme.Status.Unknown.Render(line)
		}
		if key, value, ok := strings.Cut(line, ": "); ok {
			return p.Theme.Data.Key.Render(key) + ": " + p.Theme.Data.Number.Render(value)
		}
		return line
	})
}

func (p *TopPrinter) row(line string, cols []column) string {
	var b strings.Builder
	for i, cell := range splitCells(line, cols) {
		content := strings.TrimRightFunc(cell, unicode.IsSpace)
		switch cols[i].name {
		case "TYPE":
			b.WriteString(p.Theme.Base.Secondary.Render(content))
		case "NAMESPACE":
			b.WriteString(p.Theme.Base.Muted.Render(content))
		case "CPU(cores)", "MEMORY":
			b.WriteString(p.Theme.Data.Number.Render(content))
		default:
			b.WriteString(content)
		}
		b.WriteString(cell[len(content):])
	}
	return b.String()
}
