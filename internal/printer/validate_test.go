package printer

import "testing"

func TestValidatePrinterGolden(t *testing.T) {
	runGolden(t, "validate", func() Printer { return &ValidatePrinter{Theme: testTheme(t)} })
}
