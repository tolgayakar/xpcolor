package printer

import (
	"bytes"
	"strings"
	"testing"
)

func TestTracePrinterGolden(t *testing.T) {
	runGolden(t, "trace", func() Printer { return &TracePrinter{Theme: testTheme(t)} })
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
