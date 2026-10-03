package command

import (
	"strconv"
)

const DefaultCrossplaneCmd = "crossplane"

type Options struct {
	ForceColors     bool
	Plain           bool
	LightBackground bool
	ShowVersion     bool
	CrossplaneCmd   string
}

func ParseArgs(args []string, getenv func(string) string) (Options, []string) {
	opts := Options{
		ForceColors:     envBool(getenv, "XPCOLOR_FORCE_COLORS"),
		Plain:           getenv("NO_COLOR") != "",
		LightBackground: envBool(getenv, "XPCOLOR_LIGHT_BACKGROUND"),
		CrossplaneCmd:   DefaultCrossplaneCmd,
	}
	if cmd := getenv("XPCOLOR_CROSSPLANE"); cmd != "" {
		opts.CrossplaneCmd = cmd
	}

	rest := make([]string, 0, len(args))
	for i, arg := range args {
		if arg == "--" {
			rest = append(rest, args[i:]...)
			break
		}
		switch arg {
		case "--force-colors":
			opts.ForceColors = true
		case "--plain":
			opts.Plain = true
		case "--light-background":
			opts.LightBackground = true
		case "--xpcolor-version":
			opts.ShowVersion = true
		default:
			rest = append(rest, arg)
		}
	}
	return opts, rest
}

func envBool(getenv func(string) string, key string) bool {
	b, err := strconv.ParseBool(getenv(key))
	return err == nil && b
}
