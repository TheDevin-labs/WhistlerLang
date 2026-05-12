package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) peek() Token {
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type == TOKEN_NEWLINE {
		p.pos++
	}
	if p.pos >= len(p.tokens) {
		return Token{Type: TOKEN_EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) peekRaw() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TOKEN_EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() Token {
	tok := p.peek()
	p.pos++
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type == TOKEN_NEWLINE {
		p.pos++
	}
	return tok
}

func (p *Parser) expect(tt TokenType) (Token, error) {
	tok := p.peek()
	if tok.Type != tt {
		return tok, fmt.Errorf("line %d: expected %s, got %s (%q)", tok.Line, tt, tok.Type, tok.Literal)
	}
	return p.advance(), nil
}

func (p *Parser) skipNewlines() {
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type == TOKEN_NEWLINE {
		p.pos++
	}
}

func (p *Parser) Parse() (*Program, error) {
	program := &Program{}
	p.skipNewlines()

	for p.peek().Type != TOKEN_EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.skipNewlines()
	}

	return program, nil
}

func (p *Parser) parseStatement() (Node, error) {
	tok := p.peek()

	switch tok.Type {
	case TOKEN_LET:
		return p.parseLetStatement()
	case TOKEN_SAY:
		return p.parseSayStatement()
	case TOKEN_FN:
		return p.parseFnStatement()
	case TOKEN_IF:
		return p.parseIfStatement()
	case TOKEN_FOR:
		return p.parseForStatement()
	case TOKEN_RETURN:
		p.advance()
		val, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		return &ReturnStatement{Value: val}, nil
	default:
		expr, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		return &ExpressionStatement{Expr: expr}, nil
	}
}

func (p *Parser) parseLetStatement() (Node, error) {
	p.advance()
	nameTok, err := p.expect(TOKEN_IDENT)
	if err != nil {
		return nil, err
	}
	_, err = p.expect(TOKEN_ASSIGN)
	if err != nil {
		return nil, err
	}
	val, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}
	return &LetStatement{Name: nameTok.Literal, Value: val}, nil
}

func (p *Parser) parseSayStatement() (Node, error) {
	p.advance()
	val, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}
	return &SayStatement{Value: val}, nil
}

func (p *Parser) parseFnStatement() (Node, error) {
	p.advance()
	nameTok, err := p.expect(TOKEN_IDENT)
	if err != nil {
		return nil, err
	}
	_, err = p.expect(TOKEN_LPAREN)
	if err != nil {
		return nil, err
	}

	var params []string
	for p.peek().Type != TOKEN_RPAREN && p.peek().Type != TOKEN_EOF {
		paramTok, err := p.expect(TOKEN_IDENT)
		if err != nil {
			return nil, err
		}
		params = append(params, paramTok.Literal)
		if p.peek().Type == TOKEN_COMMA {
			p.advance()
		}
	}
	_, err = p.expect(TOKEN_RPAREN)
	if err != nil {
		return nil, err
	}
	_, err = p.expect(TOKEN_ARROW)
	if err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &FnStatement{Name: nameTok.Literal, Params: params, Body: body}, nil
}

func (p *Parser) parseIfStatement() (Node, error) {
	p.advance()
	cond, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}
	then, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	stmt := &IfStatement{Condition: cond, Then: then}

	for p.peek().Type == TOKEN_ELIF {
		p.advance()
		elifCond, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		elifBody, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		stmt.Elifs = append(stmt.Elifs, ElifClause{Condition: elifCond, Body: elifBody})
	}

	if p.peek().Type == TOKEN_ELSE {
		p.advance()
		elseBody, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		stmt.Else = elseBody
	}

	return stmt, nil
}

func (p *Parser) parseForStatement() (Node, error) {
	p.advance()
	varTok, err := p.expect(TOKEN_IDENT)
	if err != nil {
		return nil, err
	}
	_, err = p.expect(TOKEN_IN)
	if err != nil {
		return nil, err
	}

	var iterable Node
	if p.peek().Type == TOKEN_RANGE {
		p.advance()
		_, err = p.expect(TOKEN_LPAREN)
		if err != nil {
			return nil, err
		}
		end, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		_, err = p.expect(TOKEN_RPAREN)
		if err != nil {
			return nil, err
		}
		iterable = &RangeExpr{End: end}
	} else {
		iterable, err = p.parseExpression(0)
		if err != nil {
			return nil, err
		}
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &ForStatement{Variable: varTok.Literal, Iterable: iterable, Body: body}, nil
}

func (p *Parser) parseBlock() ([]Node, error) {
	_, err := p.expect(TOKEN_LBRACE)
	if err != nil {
		return nil, err
	}
	p.skipNewlines()

	var stmts []Node
	for p.peek().Type != TOKEN_RBRACE && p.peek().Type != TOKEN_EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
		p.skipNewlines()
	}

	_, err = p.expect(TOKEN_RBRACE)
	if err != nil {
		return nil, err
	}
	return stmts, nil
}

type prefixParseFn func() (Node, error)
type infixParseFn func(Node) (Node, error)

