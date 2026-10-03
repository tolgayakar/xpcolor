package command

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

const exitNotFound = 127

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	opts, crossplaneArgs := ParseArgs(args, os.Getenv)

	if opts.ShowVersion {
		fmt.Fprintln(stdout, Version())
		return 0
	}

	path, err := exec.LookPath(opts.CrossplaneCmd)
	if err != nil {
		fmt.Fprintf(stderr, "xpcolor: %q not found in PATH, set XPCOLOR_CROSSPLANE to its location\n", opts.CrossplaneCmd)
		return exitNotFound
	}

	cmd := exec.Command(path, crossplaneArgs...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(stderr, "xpcolor: %v\n", err)
		return 1
	}

	// Forward signals so crossplane can exit cleanly and we keep its exit code.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigs)
	go func() {
		for s := range sigs {
			_ = cmd.Process.Signal(s)
		}
	}()

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
			return exitErr.ExitCode()
		}
		fmt.Fprintf(stderr, "xpcolor: %v\n", err)
		return 1
	}
	return 0
}

func ShouldColorize(opts Options, w io.Writer) bool {
	if opts.Plain {
		return false
	}
	if opts.ForceColors {
		return true
	}
	f, ok := w.(*os.File)
	return ok && term.IsTerminal(int(f.Fd()))
}
