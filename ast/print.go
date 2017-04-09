package ast

import (
	"fmt"
	"io"
	"myGo/mytoken"
	"os"
	"reflect"
)

// A FieldFilter may be provided to Fprint to control the output.
type FieldFilter func(name string, value reflect.Value) bool

// NoNilFilter returns true for field values that are not nil
func NotNilFilter(_ string, v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return !v.IsNil()
	}
	return true
}

// Fprintf prints the (sub-)tree starting at AST node x to w.
//
func Fprint(w io.Writer, file *mytoken.File, x interface{}, f FieldFilter) error {
	return fprint(w, file, x, f)
}

func fprint(w io.Writer, file *mytoken.File, x interface{}, f FieldFilter) (err error) {
	// setup printer
	p := printer{
		output: w,
		file:   file,
		filter: f,
		ptrmap: make(map[interface{}]int),
		last:   '\n',
	}
	if x == nil {
		p.printf("nil\n")
		return
	}
	p.print(reflect.ValueOf(x))
	p.printf("\n")
	return
}

// Print prints x to standard output, skipping nil fields.
func Print(file *mytoken.File, x interface{}) error {
	return Fprint(os.Stdout, file, x, NotNilFilter)
}

type printer struct {
	output io.Writer
	file   *mytoken.File
	filter FieldFilter
	ptrmap map[interface{}]int // *T -> line number
	indent int                 // current indentation level
	last   byte                // the last byte processed by Write
	line   int                 // current line number
}

var indent = []byte("*  ")

func (p *printer) Write(data []byte) (n int, err error) {
	var m int
	for i, b := range data {
		if b == '\n' {
			m, err = p.output.Write(data[n : i+1])
			n += m
			if err != nil {
				return
			}
			p.line++
		} else if p.last == '\n' {
			_, err = fmt.Fprintf(p.output, "%6d  ", p.line)
			if err != nil {
				return
			}
			for j := p.indent; j > 0; j-- {
				_, err = p.output.Write(indent)
				if err != nil {
					return
				}
			}
		}
		p.last = b
	}
	if len(data) > n {
		m, err = p.output.Write(data[n:])
		n += m
	}
	return
}

func (p *printer) printf(format string, args ...interface{}) {
	if _, err := fmt.Fprintf(p, format, args...); err != nil {
		panic(err)
	}
}

func (p *printer) print(x reflect.Value) {
	if !NotNilFilter("", x) {
		p.printf("nil")
		return
	}

	switch x.Kind() {
	case reflect.Interface:
		p.print(x.Elem())
	case reflect.Map:
		p.printf("%s (len = %d) {", x.Type(), x.Len())
		if x.Len() > 0 {
			p.indent++
			p.printf("\n")
			for _, key := range x.MapKeys() {
				p.print(key)
				p.printf(": ")
				p.print(x.MapIndex(key))
				p.printf("\n")
			}
			p.indent--
		}
		p.printf("}")

	case reflect.Ptr:
		p.printf("*")
		// type-checked ASTs may contain cycles - use ptrmap
		// to keep track of objects that have been printed
		// already and print the respective line number instead
		ptr := x.Interface()
		if line, exists := p.ptrmap[ptr]; exists {
			p.printf("(obj @ %d)", line)
		} else {
			p.ptrmap[ptr] = p.line
			p.print(x.Elem())
		}

	case reflect.Array:
		p.printf("%s {", x.Type())
		if x.Len() > 0 {
			p.indent++
			p.printf("\n")
			for i, n := 0, x.Len(); i < n; i++ {
				p.printf("%d: ", i)
				p.print(x.Index(i))
				p.printf("\n")
			}
			p.indent--
		}
		p.printf("}")

	case reflect.Slice:
		if s, ok := x.Interface().([]byte); ok {
			p.printf("%#q", s)
			return
		}
		p.printf("%s (len = %d) {", x.Type(), x.Len())
		if x.Len() > 0 {
			p.indent++
			p.printf("\n")
			for i, n := 0, x.Len(); i < n; i++ {
				p.printf("%d: ", i)
				p.print(x.Index(i))
				p.printf("\n")
			}
			p.indent--
		}
		p.printf("}")

	case reflect.Struct:
		t := x.Type()
		p.printf("%s {", t)
		p.indent++
		first := true
		for i, n := 0, t.NumField(); i < n; i++ {
			// exclude non-exported fields because their
			// values cannot be accessed via reflection
			if name := t.Field(i).Name; true {
				value := x.Field(i)
				if p.filter == nil || p.filter(name, value) {
					if first {
						p.printf("\n")
						first = false
					}
					p.printf("%s: ", name)
					p.print(value)
					p.printf("\n")
				}
			}
		}
		p.indent--
		p.printf("}")

	default:
		v := x.Interface()
		switch v := v.(type) {
		case string:
			// print strings in quotes
			p.printf("%q", v)
			return
		case mytoken.Pos:
			// position values can be printed nicely if we have a file set
			if p.file != nil {
				p.printf("%s", p.file.Position(v))
				return
			}
		}
		// default
		p.printf("%v", v)
	}
}
