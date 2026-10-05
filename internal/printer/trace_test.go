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
	th, err := config.LoadTheme(false, func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	return th
}

// Golden files store ESC as a literal \e so diffs stay readable.
func TestTracePrinterGolden(t *testing.T) {
	inputs, err := filepath.Glob("testdata/trace/*.txt")
	if err != nil {
		t.Fatal(err)
	}
	for _, in := range inputs {
		t.Run(filepath.Base(in), func(t *testing.T) {
			src, err := os.ReadFile(in)
			if err != nil {
				t.Fatal(err)
			}
			var out bytes.Buffer
			(&TracePrinter{Theme: testTheme(t)}).Print(bytes.NewReader(src), &out)

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

func TestTracePrinterRow(t *testing.T) {
	header := "NAME                       SYNCED   READY   STATUS"
	tests := []struct {
		name string
		line string
		want string
	}{
		{
			name: "error row",
			line: "XBucket/foo                -        -       Error: cannot get resource",
			want: `\e[96mXBucket\e[0m/foo                \e[90;3m-\e[0m        \e[90;3m-\e[0m       \e[31mError\e[0m: cannot get resource`,
		},
		{
			name: "deleting row",
			line: "XBucket/foo                True     True    Deleting",
			want: `\e[96mXBucket\e[0m/foo                \e[32mTrue\e[0m     \e[32mTrue\e[0m    \e[33mDeleting\e[0m`,
		},
		{
			name: "short row",
			line: "XBucket/foo",
			want: `\e[96mXBucket\e[0m/foo`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			(&TracePrinter{Theme: testTheme(t)}).Print(strings.NewReader(header+"\n"+tt.line), &out)
			lines := strings.SplitN(out.String(), "\n", 2)
			if got := strings.ReplaceAll(lines[1], "\033", `\e`); got != tt.want {
				t.Errorf("got  %s\nwant %s", got, tt.want)
			}
		})
	}
}

func TestTracePrinterNoHeader(t *testing.T) {
	in := "crossplane: error: something went wrong\n"
	var out bytes.Buffer
	(&TracePrinter{Theme: testTheme(t)}).Print(strings.NewReader(in), &out)
	if out.String() != in {
		t.Errorf("got %q, want %q", out.String(), in)
	}
}
