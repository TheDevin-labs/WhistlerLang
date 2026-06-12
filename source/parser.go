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

func NewParser(tokens []Token) *Parser { return &Parser{tokens: tokens} }

func (p *Parser) skipNL() {
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type == TOKEN_NEWLINE { p.pos++ }
}

func (p *Parser) peek() Token {
	i := p.pos
	for i < len(p.tokens) && p.tokens[i].Type == TOKEN_NEWLINE { i++ }
	if i >= len(p.tokens) { return Token{Type: TOKEN_EOF} }
	return p.tokens[i]
}

func (p *Parser) advance() Token {
	p.skipNL()
	if p.pos >= len(p.tokens) { return Token{Type: TOKEN_EOF} }
	t := p.tokens[p.pos]; p.pos++
	return t
}

func (p *Parser) expect(tt TokenType) (Token, error) {
	t := p.peek()
	if t.Type != tt {
		return t, fmt.Errorf("line %d: expected %s got %s (%q)", t.Line, tt, t.Type, t.Literal)
	}
	return p.advance(), nil
}

func (p *Parser) Parse() (*Program, error) {
	prog := &Program{}
	p.skipNL()
	for p.peek().Type != TOKEN_EOF {
		s, err := p.parseStmt()
		if err != nil { return nil, err }
		if s != nil { prog.Statements = append(prog.Statements, s) }
		p.skipNL()
	}
	return prog, nil
}

func (p *Parser) parseStmt() (Node, error) {
	switch p.peek().Type {
	case TOKEN_LET:       return p.parseLet()
	case TOKEN_SAY:       return p.parseSay()
	case TOKEN_FN:        return p.parseFn()
	case TOKEN_IF:        return p.parseIf()
	case TOKEN_FOR:       return p.parseFor()
	case TOKEN_BLOCKROCK: return p.parseBlockrock()
	case TOKEN_KNOWNUSE:  return p.parseKnownUse()
	case TOKEN_RETURN:
		p.advance()
		v, err := p.parseExpr(0)
		if err != nil { return nil, err }
		return &ReturnStatement{Value: v}, nil
	default:
		e, err := p.parseExpr(0)
		if err != nil { return nil, err }
		return &ExpressionStatement{Expr: e}, nil
	}
}

func (p *Parser) parseLet() (Node, error) {
	p.advance()
	name, err := p.expect(TOKEN_IDENT)
	if err != nil { return nil, err }

	typeAnnotation := ""
	if p.peek().Type == TOKEN_COLON {
		p.advance()
		t := p.peek()
		switch t.Type {
		case TOKEN_TYPE_INT, TOKEN_TYPE_FLOAT, TOKEN_TYPE_BOOL, TOKEN_TYPE_STR,
			TOKEN_TYPE_BYTE, TOKEN_TYPE_BYTES, TOKEN_TYPE_CMPLX,
			TOKEN_TYPE_ARR, TOKEN_TYPE_MAT, TOKEN_IDENT:
			typeAnnotation = t.Literal
			p.advance()
		default:
			return nil, fmt.Errorf("line %d: expected type after ':', got %s", t.Line, t.Type)
		}
	}

	if _, err := p.expect(TOKEN_ASSIGN); err != nil { return nil, err }
	val, err := p.parseExpr(0)
	if err != nil { return nil, err }
	return &LetStatement{Name: name.Literal, TypeAnnotation: typeAnnotation, Value: val}, nil
}

func (p *Parser) parseSay() (Node, error) {
	p.advance()
	val, err := p.parseExpr(0)
	if err != nil { return nil, err }
	return &SayStatement{Value: val}, nil
}

