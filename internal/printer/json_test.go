package printer

import "testing"

func TestJSONPrinterGolden(t *testing.T) {
	runGolden(t, "json", func() Printer { return &JSONPrinter{Theme: testTheme(t)} })
}
