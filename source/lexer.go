package main

import (
	"fmt"
	"strings"
)

type TokenType string

const (
	TOKEN_INT     TokenType = "INT"
	TOKEN_FLOAT   TokenType = "FLOAT"
	TOKEN_COMPLEX TokenType = "COMPLEX"
	TOKEN_STRING  TokenType = "STRING"
	TOKEN_BOOL    TokenType = "BOOL"

	TOKEN_IDENT  TokenType = "IDENT"
	TOKEN_LET    TokenType = "LET"
	TOKEN_FN     TokenType = "FN"
	TOKEN_SAY    TokenType = "SAY"
	TOKEN_IF     TokenType = "IF"
	TOKEN_ELIF   TokenType = "ELIF"
	TOKEN_ELSE   TokenType = "ELSE"
	TOKEN_FOR    TokenType = "FOR"
	TOKEN_IN     TokenType = "IN"
	TOKEN_RANGE  TokenType = "RANGE"
	TOKEN_RETURN TokenType = "RETURN"
	TOKEN_TRUE   TokenType = "TRUE"
	TOKEN_FALSE  TokenType = "FALSE"

	TOKEN_PLUS     TokenType = "PLUS"
	TOKEN_MINUS    TokenType = "MINUS"
	TOKEN_STAR     TokenType = "STAR"
	TOKEN_SLASH    TokenType = "SLASH"
	TOKEN_PERCENT  TokenType = "PERCENT"
	TOKEN_CARET    TokenType = "CARET"
	TOKEN_ASSIGN   TokenType = "ASSIGN"
	TOKEN_EQ       TokenType = "EQ"
	TOKEN_NEQ      TokenType = "NEQ"
	TOKEN_LT       TokenType = "LT"
	TOKEN_LTE      TokenType = "LTE"
	TOKEN_GT       TokenType = "GT"
	TOKEN_GTE      TokenType = "GTE"
	TOKEN_AND      TokenType = "AND"
	TOKEN_OR       TokenType = "OR"
	TOKEN_NOT      TokenType = "NOT"
	TOKEN_ARROW    TokenType = "ARROW"

	TOKEN_LPAREN   TokenType = "LPAREN"
	TOKEN_RPAREN   TokenType = "RPAREN"
	TOKEN_LBRACE   TokenType = "LBRACE"
	TOKEN_RBRACE   TokenType = "RBRACE"
	TOKEN_LBRACKET TokenType = "LBRACKET"
	TOKEN_RBRACKET TokenType = "RBRACKET"
	TOKEN_COMMA    TokenType = "COMMA"
	TOKEN_DOT      TokenType = "DOT"

	TOKEN_NEWLINE TokenType = "NEWLINE"
	TOKEN_EOF     TokenType = "EOF"
	TOKEN_ILLEGAL TokenType = "ILLEGAL"
)

var keywords = map[string]TokenType{
	"let":    TOKEN_LET,
	"fn":     TOKEN_FN,
	"say":    TOKEN_SAY,
	"if":     TOKEN_IF,
	"elif":   TOKEN_ELIF,
	"else":   TOKEN_ELSE,
	"for":    TOKEN_FOR,
	"in":     TOKEN_IN,
	"range":  TOKEN_RANGE,
	"return": TOKEN_RETURN,
	"true":   TOKEN_TRUE,
	"false":  TOKEN_FALSE,
	"and":    TOKEN_AND,
	"or":     TOKEN_OR,
	"not":    TOKEN_NOT,
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

func (t Token) String() string {
	return fmt.Sprintf("Token{%s, %q, line:%d}", t.Type, t.Literal, t.Line)
}

type Lexer struct {
	input  string
	pos    int
	line   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, pos: 0, line: 1}
}

func (l *Lexer) Tokenize() []Token {
	var tokens []Token
	for {
		tok := l.nextToken()
		tokens = append(tokens, tok)
		if tok.Type == TOKEN_EOF {
			break
		}
	}
	return tokens
}

func (l *Lexer) peek() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) peekNext() byte {
	if l.pos+1 >= len(l.input) {
		return 0
	}
	return l.input[l.pos+1]
}

func (l *Lexer) advance() byte {
	ch := l.input[l.pos]
	l.pos++
	return ch
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && (l.input[l.pos] == ' ' || l.input[l.pos] == '\t' || l.input[l.pos] == '\r') {
		l.pos++
	}
}

func (l *Lexer) skipComment() {
	for l.pos < len(l.input) && l.input[l.pos] != '\n' {
		l.pos++
	}
}

