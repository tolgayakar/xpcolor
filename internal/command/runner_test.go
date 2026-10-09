package command

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fakeCrossplane(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "crossplane")
	script := "#!/bin/sh\necho \"args: $*\"\ncat\necho oops >&2\nexit ${FAKE_EXIT:-0}\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRunPassesThrough(t *testing.T) {
	t.Setenv("XPCOLOR_CROSSPLANE", fakeCrossplane(t))
	t.Setenv("FAKE_EXIT", "3")

	var stdout, stderr bytes.Buffer
	code := Run([]string{"--plain", "resource", "trace", "xr/foo"}, strings.NewReader("from stdin\n"), &stdout, &stderr)

	if code != 3 {
		t.Errorf("exit code = %d, want 3", code)
	}
	if want := "args: resource trace xr/foo\nfrom stdin\n"; stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
	if want := "oops\n"; stderr.String() != want {
		t.Errorf("stderr = %q, want %q", stderr.String(), want)
	}
}

func TestRunColorizesTrace(t *testing.T) {
	t.Setenv("XPCOLOR_CROSSPLANE", fakeCrossplane(t))
	trace := "NAME          SYNCED   READY   STATUS\nXBucket/foo   True     False   Creating\n"

	var stdout, stderr bytes.Buffer
	code := Run([]string{"--force-colors", "resource", "trace", "xr/foo"}, strings.NewReader(trace), &stdout, &stderr)

	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	want := "args: resource trace xr/foo\n" +
		"\033[1mNAME          SYNCED   READY   STATUS\033[0m\n" +
		"\033[96mXBucket\033[0m/foo   \033[32mTrue\033[0m     \033[31mFalse\033[0m   \033[31mCreating\033[0m\n"
	if stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
}

func TestRunInvalidThemeFallsBackToPlain(t *testing.T) {
	t.Setenv("XPCOLOR_CROSSPLANE", fakeCrossplane(t))
	t.Setenv("XPCOLOR_THEME_BASE_DANGER", "pink")
	trace := "NAME          SYNCED   READY   STATUS\n"

	var stdout, stderr bytes.Buffer
	Run([]string{"--force-colors", "resource", "trace", "xr/foo"}, strings.NewReader(trace), &stdout, &stderr)

	if want := "args: resource trace xr/foo\n" + trace; stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
	if !strings.Contains(stderr.String(), "XPCOLOR_THEME_BASE_DANGER") {
		t.Errorf("stderr = %q, want a theme error", stderr.String())
	}
}

func TestRunColorizesStderr(t *testing.T) {
	path := filepath.Join(t.TempDir(), "crossplane")
	script := "#!/bin/sh\necho 'crossplane: error: expected one of \"trace\"' >&2\nexit 80\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XPCOLOR_CROSSPLANE", path)

	var stdout, stderr bytes.Buffer
	code := Run([]string{"--force-colors", "resource"}, strings.NewReader(""), &stdout, &stderr)

	if code != 80 {
		t.Errorf("exit code = %d, want 80", code)
	}
	if want := "\033[31mcrossplane: error:\033[0m expected one of \"trace\"\n"; stderr.String() != want {
		t.Errorf("stderr = %q, want %q", stderr.String(), want)
	}
}

func TestRunMissingBinary(t *testing.T) {
	t.Setenv("XPCOLOR_CROSSPLANE", "xpcolor-does-not-exist")

	var stdout, stderr bytes.Buffer
	code := Run([]string{"version"}, strings.NewReader(""), &stdout, &stderr)

	if code != exitNotFound {
		t.Errorf("exit code = %d, want %d", code, exitNotFound)
	}
	if !strings.Contains(stderr.String(), "not found") {
		t.Errorf("stderr = %q, want a not found message", stderr.String())
	}
}

func TestRunVersion(t *testing.T) {
	old := version
	version = "v1.2.3"
	t.Cleanup(func() { version = old })

	var stdout, stderr bytes.Buffer
	code := Run([]string{"--xpcolor-version"}, strings.NewReader(""), &stdout, &stderr)

	if code != 0 || stdout.String() != "v1.2.3\n" {
		t.Errorf("got code %d, stdout %q", code, stdout.String())
	}
}

func TestShouldColorize(t *testing.T) {
	var buf bytes.Buffer
	tests := []struct {
		name string
		opts Options
		want bool
	}{
		{"non-terminal writer", Options{}, false},
		{"force colors", Options{ForceColors: true}, true},
		{"plain wins over force", Options{ForceColors: true, Plain: true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShouldColorize(tt.opts, &buf); got != tt.want {
				t.Errorf("ShouldColorize = %v, want %v", got, tt.want)
			}
		})
	}
}
