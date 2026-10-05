package printer

import "testing"

func TestYAMLPrinterGolden(t *testing.T) {
	runGolden(t, "yaml", func() Printer { return &YAMLPrinter{Theme: testTheme(t)} })
}
