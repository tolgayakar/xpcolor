package color

import (
	"fmt"
	"strconv"
	"strings"
)

const reset = "\033[0m"

type Color struct {
	source string
	code   string
}

func Parse(s string) (Color, error) {
	var codes []string
	for field := range strings.SplitSeq(s, ":") {
		field = strings.ToLower(strings.TrimSpace(field))
		if field == "" || field == "none" {
			continue
		}
		code, err := parseField(field)
		if err != nil {
			return Color{}, err
		}
		codes = append(codes, code)
	}
	return Color{source: s, code: strings.Join(codes, ";")}, nil
}

func MustParse(s string) Color {
	c, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return c
}

func (c Color) Render(s string) string {
	if c.code == "" || s == "" {
		return s
	}
	return "\033[" + c.code + "m" + s + reset
}

func (c Color) String() string {
	if c.code == "" {
		return "none"
	}
	return c.source
}

var attrs = map[string]string{
	"bold":      "1",
	"b":         "1",
	"dim":       "2",
	"italic":    "3",
	"i":         "3",
	"underline": "4",
	"u":         "4",
	"reverse":   "7",
	"strike":    "9",
}

var names = map[string]int{
	"black":   0,
	"red":     1,
	"green":   2,
	"yellow":  3,
	"blue":    4,
	"magenta": 5,
	"cyan":    6,
	"white":   7,
}

func parseField(field string) (string, error) {
	if code, ok := attrs[field]; ok {
		return code, nil
	}
	if v, ok := strings.CutPrefix(field, "fg="); ok {
		return parseColor(v, false)
	}
	if v, ok := strings.CutPrefix(field, "bg="); ok {
		return parseColor(v, true)
	}
	return parseColor(field, false)
}

func parseColor(s string, bg bool) (string, error) {
	base := 30
	if bg {
		base = 40
	}
	switch {
	case s == "gray" || s == "grey":
		return strconv.Itoa(base + 60), nil
	case strings.HasPrefix(s, "hi"):
		if n, ok := names[s[2:]]; ok {
			return strconv.Itoa(base + 60 + n), nil
		}
	case strings.HasPrefix(s, "#") && len(s) == 7:
		v, err := strconv.ParseUint(s[1:], 16, 32)
		if err == nil {
			return fmt.Sprintf("%d;2;%d;%d;%d", base+8, v>>16, v>>8&0xff, v&0xff), nil
		}
	}
	if n, ok := names[s]; ok {
		return strconv.Itoa(base + n), nil
	}
	if n, err := strconv.ParseUint(s, 10, 8); err == nil {
		return fmt.Sprintf("%d;5;%d", base+8, n), nil
	}
	return "", fmt.Errorf("invalid color %q", s)
}
