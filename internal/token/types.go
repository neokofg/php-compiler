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
	T_EQ      // =
	T_EQEQ    // ==
	T_GT      // >
	T_LT      // <
	T_AND     // &&
	T_OR      // ||
	T_NOT     // !
	T_NOTEQ   // !=
	T_GTE     // >=
	T_LTE     // <=
	T_EQEQEQ  // ===
	T_NOTEQEQ // !==

	// -- Separators --
	T_SEMI   // ;
	T_DOLLAR // $
	T_LPAREN // (
	T_RPAREN // )
	T_LBRACE // {
	T_RBRACE // }
	T_COLON  // :

	// -- Statements --
	T_ECHO     // echo
	T_IF       // if
	T_ELSE     // else
	T_WHILE    // while
	T_FOR      // for
	T_BREAK    // break
	T_CONTINUE // continue
	T_DO       // do
	T_SWITCH   // switch
	T_CASE     // case
	T_DEFAULT  // default

	// -- Literals --
	T_TRUE  // true
	T_FALSE // false

	// -- Binary Operators --
	T_DOT // .

	// -- Unary operators --
	T_INC // ++
	T_DEC // --

	// -- Compound assignment operators ---
	T_PLUS_EQ  // +=
	T_MINUS_EQ // -=
	T_MUL_EQ   // *=
	T_DIV_EQ   // /=
	T_MOD_EQ   // %=
	T_DOT_EQ   // .=

	// -- Bitwise operators --
	T_BIT_AND // &
	T_BIT_OR  // |
	T_BIT_XOR // ^
	T_BIT_NOT // ~
	T_LSHIFT  // <<
	T_RSHIFT  // >>

	// -- Modulo --
	T_MOD // %
)
