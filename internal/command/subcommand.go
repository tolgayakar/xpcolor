package command

import (
	"slices"
	"strings"

	"github.com/tolgayakar/xpcolor/internal/config"
	"github.com/tolgayakar/xpcolor/internal/printer"
)

type subcommand struct {
	path   []string
	output string
	format string
	watch  bool
	help   bool
}

var valueFlags = []string{
	"--config", "-c", "--context", "-n", "--namespace", "-o", "--output",
	"--concurrency", "--show-package-dependencies", "--show-package-revisions",
	"--as", "--as-group", "--as-uid",
}

func parseSubcommand(args []string) subcommand {
	var sc subcommand
	for i := 0; i < len(args); i++ {
		arg := args[i]
		name, value, hasValue := strings.Cut(arg, "=")
		switch {
		case arg == "--":
			return sc
		case arg == "-h" || arg == "--help":
			sc.help = true
		case name == "-w" || name == "--watch":
			sc.watch = !hasValue || value == "true"
		case name == "-o" || name == "--output":
			if !hasValue && i+1 < len(args) {
				i++
				value = args[i]
			}
			sc.output = value
		case strings.HasPrefix(arg, "-o") && !strings.HasPrefix(arg, "--"):
			sc.output = arg[2:]
		case name == "--format":
			if !hasValue && i+1 < len(args) {
				i++
				value = args[i]
			}
			sc.format = value
		case slices.Contains(valueFlags, name):
			if !hasValue {
				i++
			}
		case strings.HasPrefix(arg, "-"):
		default:
			sc.path = append(sc.path, arg)
		}
	}
	return sc
}

func (sc subcommand) is(path ...string) bool {
	return len(sc.path) >= len(path) && slices.Equal(sc.path[:len(path)], path)
}

// trace --watch renders a full-screen TUI, so it is passed through untouched.
func printerFor(args []string, theme *config.Theme) printer.Printer {
	sc := parseSubcommand(args)
	if sc.help {
		return nil
	}
	switch {
	case sc.is("resource", "trace") || sc.is("beta", "trace"):
		if sc.watch {
			return nil
		}
		switch sc.output {
		case "", "default", "wide":
			return &printer.TracePrinter{Theme: theme}
		case "yaml":
			return &printer.YAMLPrinter{Theme: theme}
		case "json":
			return &printer.JSONPrinter{Theme: theme}
		}
	case sc.is("resource", "validate") || sc.is("beta", "validate"):
		switch sc.output {
		case "", "text":
			return &printer.ValidatePrinter{Theme: theme}
		case "yaml":
			return &printer.YAMLPrinter{Theme: theme}
		case "json":
			return &printer.JSONPrinter{Theme: theme}
		}
	case sc.is("xrd", "convert"):
		if sc.format == "jsonschema" {
			return &printer.JSONPrinter{Theme: theme}
		}
		return &printer.YAMLPrinter{Theme: theme}
	case sc.is("render"), sc.is("composition", "render"), sc.is("operation", "render"),
		sc.is("composition", "convert"), sc.is("xr", "generate"), sc.is("xr", "patch"),
		sc.is("config", "view"):
		return &printer.YAMLPrinter{Theme: theme}
	}
	return nil
}
