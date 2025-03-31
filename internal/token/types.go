// Licensed under GNU GPL v3. See LICENSE file for details.
package token

type TokenType int

const (
	// -- Shared --
	T_EOF TokenType = iota
	T_ILLEGAL

	// -- Lexems --
	T_IDENT
	T_NUMBER
	T_STRING

	// -- Arithmetic --
	T_PLUS  // +
	T_MINUS // -
	T_STAR  // *
	T_SLASH // /

	// -- Logical --
	T_EQ    // =
	T_EQEQ  // ==
	T_GT    // >
	T_LT    // <
	T_AND   // &&
	T_OR    // ||
	T_NOT   // !
	T_NOTEQ // !=

	// -- Separators --
	T_SEMI   // ;
	T_DOLLAR // $
	T_LPAREN // (
	T_RPAREN // )
	T_LBRACE // {
	T_RBRACE // }

	// -- Statements --
	T_ECHO  // echo
	T_IF    // if
	T_ELSE  // else
	T_WHILE // while
	T_FOR   // for

	// -- Literals --
	T_TRUE  // true
	T_FALSE // false
)
