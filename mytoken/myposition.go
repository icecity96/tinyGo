package mytoken

import (
	"fmt"
	"sort"
)

// Position describes an artibitrary source position
// include the file, line, and column location.
// A Position is valid if the line number is > 0.

type Position struct {
	Filename string // filename
	Offset   int    // offset, starting at 0
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1
}

// IsValid reports whether the position is valid.
//
func (pos *Position) IsValid() bool {
	return pos.Line > 0
}

// String returns a string in one of several forms:
//
//	file:line:column    valid position with file name
//	line:column         valid position without file name
//	file                invalid position with file name
//	-                   invalid position without file name
//
func (pos Position) String() string {
	s := pos.Filename
	if pos.IsValid() {
		if s != "" {
			s += ":"
		}
		s += fmt.Sprintf("%d:%d", pos.Line, pos.Column)
	}
	if s == "" {
		s = "-"
	}
	return s
}

// range [0,size]
type Pos int

const NoPos Pos = 0

func (p Pos) IsValid() bool {
	return p != NoPos
}

// TODO: needed?
//
// File
type File struct {
	name  string // file name
	base  int
	size  int // file size
	lines []int
	infos []lineInfo
}

// Name returns the file name of file
func (f *File) Name() string {
	return f.name
}

// Base returns the base offset of file, reserved
func (f *File) Base() int {
	return f.base
}

// Size returns the size of file f
func (f *File) Size() int {
	return f.size
}

// LineCount returns the numbers in file f.
func (f *File) LineCount() int {
	n := len(f.lines)
	return n
}

// AddLine adds the line offset for a new line.
func (f *File) AddLine(offset int) {
	if i := len(f.lines); (i == 0 || f.lines[i-1] < offset) && offset < f.size {
		f.lines = append(f.lines, offset)
	}
}

func (f *File) MergeLine(line int) {
	if line <= 0 {
		panic("illegal line number (line numbering starts at 1)")
	}
	if line >= len(f.lines) {
		panic("illegal line number")
	}
	copy(f.lines[line:], f.lines[line+1:])
	f.lines = f.lines[:len(f.lines)-1]
}

func (f *File) SetLines(lines []int) bool {
	// verify validity of lines table
	size := f.size
	for i, offset := range lines {
		if i > 0 && offset <= lines[i-1] || size <= offset {
			return false
		}
	}
	// set lines table
	f.lines = lines
	return true
}

func (f *File) SetLinesForContent(content []byte) {
	var lines []int
	line := 0
	for offset, b := range content {
		if line >= 0 {
			lines = append(lines, line)
		}
		line = -1
		if b == '\n' {
			line = offset + 1
		}
	}
	// set lines table
	f.lines = lines
}

type lineInfo struct {
	// fields are exported to make them accessible to gob
	Offset   int
	Filename string
	Line     int
}

func (f *File) AddLineInfo(offset int, filename string, line int) {
	if i := len(f.infos); i == 0 || f.infos[i-1].Offset < offset && offset < f.size {
		f.infos = append(f.infos, lineInfo{offset, filename, line})
	}
}

func (f *File) Pos(offset int) Pos {
	if offset > f.size {
		panic("illegal file offset")
	}
	return Pos(f.base + offset)
}

func (f *File) Offset(p Pos) int {
	if int(p) < f.base || int(p) > f.base+f.size {
		panic("illegal Pos value")
	}
	return int(p) - f.base
}

func (f *File) Line(p Pos) int {
	return f.Position(p).Line
}

func searchLineInfos(a []lineInfo, x int) int {
	return sort.Search(len(a), func(i int) bool { return a[i].Offset > x }) - 1
}

func (f *File) unpack(offset int, adjusted bool) (filename string, line, column int) {
	filename = f.name
	if i := searchInts(f.lines, offset); i >= 0 {
		line, column = i+1, offset-f.lines[i]+1
	}
	if adjusted && len(f.infos) > 0 {
		// almost no files have extra line infos
		if i := searchLineInfos(f.infos, offset); i >= 0 {
			alt := &f.infos[i]
			filename = alt.Filename
			if i := searchInts(f.lines, alt.Offset); i >= 0 {
				line += alt.Line - i - 1
			}
		}
	}
	return
}

func (f *File) position(p Pos, adjusted bool) (pos Position) {
	offset := int(p) - f.base
	pos.Offset = offset
	pos.Filename, pos.Line, pos.Column = f.unpack(offset, adjusted)
	return
}

func (f *File) PositionFor(p Pos, adjusted bool) (pos Position) {
	if p != NoPos {
		if int(p) < f.base || int(p) > f.base+f.size {
			panic("illegal Pos value")
		}
		pos = f.position(p, adjusted)
	}
	return
}

func (f *File) Position(p Pos) (pos Position) {
	return f.PositionFor(p, true)
}

func Newfile(name string, base, size int) *File {
	f := &File{name, base, size, []int{0}, nil}
	return f
}
func searchInts(a []int, x int) int {
	i, j := 0, len(a)
	for i < j {
		h := i + (j-i)/2 // avoid overflow when computing h
		// i ≤ h < j
		if a[h] <= x {
			i = h + 1
		} else {
			j = h
		}
	}
	return i - 1
}
