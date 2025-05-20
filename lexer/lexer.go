package lexer

import (
	"rog/token"
)

type Lexer struct {
	input 		 string
	position 	 int 		// 当前字符位置
	readPosition int		// 下一个字符位置
	ch			 byte		// 当前字符
}

// 创建Lexer
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// 成员函数 读取下一个字符
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {	
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// 成员函数 根据当前ch生成Token并返回
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()

	switch l.ch {
	case '=':
		t, is_ok := l.makeTwoCharToken()
		if is_ok {
			tok = t
		} else {			
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		t, is_ok := l.makeTwoCharToken()
		if is_ok {
			tok = t
		} else {			
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

/* 获取整数，只支持最普通的十进制数，不支持八进制、十六进制等 */
func (l *Lexer) readNumber() string {
	pos := l.position
	for isDigit(l.ch) {		// 不支持负数
		l.readChar()
	}
	return l.input[pos:l.position]
}

// 获取整个标识符
func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// 普通函数判断ch是否为字母或下划线
func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <='z') || (ch >= 'A' && ch <= 'Z') || (ch == '_')
}

// 根据TokenType和ch生成一个Token
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}


func (l *Lexer)makeTwoCharToken() (token.Token, bool) {
	if (l.ch == '=' && l.peekChar() == '=') {
		l.readChar()
		return token.Token{Type: token.EQ, Literal: "=="}, true
	} 
	
	if (l.ch == '!' && l.peekChar() == '=') {
		l.readChar()
		return token.Token{Type: token.NOT_EQ, Literal: "!="}, true
	}

	return token.Token{}, false
}