package parser

import (
	"fmt"
	"rog/ast"
	"rog/lexer"
	"rog/token"
)

type Parser struct {
	l *lexer.Lexer

	errors []string

	curToken token.Token
	peekToken token.Token
}


func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
		errors: []string{}}
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// parse主逻辑
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {		// 反复查看Token并构造Statement，直至遇到EOF
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	
	if !p.expectPeek(token.IDENT) { 	// let后必须跟identifier
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {	// 然后跟赋值符号
		return nil
	}

	// 目前先跳过expression直至;出现
	for !p.curTokenIs(token.SEMICOLON) {	// 暂时跳过表达式
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	
	// 目前先跳过expression直至;出现
	for !p.curTokenIs(token.SEMICOLON) {	// 暂时跳过表达式
		p.nextToken()
	}

	return stmt
}


func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}