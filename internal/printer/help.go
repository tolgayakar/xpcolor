package printer

import (
	"io"
	"regexp"
	"strings"

	"github.com/tolgayakar/xpcolor/internal/config"
)

type HelpPrinter struct {
	Theme *config.Theme
}

var (
	helpSection  = regexp.MustCompile(`^(Arguments|Flags|Commands):$`)
	helpHeading  = regexp.MustCompile(`^(\s*)(#{2,} .*)$`)
	helpFlag     = regexp.MustCompile(`^(\s+)((?:-\w, )?--[\w-]+)(=\S*)?(.*)$`)
	helpArgument = regexp.MustCompile(`^(  )(\[?<[^>]+>\]?(?:\.\.\.)?)(.*)$`)
	helpCommand  = regexp.MustCompile(`^(  )(\S.*?)(\s{2,}.*)$`)
	maturityTag  = regexp.MustCompile(`\[(BETA|ALPHA)\]`)
)

func (p *HelpPrinter) Print(r io.Reader, w io.Writer) {
	section := ""
	eachLine(r, w, func(line string) string {
		if rest, ok := strings.CutPrefix(line, "Usage: "); ok {
			return p.Theme.Table.Header.Render("Usage:") + " " + p.Theme.Base.Key.Render(rest)
		}
		if m := helpSection.FindStringSubmatch(line); m != nil {
			section = m[1]
			return p.Theme.Table.Header.Render(line)
		}
		if m := helpHeading.FindStringSubmatch(line); m != nil {
			return m[1] + p.Theme.Table.Header.Render(m[2])
		}
		if strings.HasPrefix(line, "Run \"") {
			return p.Theme.Base.Muted.Render(line)
		}

		switch section {
		case "Flags":
			if m := helpFlag.FindStringSubmatch(line); m != nil {
				return m[1] + p.Theme.Base.Secondary.Render(m[2]) + p.Theme.Base.Muted.Render(m[3]) + p.tags(m[4])
			}
		case "Arguments":
			if m := helpArgument.FindStringSubmatch(line); m != nil {
				return m[1] + p.Theme.Base.Primary.Render(m[2]) + p.tags(m[3])
			}
		case "Commands":
			if m := helpCommand.FindStringSubmatch(line); m != nil {
				return m[1] + p.Theme.Base.Key.Render(m[2]) + p.tags(m[3])
			}
		}
		return p.tags(line)
	})
}

func (p *HelpPrinter) tags(s string) string {
	return maturityTag.ReplaceAllStringFunc(s, p.Theme.Status.Unknown.Render)
}