func (l *Lexer) readString() string {
	var sb strings.Builder
	for l.pos < len(l.input) && l.input[l.pos] != '"' {
		if l.input[l.pos] == '\\' {
			l.pos++
			switch l.input[l.pos] {
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case '"':
				sb.WriteByte('"')
			case '\\':
				sb.WriteByte('\\')
			}
		} else {
			sb.WriteByte(l.input[l.pos])
		}
		l.pos++
	}
	l.pos++
	return sb.String()
}

func (l *Lexer) readNumber() Token {
	start := l.pos
	isFloat := false

	for l.pos < len(l.input) && (isDigit(l.input[l.pos]) || l.input[l.pos] == '.') {
		if l.input[l.pos] == '.' {
			isFloat = true
		}
		l.pos++
	}

	if l.pos < len(l.input) && l.input[l.pos] == 'i' {
		l.pos++
		return Token{TOKEN_COMPLEX, l.input[start:l.pos], l.line}
	}

	if isFloat {
		return Token{TOKEN_FLOAT, l.input[start:l.pos], l.line}
	}
	return Token{TOKEN_INT, l.input[start:l.pos], l.line}
}

func (l *Lexer) readIdent() string {
	start := l.pos
	for l.pos < len(l.input) && (isLetter(l.input[l.pos]) || isDigit(l.input[l.pos]) || l.input[l.pos] == '_') {
		l.pos++
	}
	return l.input[start:l.pos]
}

func (l *Lexer) nextToken() Token {
	l.skipWhitespace()

	if l.pos >= len(l.input) {
		return Token{TOKEN_EOF, "", l.line}
	}

	ch := l.input[l.pos]

	if ch == '-' && l.peekNext() == '-' {
		l.skipComment()
		return l.nextToken()
	}

	if ch == '\n' {
		l.pos++
		l.line++
		return Token{TOKEN_NEWLINE, "\n", l.line}
	}

	if ch == '"' {
		l.pos++
		str := l.readString()
		return Token{TOKEN_STRING, str, l.line}
	}

	if isDigit(ch) {
		return l.readNumber()
	}

	if isLetter(ch) || ch == '_' {
		ident := l.readIdent()
		if tt, ok := keywords[ident]; ok {
			return Token{tt, ident, l.line}
		}
		return Token{TOKEN_IDENT, ident, l.line}
	}

	l.pos++
	switch ch {
	case '+':
		return Token{TOKEN_PLUS, "+", l.line}
	case '-':
		if l.peek() == '>' {
			l.pos++
			return Token{TOKEN_ARROW, "->", l.line}
		}
		return Token{TOKEN_MINUS, "-", l.line}
	case '*':
		return Token{TOKEN_STAR, "*", l.line}
	case '/':
		return Token{TOKEN_SLASH, "/", l.line}
	case '%':
		return Token{TOKEN_PERCENT, "%", l.line}
	case '^':
		return Token{TOKEN_CARET, "^", l.line}
	case '=':
		if l.peek() == '=' {
			l.pos++
			return Token{TOKEN_EQ, "==", l.line}
		}
		return Token{TOKEN_ASSIGN, "=", l.line}
	case '!':
		if l.peek() == '=' {
			l.pos++
			return Token{TOKEN_NEQ, "!=", l.line}
		}
		return Token{TOKEN_NOT, "!", l.line}
	case '<':
		if l.peek() == '=' {
			l.pos++
			return Token{TOKEN_LTE, "<=", l.line}
		}
		return Token{TOKEN_LT, "<", l.line}
	case '>':
		if l.peek() == '=' {
			l.pos++
			return Token{TOKEN_GTE, ">=", l.line}
		}
		return Token{TOKEN_GT, ">", l.line}
	case '(':
		return Token{TOKEN_LPAREN, "(", l.line}
	case ')':
		return Token{TOKEN_RPAREN, ")", l.line}
	case '{':
		return Token{TOKEN_LBRACE, "{", l.line}
	case '}':
		return Token{TOKEN_RBRACE, "}", l.line}
	case '[':
		return Token{TOKEN_LBRACKET, "[", l.line}
	case ']':
		return Token{TOKEN_RBRACKET, "]", l.line}
	case ',':
		return Token{TOKEN_COMMA, ",", l.line}
	case '.':
		return Token{TOKEN_DOT, ".", l.line}
	}

	return Token{TOKEN_ILLEGAL, string(ch), l.line}
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

