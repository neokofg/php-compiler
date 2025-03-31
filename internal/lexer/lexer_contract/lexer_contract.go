package lexer_contract

type LexerLike interface {
	GetPos() int
	SetPos(pos int)

	Next() rune
	Peek() rune
	PeekNext() rune
	ReadWhile(cond func(rune) bool) string
}
