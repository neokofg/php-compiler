package parser

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/context"
	"github.com/neokofg/php-compiler/internal/parser/expr"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/parser/stmt"
	"github.com/neokofg/php-compiler/internal/token"
)

type Parser struct {
	context    *context.ParserContext
	exprParser interfaces.ExpressionParser
	stmtParser interfaces.StatementParser
}

func NewParser(tokens []token.Token) *Parser {
	ctx := context.NewParserContext(tokens)

	exprParser := expr.NewParser(ctx)
	stmtParser := stmt.NewParser(ctx, exprParser)

	return &Parser{
		context:    ctx,
		exprParser: exprParser,
		stmtParser: stmtParser,
	}
}

func (p *Parser) Next() token.Token {
	return p.context.Next()
}

func (p *Parser) Peek() token.Token {
	return p.context.Peek()
}

func (p *Parser) Expect(t token.TokenType) (token.Token, error) {
	return p.context.Expect(t)
}

func (p *Parser) GetPos() int {
	return p.context.GetPos()
}

func (p *Parser) SetPos(pos int) {
	p.context.SetPos(pos)
}

func (p *Parser) ParseBlock() ([]ast.Stmt, error) {
	return p.stmtParser.ParseBlock()
}

func (p *Parser) ParseOptionalExpression(terminator token.TokenType) (ast.Expr, error) {
	return p.stmtParser.ParseOptionalExpression(terminator)
}

func (p *Parser) Parse() ([]ast.Stmt, error) {
	var stmts []ast.Stmt

	for p.Peek().Type != token.T_EOF {
		stmt, err := p.stmtParser.ParseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
	}

	return stmts, nil
}
