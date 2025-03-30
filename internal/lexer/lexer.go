package lexer

import (
	"strings"
	"unicode"
	"fmt"
	"github.com/neokofg/php-compiler/internal/token"
)

type Lexer struct {
	input []rune
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: []rune(input)}
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	ch := l.input[l.pos]
	l.pos++
	return ch
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) peekNext() rune {
	if l.pos+1 >= len(l.input) {
		return 0
	}
	return l.input[l.pos+1]
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.peek()) {
		l.next()
	}
}

func (l *Lexer) readWhile(cond func(rune) bool) string {
	start := l.pos
	for cond(l.peek()) {
		l.next()
	}
	return string(l.input[start:l.pos])
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	ch := l.peek()
	if ch == 0 {
		return token.Token{Type: token.T_EOF, Value: ""}
	}

	switch ch {
	case '+':
		l.next()
		return token.Token{token.T_PLUS, "+"}
	case '-':
		l.next()
		return token.Token{token.T_MINUS, "-"}
	case '*':
		l.next()
		return token.Token{token.T_STAR, "*"}
	case '/':
		l.next()
		// TODO: Добавить обработку комментариев // и /* */ ?
		return token.Token{token.T_SLASH, "/"}
	case '=':
		l.next()
		if l.peek() == '=' {
			l.next()
			return token.Token{token.T_EQEQ, "=="}
		}
		return token.Token{token.T_EQ, "="}	
	case ';':
		l.next()
		return token.Token{token.T_SEMI, ";"}
	case '$':
		l.next()
		return token.Token{token.T_DOLLAR, "$"}
	case '(':
		l.next()
		return token.Token{token.T_LPAREN, "("}
	case ')':
		l.next()
		return token.Token{token.T_RPAREN, ")"}
	case '{':
		l.next()
		return token.Token{token.T_LBRACE, "{"}
	case '}':
		l.next()
		return token.Token{token.T_RBRACE, "}"}	
	case '>':
		l.next()
		// TODO: Добавить >= ?
		return token.Token{token.T_GT, ">"}	
	case '<':
		l.next()
		// TODO: Добавить <= ?
		return token.Token{token.T_LT, "<"}
	case '&':
		l.next()
		if l.peek() == '&' {
			l.next()
			return token.Token{token.T_AND, "&&"}
		} else {
			return token.Token{token.T_ILLEGAL, "&"}
		}
	case '|':
		l.next()
		if l.peek() == '|' {
			l.next()
			return token.Token{token.T_OR, "||"}
		} else {
			return token.Token{token.T_ILLEGAL, "|"}
		}
		case '"': // FIX: Улучшенная обработка строк с экранированием
		l.next() // съели открывающую "
		var sb strings.Builder
		startPos := l.pos
		for {
			ch := l.peek()
			if ch == '"' {
				l.next()
				break
			}
			if ch == 0 {
				l.pos = startPos
				return token.Token{token.T_ILLEGAL, "Unexpected end of string"}
			}
			if ch == '\\' {
				l.next()
				nextCh := l.peek()
				switch nextCh {
				case 'n':
					sb.WriteRune('\n')
					l.next()
				case 't':
					sb.WriteRune('\t')
					l.next()
				case 'r':
					sb.WriteRune('\r')
					l.next()
				case '"':
					sb.WriteRune('"')
					l.next()
				case '\\':
					sb.WriteRune('\\')
					l.next()
				default:
					sb.WriteRune('\\')
					sb.WriteRune(l.next())
				}
			} else {
				sb.WriteRune(l.next())
			}
		}
		return token.Token{token.T_STRING, sb.String()}
	default:
		if unicode.IsDigit(ch) {
			// TODO: Добавить поддержку float?
			val := l.readWhile(unicode.IsDigit)
			return token.Token{token.T_NUMBER, val}
		}
	
		if unicode.IsLetter(ch) || ch == '_' {
			val := l.readWhile(func(r rune) bool {
				return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
			})
			switch val {
			case "echo":
				return token.Token{token.T_ECHO, val}
			case "if":
				return token.Token{token.T_IF, val}
			case "else":
				return token.Token{token.T_ELSE, val}
			case "while":
				return token.Token{token.T_WHILE, val}
			case "for":
				return token.Token{token.T_FOR, val}
			// TODO: Добавить другие ключевые слова (while, for, function, ...)
			}		
			return token.Token{token.T_IDENT, val}
		}

		charStr := string(ch)
		l.next()
		return token.Token{token.T_ILLEGAL, fmt.Sprintf("Undefined symbol: '%s'", charStr)}
	}
}