func (p *Parser) parseFn() (Node, error) {
	p.advance()
	name, err := p.expect(TOKEN_IDENT)
	if err != nil { return nil, err }
	if _, err := p.expect(TOKEN_LPAREN); err != nil { return nil, err }

	var params []FnParam
	for p.peek().Type != TOKEN_RPAREN && p.peek().Type != TOKEN_EOF {
		paramName, err := p.expect(TOKEN_IDENT)
		if err != nil { return nil, err }
		typeAnn := ""
		if p.peek().Type == TOKEN_COLON {
			p.advance()
			t := p.peek()
			switch t.Type {
			case TOKEN_TYPE_INT, TOKEN_TYPE_FLOAT, TOKEN_TYPE_BOOL, TOKEN_TYPE_STR,
				TOKEN_TYPE_BYTE, TOKEN_TYPE_BYTES, TOKEN_TYPE_CMPLX,
				TOKEN_TYPE_ARR, TOKEN_TYPE_MAT, TOKEN_IDENT:
				typeAnn = t.Literal
				p.advance()
			}
		}
		params = append(params, FnParam{Name: paramName.Literal, TypeAnnotation: typeAnn})
		if p.peek().Type == TOKEN_COMMA { p.advance() }
	}
	if _, err := p.expect(TOKEN_RPAREN); err != nil { return nil, err }
	if _, err := p.expect(TOKEN_ARROW); err != nil { return nil, err }

	returnType := ""
	t := p.peek()
	switch t.Type {
	case TOKEN_TYPE_INT, TOKEN_TYPE_FLOAT, TOKEN_TYPE_BOOL, TOKEN_TYPE_STR,
		TOKEN_TYPE_BYTE, TOKEN_TYPE_BYTES, TOKEN_TYPE_CMPLX,
		TOKEN_TYPE_ARR, TOKEN_TYPE_MAT, TOKEN_IDENT:
		returnType = t.Literal
		p.advance()
	}

	body, err := p.parseBlock()
	if err != nil { return nil, err }
	return &FnStatement{Name: name.Literal, Params: params, ReturnType: returnType, Body: body}, nil
}

func (p *Parser) parseIf() (Node, error) {
	p.advance()
	cond, err := p.parseExpr(0)
	if err != nil { return nil, err }
	then, err := p.parseBlock()
	if err != nil { return nil, err }
	stmt := &IfStatement{Condition: cond, Then: then}
	for p.peek().Type == TOKEN_ELIF {
		p.advance()
		ec, err := p.parseExpr(0)
		if err != nil { return nil, err }
		eb, err := p.parseBlock()
		if err != nil { return nil, err }
		stmt.Elifs = append(stmt.Elifs, ElifClause{Condition: ec, Body: eb})
	}
	if p.peek().Type == TOKEN_ELSE {
		p.advance()
		eb, err := p.parseBlock()
		if err != nil { return nil, err }
		stmt.Else = eb
	}
	return stmt, nil
}

func (p *Parser) parseFor() (Node, error) {
	p.advance()
	varTok, err := p.expect(TOKEN_IDENT)
	if err != nil { return nil, err }
	if _, err := p.expect(TOKEN_IN); err != nil { return nil, err }
	var iterable Node
	if p.peek().Type == TOKEN_RANGE {
		p.advance()
		if _, err := p.expect(TOKEN_LPAREN); err != nil { return nil, err }
		end, err := p.parseExpr(0)
		if err != nil { return nil, err }
		if _, err := p.expect(TOKEN_RPAREN); err != nil { return nil, err }
		iterable = &RangeExpr{End: end}
	} else {
		iterable, err = p.parseExpr(0)
		if err != nil { return nil, err }
	}
	body, err := p.parseBlock()
	if err != nil { return nil, err }
	return &ForStatement{Variable: varTok.Literal, Iterable: iterable, Body: body}, nil
}

func (p *Parser) parseBlockrock() (Node, error) {
	p.advance()
	body, err := p.parseBlock()
	if err != nil { return nil, err }
	var panicBody []Node
	if p.peek().Type == TOKEN_PANIC {
		p.advance()
		panicBody, err = p.parseBlock()
		if err != nil { return nil, err }
	}
	return &BlockrockStatement{Body: body, PanicBody: panicBody}, nil
}

