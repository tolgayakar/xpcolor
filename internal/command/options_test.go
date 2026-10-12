package command

import (
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		env      map[string]string
		wantOpts Options
		wantArgs []string
	}{
		{
			name:     "no flags",
			args:     []string{"resource", "trace", "xr/foo"},
			wantOpts: Options{CrossplaneCmd: "crossplane"},
			wantArgs: []string{"resource", "trace", "xr/foo"},
		},
		{
			name:     "xpcolor flags are stripped",
			args:     []string{"--force-colors", "resource", "trace", "--light-background", "xr/foo", "--plain", "--xpcolor-version"},
			wantOpts: Options{ForceColors: true, Plain: true, LightBackground: true, ShowVersion: true, CrossplaneCmd: "crossplane"},
			wantArgs: []string{"resource", "trace", "xr/foo"},
		},
		{
			name:     "args after separator are untouched",
			args:     []string{"composition", "render", "--", "--plain"},
			wantOpts: Options{CrossplaneCmd: "crossplane"},
			wantArgs: []string{"composition", "render", "--", "--plain"},
		},
		{
			name:     "environment variables",
			args:     []string{"version"},
			env:      map[string]string{"XPCOLOR_FORCE_COLORS": "true", "XPCOLOR_LIGHT_BACKGROUND": "1", "NO_COLOR": "1", "XPCOLOR_CROSSPLANE": "/opt/crossplane"},
			wantOpts: Options{ForceColors: true, Plain: true, LightBackground: true, CrossplaneCmd: "/opt/crossplane"},
			wantArgs: []string{"version"},
		},
		{
			name:     "theme flag and env",
			args:     []string{"--xpcolor-theme", "light", "version", "--xpcolor-theme=dark"},
			env:      map[string]string{"XPCOLOR_PRESET": "light"},
			wantOpts: Options{Theme: "dark", CrossplaneCmd: "crossplane"},
			wantArgs: []string{"version"},
		},
		{
			name:     "invalid boolean env is ignored",
			args:     []string{},
			env:      map[string]string{"XPCOLOR_FORCE_COLORS": "yes please"},
			wantOpts: Options{CrossplaneCmd: "crossplane"},
			wantArgs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getenv := func(k string) string { return tt.env[k] }
			opts, args := ParseArgs(tt.args, getenv)
			if opts != tt.wantOpts {
				t.Errorf("opts = %+v, want %+v", opts, tt.wantOpts)
			}
			if !reflect.DeepEqual(args, tt.wantArgs) {
				t.Errorf("args = %q, want %q", args, tt.wantArgs)
			}
		})
	}
}

func TestOptionsPreset(t *testing.T) {
	tests := []struct {
		opts Options
		want string
	}{
		{opts: Options{}, want: ""},
		{opts: Options{LightBackground: true}, want: "light"},
		{opts: Options{LightBackground: true, Theme: "dark"}, want: "dark"},
	}
	for _, tt := range tests {
		if got := tt.opts.Preset(); got != tt.want {
			t.Errorf("%+v.Preset() = %q, want %q", tt.opts, got, tt.want)
		}
	}
}
