package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Printer interface {
	Print(r io.Reader, w io.Writer)
}

func eachLine(r io.Reader, w io.Writer, fn func(line string) string) {
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadString('\n')
		if line != "" {
			if text, ok := strings.CutSuffix(line, "\n"); ok {
				fmt.Fprintln(w, fn(text))
			} else {
				fmt.Fprint(w, fn(text))
			}
		}
		if err != nil {
			return
		}
	}
}