func (p *Parser) parseExpression(minPrec int) (Node, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for {
		tok := p.peek()
		prec := infixPrec(tok.Type)
		if prec <= minPrec {
			break
		}
		op := p.advance()
		right, err := p.parseExpression(prec)
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Operator: op.Literal, Right: right}
	}

	return left, nil
}

func infixPrec(tt TokenType) int {
	switch tt {
	case TOKEN_OR:
		return 1
	case TOKEN_AND:
		return 2
	case TOKEN_EQ, TOKEN_NEQ:
		return 3
	case TOKEN_LT, TOKEN_LTE, TOKEN_GT, TOKEN_GTE:
		return 4
	case TOKEN_PLUS, TOKEN_MINUS:
		return 5
	case TOKEN_STAR, TOKEN_SLASH, TOKEN_PERCENT:
		return 6
	case TOKEN_CARET:
		return 7
	}
	return 0
}

func (p *Parser) parsePrimary() (Node, error) {
	tok := p.peek()

	switch tok.Type {
	case TOKEN_INT:
		p.advance()
		v, _ := strconv.ParseInt(tok.Literal, 10, 64)
		return &IntLiteral{Value: v}, nil

	case TOKEN_FLOAT:
		p.advance()
		v, _ := strconv.ParseFloat(tok.Literal, 64)
		return &FloatLiteral{Value: v}, nil

	case TOKEN_COMPLEX:
		p.advance()
		raw := strings.TrimSuffix(tok.Literal, "i")
		v, _ := strconv.ParseFloat(raw, 64)
		return &ComplexLiteral{Real: 0, Imag: v}, nil

	case TOKEN_STRING:
		p.advance()
		return &StringLiteral{Value: tok.Literal}, nil

	case TOKEN_TRUE:
		p.advance()
		return &BoolLiteral{Value: true}, nil

	case TOKEN_FALSE:
		p.advance()
		return &BoolLiteral{Value: false}, nil

	case TOKEN_MINUS:
		p.advance()
		operand, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{Operator: "-", Operand: operand}, nil

	case TOKEN_NOT:
		p.advance()
		operand, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{Operator: "not", Operand: operand}, nil

	case TOKEN_LPAREN:
		p.advance()
		expr, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		_, err = p.expect(TOKEN_RPAREN)
		if err != nil {
			return nil, err
		}
		return expr, nil

	case TOKEN_LBRACKET:
		return p.parseArrayOrMatrix()

	case TOKEN_IDENT:
		p.advance()
		if p.peek().Type == TOKEN_LPAREN {
			p.advance()
			var args []Node
			for p.peek().Type != TOKEN_RPAREN && p.peek().Type != TOKEN_EOF {
				arg, err := p.parseExpression(0)
				if err != nil {
					return nil, err
				}
				args = append(args, arg)
				if p.peek().Type == TOKEN_COMMA {
					p.advance()
				}
			}
			_, err := p.expect(TOKEN_RPAREN)
			if err != nil {
				return nil, err
			}
			node := &CallExpr{Function: tok.Literal, Args: args}
			return p.parseIndexSuffix(node)
		}
		node := &Identifier{Name: tok.Literal}
		return p.parseIndexSuffix(node)
	}

	return nil, fmt.Errorf("line %d: unexpected token %s (%q)", tok.Line, tok.Type, tok.Literal)
}

func (p *Parser) parseIndexSuffix(node Node) (Node, error) {
	for p.peek().Type == TOKEN_LBRACKET {
		p.advance()
		idx, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		_, err = p.expect(TOKEN_RBRACKET)
		if err != nil {
			return nil, err
		}
		node = &IndexExpr{Object: node, Index: idx}
	}
	return node, nil
}

func (p *Parser) parseArrayOrMatrix() (Node, error) {
	p.advance()
	p.skipNewlines()

	if p.peek().Type == TOKEN_RBRACKET {
		p.advance()
		return &ArrayLiteral{Elements: nil}, nil
	}

	if p.peek().Type == TOKEN_LBRACKET {
		var rows [][]Node
		for p.peek().Type == TOKEN_LBRACKET {
			p.advance()
			var row []Node
			for p.peek().Type != TOKEN_RBRACKET && p.peek().Type != TOKEN_EOF {
				el, err := p.parseExpression(0)
				if err != nil {
					return nil, err
				}
				row = append(row, el)
				if p.peek().Type == TOKEN_COMMA {
					p.advance()
				}
			}
			_, err := p.expect(TOKEN_RBRACKET)
			if err != nil {
				return nil, err
			}
			rows = append(rows, row)
			p.skipNewlines()
			if p.peek().Type == TOKEN_COMMA {
				p.advance()
				p.skipNewlines()
			}
		}
		_, err := p.expect(TOKEN_RBRACKET)
		if err != nil {
			return nil, err
		}
		return &MatrixLiteral{Rows: rows}, nil
	}

	var elements []Node
	for p.peek().Type != TOKEN_RBRACKET && p.peek().Type != TOKEN_EOF {
		el, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		elements = append(elements, el)
		if p.peek().Type == TOKEN_COMMA {
			p.advance()
		}
	}
	_, err := p.expect(TOKEN_RBRACKET)
	if err != nil {
		return nil, err
	}
	return &ArrayLiteral{Elements: elements}, nil
}

