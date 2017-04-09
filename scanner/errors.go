package scanner

import (
	//"fmt"
	"myGo/mytoken"
	//	"io"
	//	"sort"
)

type Error struct {
	Pos mytoken.Position
	Msg string
}

func (e Error) Error() string {
	if e.Pos.Filename != "" || e.Pos.IsValid() {
		return e.Pos.String() + ": " + e.Msg
	}
	return e.Msg
}

// ErrorList is a list of *Errors.
// The zero value for an ErrorList is an empty ErrorList ready to use
//
type ErrorList []*Error

// add an error to an errorlist
func (p *ErrorList) Add(pos mytoken.Position, msg string) {
	*p = append(*p, &Error{pos, msg})
}
