package command

import (
	"strings"
	"testing"

	"github.com/tolgayakar/xpcolor/internal/config"
	"github.com/tolgayakar/xpcolor/internal/printer"
)

func TestPrinterFor(t *testing.T) {
	tests := []struct {
		args []string
		want string
	}{
		{args: []string{"resource", "trace", "xr/foo"}, want: "trace"},
		{args: []string{"beta", "trace", "xr/foo"}, want: "trace"},
		{args: []string{"resource", "trace", "xr/foo", "-o", "wide"}, want: "trace"},
		{args: []string{"resource", "trace", "xr/foo", "--output=default"}, want: "trace"},
		{args: []string{"--config", "/tmp/cfg", "resource", "trace", "xr/foo"}, want: "trace"},
		{args: []string{"resource", "trace", "-n", "ns", "xr", "foo"}, want: "trace"},
		{args: []string{"resource", "trace", "xr/foo", "-o", "yaml"}, want: "yaml"},
		{args: []string{"resource", "trace", "xr/foo", "-ojson"}, want: "json"},
		{args: []string{"resource", "trace", "xr/foo", "--output=dot"}, want: ""},
		{args: []string{"resource", "trace", "xr/foo", "-w"}, want: ""},
		{args: []string{"resource", "trace", "xr/foo", "--watch"}, want: ""},
		{args: []string{"resource", "trace", "--help"}, want: "help"},
		{args: []string{"resource", "--", "trace"}, want: "help"},
		{args: []string{"composition", "render", "xr.yaml", "comp.yaml", "fn.yaml"}, want: "yaml"},
		{args: []string{"render", "xr.yaml", "comp.yaml", "fn.yaml"}, want: "yaml"},
		{args: []string{"composition", "render", "-o", "observed.yaml", "xr.yaml"}, want: "yaml"},
		{args: []string{"operation", "render", "op.yaml", "fn.yaml"}, want: "yaml"},
		{args: []string{"composition", "convert", "composition-environment", "c.yaml"}, want: "yaml"},
		{args: []string{"xrd", "convert", "xrd.yaml"}, want: "yaml"},
		{args: []string{"xrd", "convert", "xrd.yaml", "--format", "jsonschema"}, want: "json"},
		{args: []string{"xr", "generate", "claim.yaml"}, want: "yaml"},
		{args: []string{"xr", "patch", "xr.yaml"}, want: "yaml"},
		{args: []string{"config", "view"}, want: "yaml"},
		{args: []string{"composition", "render", "--help"}, want: "help"},
		{args: []string{"-h"}, want: "help"},
		{args: []string{"resource"}, want: "help"},
		{args: []string{"--config", "/tmp/cfg", "xrd"}, want: "help"},
		{args: []string{"resource", "validate", "ext.yaml", "res.yaml"}, want: "validate"},
		{args: []string{"beta", "validate", "ext.yaml", "res.yaml", "-o", "text"}, want: "validate"},
		{args: []string{"resource", "validate", "ext.yaml", "res.yaml", "--output=json"}, want: "json"},
		{args: []string{"resource", "validate", "ext.yaml", "res.yaml", "-o", "yaml"}, want: "yaml"},
		{args: []string{"cluster", "top"}, want: "top"},
		{args: []string{"cluster", "top", "-s", "-n", "crossplane-system"}, want: "top"},
		{args: []string{"beta", "top"}, want: "top"},
		{args: []string{"version"}, want: ""},
		{args: []string{}, want: "help"},
		{args: []string{"xpkg", "login"}, want: ""},
	}
	theme, err := config.LoadTheme("", func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			if got := printerName(printerFor(tt.args, theme)); got != tt.want {
				t.Errorf("printerFor(%q) = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}

func printerName(p printer.Printer) string {
	switch p.(type) {
	case *printer.TracePrinter:
		return "trace"
	case *printer.YAMLPrinter:
		return "yaml"
	case *printer.JSONPrinter:
		return "json"
	case *printer.ValidatePrinter:
		return "validate"
	case *printer.TopPrinter:
		return "top"
	case *printer.HelpPrinter:
		return "help"
	}
	return ""
}
