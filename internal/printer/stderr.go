package printer

import (
	"io"
	"strings"

	"github.com/tolgayakar/xpcolor/internal/config"
)

type StderrPrinter struct {
	Theme *config.Theme
}

func (p *StderrPrinter) Print(r io.Reader, w io.Writer) {
	eachLine(r, w, func(line string) string {
		if rest, ok := strings.CutPrefix(line, "crossplane: error:"); ok {
			return p.Theme.Status.False.Render("crossplane: error:") + rest
		}
		return line
	})
}
