// Licensed under GNU GPL v3. See LICENSE file for details.
package compiler

import (
	"fmt"
)

type CompilationError struct {
	Msg string
	Pos int
}

func (e *CompilationError) Error() string {
	if e.Pos >= 0 {
		return fmt.Sprintf("compilation error at position %d: %s", e.Pos, e.Msg)
	}
	return fmt.Sprintf("compilation error: %s", e.Msg)
}

func NewError(msg string) error {
	return &CompilationError{
		Msg: msg,
		Pos: -1,
	}
}

func NewErrorAtPos(msg string, pos int) error {
	return &CompilationError{
		Msg: msg,
		Pos: pos,
	}
}
