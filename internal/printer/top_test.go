package printer

import "testing"

func TestTopPrinterGolden(t *testing.T) {
	runGolden(t, "top", func() Printer { return &TopPrinter{Theme: testTheme(t)} })
}
