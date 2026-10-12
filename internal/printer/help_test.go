package printer

import "testing"

func TestHelpPrinterGolden(t *testing.T) {
	runGolden(t, "help", func() Printer { return &HelpPrinter{Theme: testTheme(t)} })
}

func TestStderrPrinterGolden(t *testing.T) {
	runGolden(t, "stderr", func() Printer { return &StderrPrinter{Theme: testTheme(t)} })
}
