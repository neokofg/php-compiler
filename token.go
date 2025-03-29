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
)

type Token struct {
	Type  TokenType
	Value string
}
