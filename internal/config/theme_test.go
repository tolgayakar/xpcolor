package config

import "testing"

func TestLoadTheme(t *testing.T) {
	tests := []struct {
		name  string
		light bool
		env   map[string]string
		check func(*Theme) (got, want string)
	}{
		{
			name:  "dark preset",
			check: func(th *Theme) (string, string) { return th.Base.Key.String(), "hicyan" },
		},
		{
			name:  "light preset",
			light: true,
			check: func(th *Theme) (string, string) { return th.Base.Key.String(), "hiblue" },
		},
		{
			name:  "data defaults from base",
			check: func(th *Theme) (string, string) { return th.Data.True.String(), "green" },
		},
		{
			name:  "preset data value wins over base",
			check: func(th *Theme) (string, string) { return th.Data.String.String(), "hiyellow" },
		},
		{
			name:  "env overrides base",
			env:   map[string]string{"XPCOLOR_THEME_BASE_SUCCESS": "higreen:bold"},
			check: func(th *Theme) (string, string) { return th.Base.Success.String(), "higreen:bold" },
		},
		{
			name:  "base override flows into data",
			env:   map[string]string{"XPCOLOR_THEME_BASE_SUCCESS": "higreen"},
			check: func(th *Theme) (string, string) { return th.Data.True.String(), "higreen" },
		},
		{
			name: "env overrides data",
			env: map[string]string{
				"XPCOLOR_THEME_BASE_SUCCESS": "higreen",
				"XPCOLOR_THEME_DATA_TRUE":    "blue",
			},
			check: func(th *Theme) (string, string) { return th.Data.True.String(), "blue" },
		},
		{
			name:  "none disables a color",
			env:   map[string]string{"XPCOLOR_THEME_DATA_NULL": "none"},
			check: func(th *Theme) (string, string) { return th.Data.Null.Render("x"), "x" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th, err := LoadTheme(tt.light, func(k string) string { return tt.env[k] })
			if err != nil {
				t.Fatal(err)
			}
			if got, want := tt.check(th); got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	}
}

func TestLoadThemeInvalidEnv(t *testing.T) {
	env := map[string]string{"XPCOLOR_THEME_BASE_DANGER": "pink"}
	_, err := LoadTheme(false, func(k string) string { return env[k] })
	if err == nil {
		t.Fatal("expected error")
	}
	if want := `XPCOLOR_THEME_BASE_DANGER: invalid color "pink"`; err.Error() != want {
		t.Errorf("err = %q, want %q", err, want)
	}
}

func TestThemeFullyResolved(t *testing.T) {
	for _, light := range []bool{false, true} {
		th, err := LoadTheme(light, func(string) string { return "" })
		if err != nil {
			t.Fatal(err)
		}
		for _, f := range th.fields() {
			if f.ptr.IsZero() {
				t.Errorf("light=%v: %s is unset", light, f.key)
			}
		}
	}
}
