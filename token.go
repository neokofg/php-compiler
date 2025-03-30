package main

type TokenType int

const (
	T_EOF TokenType = iota
	T_IDENT
	T_NUMBER
	T_STRING
	T_PLUS
	T_MINUS
	T_STAR
	T_SLASH
	T_EQ
	T_SEMI
	T_DOLLAR
	T_LPAREN
	T_RPAREN
	T_ECHO
	T_IF
	T_ELSE
	T_LBRACE
	T_RBRACE
	T_GT
	T_LT
	T_AND
	T_OR
	T_EQEQ
	T_ILLEGAL
	T_WHILE
	T_FOR
)

type Token struct {
	Type  TokenType
	Value string
}