func (p *Parser) parseBlock() ([]Node, error) {
	if _, err := p.expect(TOKEN_LBRACE); err != nil { return nil, err }
	p.skipNL()
	var stmts []Node
	for p.peek().Type != TOKEN_RBRACE && p.peek().Type != TOKEN_EOF {
		s, err := p.parseStmt()
		if err != nil { return nil, err }
		if s != nil { stmts = append(stmts, s) }
		p.skipNL()
	}
	if _, err := p.expect(TOKEN_RBRACE); err != nil { return nil, err }
	return stmts, nil
}

func prec(tt TokenType) int {
	switch tt {
	case TOKEN_OR:                                  return 1
	case TOKEN_AND:                                 return 2
	case TOKEN_EQ, TOKEN_NEQ:                       return 3
	case TOKEN_LT, TOKEN_LTE, TOKEN_GT, TOKEN_GTE:  return 4
	case TOKEN_PLUS, TOKEN_MINUS:                   return 5
	case TOKEN_STAR, TOKEN_SLASH, TOKEN_PERCENT:    return 6
	case TOKEN_CARET:                               return 7
	}
	return 0
}

func (p *Parser) parseExpr(minPrec int) (Node, error) {
	left, err := p.parsePrimary()
	if err != nil { return nil, err }
	for {
		t := p.peek()
		pr := prec(t.Type)
		if pr <= minPrec { break }
		op := p.advance()
		right, err := p.parseExpr(pr)
		if err != nil { return nil, err }
		left = &BinaryExpr{Left: left, Operator: op.Literal, Right: right}
	}
	return left, nil
}

func (p *Parser) parsePrimary() (Node, error) {
	t := p.peek()

	switch t.Type {
	case TOKEN_INT:
		p.advance()
		v, _ := strconv.ParseInt(t.Literal, 10, 64)
		return &IntLiteral{Value: v}, nil

	case TOKEN_HEX:
		p.advance()
		v, _ := strconv.ParseUint(t.Literal, 16, 64)
		return &ByteLiteral{Value: uint8(v)}, nil

	case TOKEN_FLOAT:
		p.advance()
		v, _ := strconv.ParseFloat(t.Literal, 64)
		return &FloatLiteral{Value: v}, nil

	case TOKEN_COMPLEX:
		p.advance()
		raw := strings.TrimSuffix(t.Literal, "i")
		v, _ := strconv.ParseFloat(raw, 64)
		return &ComplexLiteral{Imag: v}, nil

	case TOKEN_STRING:
		p.advance()
		return &StringLiteral{Value: t.Literal}, nil

	case TOKEN_TRUE:
		p.advance(); return &BoolLiteral{Value: true}, nil

	case TOKEN_FALSE:
		p.advance(); return &BoolLiteral{Value: false}, nil

	case TOKEN_MINUS:
		p.advance()
		operand, err := p.parsePrimary()
		if err != nil { return nil, err }
		return &UnaryExpr{Operator: "-", Operand: operand}, nil

	case TOKEN_NOT:
		p.advance()
		operand, err := p.parsePrimary()
		if err != nil { return nil, err }
		return &UnaryExpr{Operator: "not", Operand: operand}, nil

	case TOKEN_LPAREN:
		p.advance()
		expr, err := p.parseExpr(0)
		if err != nil { return nil, err }
		if _, err := p.expect(TOKEN_RPAREN); err != nil { return nil, err }
		return expr, nil

	case TOKEN_LBRACKET:
		return p.parseArrayOrMatrix()

	case TOKEN_CSV:
		return p.parseCsvExpr()

	case TOKEN_IDENT:
		p.advance()
		if p.peek().Type == TOKEN_LPAREN {
			p.advance()
			var args []Node
			for p.peek().Type != TOKEN_RPAREN && p.peek().Type != TOKEN_EOF {
				a, err := p.parseExpr(0)
				if err != nil { return nil, err }
				args = append(args, a)
				if p.peek().Type == TOKEN_COMMA { p.advance() }
			}
			if _, err := p.expect(TOKEN_RPAREN); err != nil { return nil, err }
			return p.indexSuffix(&CallExpr{Function: t.Literal, Args: args})
		}
		return p.indexSuffix(&Identifier{Name: t.Literal})
	}

	return nil, fmt.Errorf("line %d: unexpected token %s (%q)", t.Line, t.Type, t.Literal)
}

