package printer

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/tolgayakar/xpcolor/internal/config"
)

var update = flag.Bool("update", false, "update golden files")

var ansi = regexp.MustCompile("\033\\[[0-9;]*m")

func testTheme(t *testing.T) *config.Theme {
	t.Helper()
	th, err := config.LoadTheme("", func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	return th
}

// Golden files store ESC as a literal \e so diffs stay readable.
func runGolden(t *testing.T, dir string, newPrinter func() Printer) {
	t.Helper()
	inputs, err := filepath.Glob(filepath.Join("testdata", dir, "*.txt"))
	if err != nil || len(inputs) == 0 {
		t.Fatalf("no inputs in testdata/%s: %v", dir, err)
	}
	for _, in := range inputs {
		t.Run(filepath.Base(in), func(t *testing.T) {
			src, err := os.ReadFile(in)
			if err != nil {
				t.Fatal(err)
			}
			var out bytes.Buffer
			newPrinter().Print(bytes.NewReader(src), &out)

			if got := ansi.ReplaceAllString(out.String(), ""); got != string(src) {
				t.Fatalf("output without colors differs from input:\n%s", got)
			}

			got := strings.ReplaceAll(out.String(), "\033", `\e`)
			golden := strings.TrimSuffix(in, ".txt") + ".golden"
			if *update {
				if err := os.WriteFile(golden, []byte(got), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			want, err := os.ReadFile(golden)
			if err != nil {
				t.Fatal(err)
			}
			if got != string(want) {
				t.Errorf("output differs from %s (run with -update to regenerate):\n%s", golden, got)
			}
		})
	}
}
