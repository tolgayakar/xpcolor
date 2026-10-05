package command

import (
	"strings"
	"testing"

	"github.com/tolgayakar/xpcolor/internal/config"
)

func TestPrinterFor(t *testing.T) {
	tests := []struct {
		args []string
		want bool
	}{
		{args: []string{"resource", "trace", "xr/foo"}, want: true},
		{args: []string{"beta", "trace", "xr/foo"}, want: true},
		{args: []string{"resource", "trace", "xr/foo", "-o", "wide"}, want: true},
		{args: []string{"resource", "trace", "xr/foo", "--output=default"}, want: true},
		{args: []string{"--config", "/tmp/cfg", "resource", "trace", "xr/foo"}, want: true},
		{args: []string{"resource", "trace", "-n", "ns", "xr", "foo"}, want: true},
		{args: []string{"resource", "trace", "xr/foo", "-o", "yaml"}, want: false},
		{args: []string{"resource", "trace", "xr/foo", "-ojson"}, want: false},
		{args: []string{"resource", "trace", "xr/foo", "--output=dot"}, want: false},
		{args: []string{"resource", "trace", "xr/foo", "-w"}, want: false},
		{args: []string{"resource", "trace", "xr/foo", "--watch"}, want: false},
		{args: []string{"resource", "trace", "--help"}, want: false},
		{args: []string{"resource", "--", "trace"}, want: false},
		{args: []string{"version"}, want: false},
		{args: []string{}, want: false},
	}
	theme, err := config.LoadTheme(false, func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			if got := printerFor(tt.args, theme) != nil; got != tt.want {
				t.Errorf("printerFor(%q) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}