func (p *Parser) parseCsvExpr() (Node, error) {
	p.advance() // consume 'csv'
	if _, err := p.expect(TOKEN_DOT); err != nil { return nil, err }
	method := p.advance()
	if _, err := p.expect(TOKEN_LPAREN); err != nil { return nil, err }
	path, err := p.parseExpr(0)
	if err != nil { return nil, err }
	if _, err := p.expect(TOKEN_RPAREN); err != nil { return nil, err }
	switch method.Literal {
	case "open": return &CsvOpenExpr{Path: path}, nil
	case "line": return &CsvLineExpr{Path: path}, nil
	}
	return nil, fmt.Errorf("unknown csv method: %s", method.Literal)
}

func (p *Parser) indexSuffix(node Node) (Node, error) {
	for p.peek().Type == TOKEN_LBRACKET {
		p.advance()
		idx, err := p.parseExpr(0)
		if err != nil { return nil, err }
		if _, err := p.expect(TOKEN_RBRACKET); err != nil { return nil, err }
		node = &IndexExpr{Object: node, Index: idx}
	}
	return node, nil
}

func (p *Parser) parseArrayOrMatrix() (Node, error) {
	p.advance()
	p.skipNL()
	if p.peek().Type == TOKEN_RBRACKET { p.advance(); return &ArrayLiteral{}, nil }

	allHex := true
	savedPos := p.pos

	if p.peek().Type == TOKEN_LBRACKET {
		var rows [][]Node
		for p.peek().Type == TOKEN_LBRACKET {
			p.advance()
			var row []Node
			for p.peek().Type != TOKEN_RBRACKET && p.peek().Type != TOKEN_EOF {
				el, err := p.parseExpr(0)
				if err != nil { return nil, err }
				row = append(row, el)
				if p.peek().Type == TOKEN_COMMA { p.advance() }
			}
			if _, err := p.expect(TOKEN_RBRACKET); err != nil { return nil, err }
			rows = append(rows, row)
			p.skipNL()
			if p.peek().Type == TOKEN_COMMA { p.advance(); p.skipNL() }
		}
		if _, err := p.expect(TOKEN_RBRACKET); err != nil { return nil, err }
		return &MatrixLiteral{Rows: rows}, nil
	}

	_ = savedPos
	var elements []Node
	for p.peek().Type != TOKEN_RBRACKET && p.peek().Type != TOKEN_EOF {
		el, err := p.parseExpr(0)
		if err != nil { return nil, err }
		if _, ok := el.(*ByteLiteral); !ok { allHex = false }
		elements = append(elements, el)
		if p.peek().Type == TOKEN_COMMA { p.advance() }
	}
	if _, err := p.expect(TOKEN_RBRACKET); err != nil { return nil, err }

	if allHex && len(elements) > 0 {
		var byteVals []uint8
		for _, el := range elements {
			if b, ok := el.(*ByteLiteral); ok {
				byteVals = append(byteVals, b.Value)
			}
		}
		return &BytesLiteral{Values: byteVals}, nil
	}
	return &ArrayLiteral{Elements: elements}, nil
}

func (p *Parser) parseKnownUse() (Node, error) {
	p.advance()
	body, err := p.parseBlock()
	if err != nil { return nil, err }
	return &KnownUseStatement{Body: body}, nil
}
