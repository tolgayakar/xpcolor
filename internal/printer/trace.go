package printer

import (
	"io"
	"strings"
	"unicode"

	"github.com/tolgayakar/xpcolor/internal/color"
	"github.com/tolgayakar/xpcolor/internal/config"
)

type TracePrinter struct {
	Theme *config.Theme
}

type column struct {
	name  string
	start int
}

func (p *TracePrinter) Print(r io.Reader, w io.Writer) {
	var cols []column
	eachLine(r, w, func(line string) string {
		if strings.HasPrefix(line, "NAME ") && strings.Contains(line, "STATUS") {
			cols = parseHeader(line)
			return p.Theme.Table.Header.Render(line)
		}
		if cols == nil {
			return line
		}
		return p.row(line, cols)
	})
}

func parseHeader(line string) []column {
	var cols []column
	for i, r := range line {
		if r != ' ' && (i == 0 || line[i-1] == ' ') {
			end := strings.IndexByte(line[i:], ' ')
			if end < 0 {
				end = len(line) - i
			}
			cols = append(cols, column{name: line[i : i+end], start: i})
		}
	}
	return cols
}

// Columns are aligned by text/tabwriter, which pads by rune count.
func splitCells(line string, cols []column) []string {
	runes := []rune(line)
	cells := make([]string, len(cols))
	for i, c := range cols {
		start := min(c.start, len(runes))
		end := len(runes)
		if i+1 < len(cols) {
			end = min(cols[i+1].start, len(runes))
		}
		cells[i] = string(runes[start:end])
	}
	return cells
}

func (p *TracePrinter) row(line string, cols []column) string {
	cells := splitCells(line, cols)
	health := rowHealth(cells, cols)

	var b strings.Builder
	for i, cell := range cells {
		content := strings.TrimRightFunc(cell, unicode.IsSpace)
		b.WriteString(p.cell(cols[i].name, content, health))
		b.WriteString(cell[len(content):])
	}
	return b.String()
}

func (p *TracePrinter) cell(name, s, health string) string {
	switch name {
	case "NAME":
		return p.name(s)
	case "SYNCED", "READY", "INSTALLED", "HEALTHY":
		return p.condition(s).Render(s)
	case "STATE":
		if s == "Active" {
			return p.Theme.Status.True.Render(s)
		}
		return p.Theme.Base.Muted.Render(s)
	case "VERSION":
		return p.Theme.Base.Secondary.Render(s)
	case "STATUS":
		reason, msg, ok := strings.Cut(s, ": ")
		reason = p.condition(health).Render(reason)
		if !ok {
			return reason
		}
		return reason + ": " + msg
	}
	return s
}

func (p *TracePrinter) name(s string) string {
	nameStart := strings.IndexFunc(s, func(r rune) bool { return !strings.ContainsRune(" │├└─", r) })
	if nameStart < 0 {
		return p.Theme.Trace.Tree.Render(s)
	}
	tree, rest := s[:nameStart], s[nameStart:]

	ns := ""
	if i := strings.LastIndex(rest, " ("); i >= 0 && strings.HasSuffix(rest, ")") {
		rest, ns = rest[:i], rest[i:]
	}
	kind, name, _ := strings.Cut(rest, "/")

	out := p.Theme.Trace.Tree.Render(tree) + p.Theme.Trace.Kind.Render(kind)
	if rest != kind {
		out += "/" + name
	}
	return out + p.Theme.Trace.Namespace.Render(ns)
}

func (p *TracePrinter) condition(s string) color.Color {
	switch s {
	case "True":
		return p.Theme.Status.True
	case "False":
		return p.Theme.Status.False
	case "Unknown":
		return p.Theme.Status.Unknown
	}
	return p.Theme.Base.Muted
}

func rowHealth(cells []string, cols []column) string {
	health := "True"
	for i, c := range cols {
		v := strings.TrimSpace(cells[i])
		switch c.name {
		case "SYNCED", "READY", "INSTALLED", "HEALTHY":
			if v == "False" {
				return "False"
			}
			if v == "Unknown" {
				health = "Unknown"
			}
		case "STATUS":
			if strings.HasPrefix(v, "Error") {
				return "False"
			}
			if v == "Deleting" {
				health = "Unknown"
			}
		}
	}
	return health
}
