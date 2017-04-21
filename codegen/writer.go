package codegen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
)

type Writer struct {
	bytes.Buffer
}

func (w *Writer) Line(text string) {
	io.WriteString(w, text+"\n")
}

func (w *Writer) Linef(format string, a ...interface{}) {
	fmt.Fprintf(w, format+"\n", a...)
}

func (w *Writer) Raw() []byte {
	return w.Bytes()
}

func (w *Writer) Fmt() ([]byte, error) {
	return format.Source(w.Raw())
}
