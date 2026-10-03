package main

import (
	"os"

	"github.com/tolgayakar/xpcolor/internal/command"
)

func main() {
	os.Exit(command.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
