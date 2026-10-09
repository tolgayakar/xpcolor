package command

import (
	"strconv"
	"strings"
)

const DefaultCrossplaneCmd = "crossplane"

type Options struct {
	ForceColors     bool
	Plain           bool
	LightBackground bool
	ShowVersion     bool
	Theme           string
	CrossplaneCmd   string
}

func (o Options) Preset() string {
	if o.Theme == "" && o.LightBackground {
		return "light"
	}
	return o.Theme
}

func ParseArgs(args []string, getenv func(string) string) (Options, []string) {
	opts := Options{
		ForceColors:     envBool(getenv, "XPCOLOR_FORCE_COLORS"),
		Plain:           getenv("NO_COLOR") != "",
		LightBackground: envBool(getenv, "XPCOLOR_LIGHT_BACKGROUND"),
		Theme:           getenv("XPCOLOR_PRESET"),
		CrossplaneCmd:   DefaultCrossplaneCmd,
	}
	if cmd := getenv("XPCOLOR_CROSSPLANE"); cmd != "" {
		opts.CrossplaneCmd = cmd
	}

	rest := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			rest = append(rest, args[i:]...)
			break
		}
		if v, ok := strings.CutPrefix(arg, "--xpcolor-theme="); ok {
			opts.Theme = v
			continue
		}
		switch arg {
		case "--xpcolor-theme":
			if i+1 < len(args) {
				i++
				opts.Theme = args[i]
			}
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
