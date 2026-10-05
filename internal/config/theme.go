package config

import (
	"fmt"
	"strings"

	"github.com/tolgayakar/xpcolor/internal/color"
)

type Theme struct {
	Base ThemeBase
	Data ThemeData
}

type ThemeBase struct {
	Primary   color.Color
	Secondary color.Color
	Success   color.Color
	Warning   color.Color
	Danger    color.Color
	Info      color.Color
	Muted     color.Color
	Key       color.Color
}

type ThemeData struct {
	Key    color.Color
	String color.Color
	Number color.Color
	True   color.Color
	False  color.Color
	Null   color.Color
}

func darkTheme() *Theme {
	return &Theme{
		Base: ThemeBase{
			Primary:   color.MustParse("magenta"),
			Secondary: color.MustParse("cyan"),
			Success:   color.MustParse("green"),
			Warning:   color.MustParse("yellow"),
			Danger:    color.MustParse("red"),
			Info:      color.MustParse("white"),
			Muted:     color.MustParse("gray:italic"),
			Key:       color.MustParse("hicyan"),
		},
		Data: ThemeData{
			String: color.MustParse("hiyellow"),
		},
	}
}

func lightTheme() *Theme {
	return &Theme{
		Base: ThemeBase{
			Primary:   color.MustParse("magenta"),
			Secondary: color.MustParse("blue"),
			Success:   color.MustParse("green"),
			Warning:   color.MustParse("yellow"),
			Danger:    color.MustParse("red"),
			Info:      color.MustParse("black"),
			Muted:     color.MustParse("gray:italic"),
			Key:       color.MustParse("hiblue"),
		},
		Data: ThemeData{
			String: color.MustParse("yellow"),
		},
	}
}

type field struct {
	key  string
	ptr  *color.Color
	from string
}

// Base fields must come first so data fields can default from them.
func (t *Theme) fields() []field {
	return []field{
		{key: "base.primary", ptr: &t.Base.Primary},
		{key: "base.secondary", ptr: &t.Base.Secondary},
		{key: "base.success", ptr: &t.Base.Success},
		{key: "base.warning", ptr: &t.Base.Warning},
		{key: "base.danger", ptr: &t.Base.Danger},
		{key: "base.info", ptr: &t.Base.Info},
		{key: "base.muted", ptr: &t.Base.Muted},
		{key: "base.key", ptr: &t.Base.Key},
		{key: "data.key", ptr: &t.Data.Key, from: "base.key"},
		{key: "data.string", ptr: &t.Data.String, from: "base.info"},
		{key: "data.number", ptr: &t.Data.Number, from: "base.primary"},
		{key: "data.true", ptr: &t.Data.True, from: "base.success"},
		{key: "data.false", ptr: &t.Data.False, from: "base.danger"},
		{key: "data.null", ptr: &t.Data.Null, from: "base.muted"},
	}
}

func LoadTheme(lightBackground bool, getenv func(string) string) (*Theme, error) {
	t := darkTheme()
	if lightBackground {
		t = lightTheme()
	}

	resolved := map[string]color.Color{}
	for _, f := range t.fields() {
		env := envName(f.key)
		if v := getenv(env); v != "" {
			c, err := color.Parse(v)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", env, err)
			}
			*f.ptr = c
		} else if f.ptr.IsZero() && f.from != "" {
			*f.ptr = resolved[f.from]
		}
		resolved[f.key] = *f.ptr
	}
	return t, nil
}

func envName(key string) string {
	return "XPCOLOR_THEME_" + strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
}
