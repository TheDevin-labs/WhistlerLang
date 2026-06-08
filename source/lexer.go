package main

type TokenType string

const (
	TOKEN_INT        TokenType = "INT"
	TOKEN_FLOAT      TokenType = "FLOAT"
	TOKEN_COMPLEX    TokenType = "COMPLEX"
	TOKEN_STRING     TokenType = "STRING"
	TOKEN_HEX        TokenType = "HEX"
	TOKEN_IDENT      TokenType = "IDENT"
	TOKEN_LET        TokenType = "LET"
	TOKEN_FN         TokenType = "FN"
	TOKEN_SAY        TokenType = "SAY"
	TOKEN_IF         TokenType = "IF"
	TOKEN_ELIF       TokenType = "ELIF"
	TOKEN_ELSE       TokenType = "ELSE"
	TOKEN_FOR        TokenType = "FOR"
	TOKEN_IN         TokenType = "IN"
	TOKEN_RANGE      TokenType = "RANGE"
	TOKEN_RETURN     TokenType = "RETURN"
	TOKEN_TRUE       TokenType = "TRUE"
	TOKEN_FALSE      TokenType = "FALSE"
	TOKEN_AND        TokenType = "AND"
	TOKEN_OR         TokenType = "OR"
	TOKEN_NOT        TokenType = "NOT"
	TOKEN_BLOCKROCK  TokenType = "BLOCKROCK"
	TOKEN_PANIC      TokenType = "PANIC"
	TOKEN_CSV        TokenType = "CSV"
	TOKEN_KNOWNUSE   TokenType = "KNOWNUSE"
	TOKEN_UNSAFE     TokenType = "UNSAFE"
	TOKEN_TYPE_INT   TokenType = "TYPE_INT"
	TOKEN_TYPE_FLOAT TokenType = "TYPE_FLOAT"
	TOKEN_TYPE_BOOL  TokenType = "TYPE_BOOL"
	TOKEN_TYPE_STR   TokenType = "TYPE_STR"
	TOKEN_TYPE_BYTE  TokenType = "TYPE_BYTE"
	TOKEN_TYPE_BYTES TokenType = "TYPE_BYTES"
	TOKEN_TYPE_CMPLX TokenType = "TYPE_CMPLX"
	TOKEN_TYPE_ARR   TokenType = "TYPE_ARR"
	TOKEN_TYPE_MAT   TokenType = "TYPE_MAT"
	TOKEN_PLUS       TokenType = "PLUS"
	TOKEN_MINUS      TokenType = "MINUS"
	TOKEN_STAR       TokenType = "STAR"
	TOKEN_SLASH      TokenType = "SLASH"
	TOKEN_PERCENT    TokenType = "PERCENT"
	TOKEN_CARET      TokenType = "CARET"
	TOKEN_ASSIGN     TokenType = "ASSIGN"
	TOKEN_EQ         TokenType = "EQ"
	TOKEN_NEQ        TokenType = "NEQ"
	TOKEN_LT         TokenType = "LT"
	TOKEN_LTE        TokenType = "LTE"
	TOKEN_GT         TokenType = "GT"
	TOKEN_GTE        TokenType = "GTE"
	TOKEN_ARROW      TokenType = "ARROW"
	TOKEN_COLON      TokenType = "COLON"
	TOKEN_DOT        TokenType = "DOT"
	TOKEN_LPAREN     TokenType = "LPAREN"
	TOKEN_RPAREN     TokenType = "RPAREN"
	TOKEN_LBRACE     TokenType = "LBRACE"
	TOKEN_RBRACE     TokenType = "RBRACE"
	TOKEN_LBRACKET   TokenType = "LBRACKET"
	TOKEN_RBRACKET   TokenType = "RBRACKET"
	TOKEN_COMMA      TokenType = "COMMA"
	TOKEN_NEWLINE    TokenType = "NEWLINE"
	TOKEN_EOF        TokenType = "EOF"
	TOKEN_ILLEGAL    TokenType = "ILLEGAL"
)

