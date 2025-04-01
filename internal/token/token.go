// Licensed under GNU GPL v3. See LICENSE file for details.
package token

import "fmt"

type Token struct {
	Type  TokenType
	Value string
}

type UnexpectedTokenError struct {
	Expected string
	Found    Token
	Pos      int
}

func (e *UnexpectedTokenError) Error() string {
	return fmt.Sprintf("Position %d: expected %s, but found %v (%q)",
		e.Pos, e.Expected, e.Found.Type, e.Found.Value)
}
