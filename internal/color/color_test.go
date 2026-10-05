package color

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{in: "", want: "x"},
		{in: "none", want: "x"},
		{in: "red", want: "\033[31mx\033[0m"},
		{in: "hicyan", want: "\033[96mx\033[0m"},
		{in: "gray:italic", want: "\033[90;3mx\033[0m"},
		{in: "Green : Bold", want: "\033[32;1mx\033[0m"},
		{in: "fg=white:bg=red", want: "\033[37;41mx\033[0m"},
		{in: "208", want: "\033[38;5;208mx\033[0m"},
		{in: "bg=#c2270a", want: "\033[48;2;194;39;10mx\033[0m"},
		{in: "pink", wantErr: true},
		{in: "#zzzzzz", wantErr: true},
		{in: "256", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			c, err := Parse(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse(%q) err = %v, wantErr %v", tt.in, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got := c.Render("x"); got != tt.want {
				t.Errorf("Render = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRenderEmpty(t *testing.T) {
	if got := MustParse("red").Render(""); got != "" {
		t.Errorf("Render(\"\") = %q, want empty", got)
	}
}

func TestString(t *testing.T) {
	if got := MustParse("red:bold").String(); got != "red:bold" {
		t.Errorf("String = %q", got)
	}
	if got := MustParse("").String(); got != "none" {
		t.Errorf("String = %q, want none", got)
	}
}