var keywords = map[string]TokenType{
	"let":      TOKEN_LET,
	"fn":       TOKEN_FN,
	"say":      TOKEN_SAY,
	"if":       TOKEN_IF,
	"elif":     TOKEN_ELIF,
	"else":     TOKEN_ELSE,
	"for":      TOKEN_FOR,
	"in":       TOKEN_IN,
	"range":    TOKEN_RANGE,
	"return":   TOKEN_RETURN,
	"true":     TOKEN_TRUE,
	"false":    TOKEN_FALSE,
	"and":      TOKEN_AND,
	"or":       TOKEN_OR,
	"not":      TOKEN_NOT,
	"blockrock": TOKEN_BLOCKROCK,
	"panic":    TOKEN_PANIC,
	"csv":      TOKEN_CSV,
	"_knownuse": TOKEN_KNOWNUSE,
	"unsafe":    TOKEN_UNSAFE,
	"int":      TOKEN_TYPE_INT,
	"float":    TOKEN_TYPE_FLOAT,
	"bool":     TOKEN_TYPE_BOOL,
	"string":   TOKEN_TYPE_STR,
	"byte":     TOKEN_TYPE_BYTE,
	"bytes":    TOKEN_TYPE_BYTES,
	"complex":  TOKEN_TYPE_CMPLX,
	"array":    TOKEN_TYPE_ARR,
	"matrix":   TOKEN_TYPE_MAT,
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

type Lexer struct {
	input string
	pos   int
	line  int
}

func NewLexer(input string) *Lexer { return &Lexer{input: input, line: 1} }

func (l *Lexer) Tokenize() []Token {
	var tokens []Token
	for {
		t := l.next()
		tokens = append(tokens, t)
		if t.Type == TOKEN_EOF { break }
	}
	return tokens
}

func (l *Lexer) ch() byte {
	if l.pos >= len(l.input) { return 0 }
	return l.input[l.pos]
}

func (l *Lexer) peek2() byte {
	if l.pos+1 >= len(l.input) { return 0 }
	return l.input[l.pos+1]
}

func (l *Lexer) next() Token {
	for l.ch() == ' ' || l.ch() == '\t' || l.ch() == '\r' { l.pos++ }

	if l.pos >= len(l.input) { return Token{TOKEN_EOF, "", l.line} }

	c := l.ch()

	if c == '-' && l.peek2() == '-' {
		for l.ch() != '\n' && l.ch() != 0 { l.pos++ }
		return l.next()
	}

	if c == '\n' { l.pos++; l.line++; return Token{TOKEN_NEWLINE, "\n", l.line} }

	if c == '"' {
		l.pos++
		start := l.pos
		for l.ch() != '"' && l.ch() != 0 { l.pos++ }
		s := l.input[start:l.pos]
		l.pos++
		return Token{TOKEN_STRING, s, l.line}
	}

	if c == '0' && l.peek2() == 'x' {
		l.pos += 2
		start := l.pos
		for isHexDigit(l.ch()) { l.pos++ }
		return Token{TOKEN_HEX, l.input[start:l.pos], l.line}
	}

	if isDigit(c) {
		start := l.pos
		isFloat := false
		for isDigit(l.ch()) || l.ch() == '.' {
			if l.ch() == '.' { isFloat = true }
			l.pos++
		}
		if l.ch() == 'i' { l.pos++; return Token{TOKEN_COMPLEX, l.input[start:l.pos], l.line} }
		if isFloat { return Token{TOKEN_FLOAT, l.input[start:l.pos], l.line} }
		return Token{TOKEN_INT, l.input[start:l.pos], l.line}
	}

	if isLetter(c) || c == '_' {
		start := l.pos
		for isLetter(l.ch()) || isDigit(l.ch()) || l.ch() == '_' { l.pos++ }
		word := l.input[start:l.pos]
		if tt, ok := keywords[word]; ok { return Token{tt, word, l.line} }
		return Token{TOKEN_IDENT, word, l.line}
	}

	l.pos++
	switch c {
	case '+': return Token{TOKEN_PLUS, "+", l.line}
	case '-':
		if l.ch() == '>' { l.pos++; return Token{TOKEN_ARROW, "->", l.line} }
		return Token{TOKEN_MINUS, "-", l.line}
	case '*': return Token{TOKEN_STAR, "*", l.line}
	case '/': return Token{TOKEN_SLASH, "/", l.line}
	case '%': return Token{TOKEN_PERCENT, "%", l.line}
	case '^': return Token{TOKEN_CARET, "^", l.line}
	case '=':
		if l.ch() == '=' { l.pos++; return Token{TOKEN_EQ, "==", l.line} }
		return Token{TOKEN_ASSIGN, "=", l.line}
	case '!':
		if l.ch() == '=' { l.pos++; return Token{TOKEN_NEQ, "!=", l.line} }
		return Token{TOKEN_NOT, "!", l.line}
	case '<':
		if l.ch() == '=' { l.pos++; return Token{TOKEN_LTE, "<=", l.line} }
		return Token{TOKEN_LT, "<", l.line}
	case '>':
		if l.ch() == '=' { l.pos++; return Token{TOKEN_GTE, ">=", l.line} }
		return Token{TOKEN_GT, ">", l.line}
	case ':': return Token{TOKEN_COLON, ":", l.line}
	case '.': return Token{TOKEN_DOT, ".", l.line}
	case '(': return Token{TOKEN_LPAREN, "(", l.line}
	case ')': return Token{TOKEN_RPAREN, ")", l.line}
	case '{': return Token{TOKEN_LBRACE, "{", l.line}
	case '}': return Token{TOKEN_RBRACE, "}", l.line}
	case '[': return Token{TOKEN_LBRACKET, "[", l.line}
	case ']': return Token{TOKEN_RBRACKET, "]", l.line}
	case ',': return Token{TOKEN_COMMA, ",", l.line}
	}
	return Token{TOKEN_ILLEGAL, string(c), l.line}
}

func isDigit(c byte) bool    { return c >= '0' && c <= '9' }
func isLetter(c byte) bool   { return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') }
func isHexDigit(c byte) bool { return isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') }
