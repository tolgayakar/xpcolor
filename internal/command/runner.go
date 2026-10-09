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

	"github.com/tolgayakar/xpcolor/internal/config"
	"github.com/tolgayakar/xpcolor/internal/printer"
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

	var outP, errP printer.Printer
	if ShouldColorize(opts, stdout) {
		theme, err := config.LoadTheme(opts.LightBackground, os.Getenv)
		if err != nil {
			fmt.Fprintf(stderr, "xpcolor: %v, printing without colors\n", err)
		} else {
			outP = printerFor(crossplaneArgs, theme)
		}
		// Stderr is line buffered once colored, so only do it for known non-interactive commands.
		if outP != nil && ShouldColorize(opts, stderr) {
			errP = &printer.StderrPrinter{Theme: theme}
		}
	}

	cmd := exec.Command(path, crossplaneArgs...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	finishOut, finishErr := func() {}, func() {}
	if outP != nil {
		cmd.Stdout, finishOut = pipeTo(outP, stdout)
	}
	if errP != nil {
		cmd.Stderr, finishErr = pipeTo(errP, stderr)
	}
	finish := func() {
		finishOut()
		finishErr()
	}

	if err := cmd.Start(); err != nil {
		finish()
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

	err = cmd.Wait()
	finish()
	if err != nil {
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
			return exitErr.ExitCode()
		}
		fmt.Fprintf(stderr, "xpcolor: %v\n", err)
		return 1
	}
	return 0
}

func pipeTo(p printer.Printer, w io.Writer) (io.Writer, func()) {
	pr, pw := io.Pipe()
	done := make(chan struct{})
	go func() {
		defer close(done)
		p.Print(pr, w)
		_, _ = io.Copy(io.Discard, pr)
	}()
	return pw, func() {
		_ = pw.Close()
		<-done
	}
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
