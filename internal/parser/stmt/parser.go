// Licensed under GNU GPL v3. See LICENSE file for details.
package stmt

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type Parser struct {
	context    interfaces.TokenReader
	exprParser interfaces.ExpressionParser

	assignParser       *AssignParser
	echoParser         *EchoParser
	ifParser           *IfParser
	whileParser        *WhileParser
	forParser          *ForParser
	blockParser        *BlockParser
	doWhileParser      *DoWhileParser
	switchParser       *SwitchParser
	functionParser     *FunctionParser
	returnParser       *ReturnParser
	functionCallParser *FunctionCallParser
}

func NewParser(context interfaces.TokenReader, exprParser interfaces.ExpressionParser) interfaces.StatementParser {
	parser := &Parser{
		context:    context,
		exprParser: exprParser,
	}

	parser.blockParser = NewBlockParser(context, parser)

	parser.assignParser = NewAssignParser(context, exprParser)
	parser.echoParser = NewEchoParser(context, exprParser)
	parser.ifParser = NewIfParser(context, exprParser, parser.blockParser)
	parser.whileParser = NewWhileParser(context, exprParser, parser.blockParser)
	parser.forParser = NewForParser(context, exprParser, parser.blockParser)
	parser.doWhileParser = NewDoWhileParser(context, exprParser, parser.blockParser)
	parser.switchParser = NewSwitchParser(context, exprParser, parser)
	parser.functionParser = NewFunctionParser(context, exprParser, parser.blockParser)
	parser.returnParser = NewReturnParser(context, exprParser)
	parser.functionCallParser = NewFunctionCallParser(context, exprParser)

	return parser
}

func (p *Parser) ParseStatement() (ast.Stmt, error) {
	peekedToken := p.context.Peek()

	switch peekedToken.Type {
	case token.T_DOLLAR:
		return p.assignParser.Parse()
	case token.T_ECHO:
		return p.echoParser.Parse()
	case token.T_IF:
		return p.ifParser.Parse()
	case token.T_WHILE:
		return p.whileParser.Parse()
	case token.T_FOR:
		return p.forParser.Parse()
	case token.T_BREAK:
		p.context.Next()
		_, err := p.context.Expect(token.T_SEMI)
		if err != nil {
			return nil, err
		}
		return &ast.BreakStmt{}, nil
	case token.T_CONTINUE:
		p.context.Next()
		_, err := p.context.Expect(token.T_SEMI)
		if err != nil {
			return nil, err
		}
		return &ast.ContinueStmt{}, nil
	case token.T_DO:
		return p.doWhileParser.Parse()
	case token.T_SWITCH:
		return p.switchParser.Parse()
	case token.T_FUNCTION:
		return p.functionParser.Parse()
	case token.T_RETURN:
		return p.returnParser.Parse()
	case token.T_IDENT:
		if p.context.PeekNext().Type == token.T_LPAREN {
			return p.functionCallParser.Parse()
		}
		return nil, fmt.Errorf("Position %d: unexpected identifier: %s", p.context.GetPos(), peekedToken.Value)
	default:
		if peekedToken.Type == token.T_ILLEGAL {
			return nil, fmt.Errorf("Lexer error in position %d: %s", p.context.GetPos(), peekedToken.Value)
		}
		return nil, fmt.Errorf("Position %d: unexpected token in the start of instruction: %v (%q)", p.context.GetPos(), peekedToken.Type, peekedToken.Value)
	}
}

func (p *Parser) ParseBlock() ([]ast.Stmt, error) {
	return p.blockParser.Parse()
}

func (p *Parser) ParseOptionalExpression(terminator token.TokenType) (ast.Expr, error) {
	if p.context.Peek().Type == terminator {
		return nil, nil
	}

	expr, err := p.exprParser.ParseExpression()
	if err != nil {
		return nil, fmt.Errorf("Failed to parse optional expression before %v: %w", terminator, err)
	}

	return expr, nil
}
