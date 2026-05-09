package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// EvalExpressionWithVars evaluates a math expression string, substituting
// any known variable names before parsing.
func EvalExpressionWithVars(expr string, vars map[string]Value) (float64, error) {
	expr = strings.TrimSpace(expr)

	// Substitute variables (longest-name-first to avoid partial matches)
	// We do a simple token-by-token substitution.
	tokens := tokenize(expr)
	out := strings.Builder{}
	for _, tok := range tokens {
		if v, ok := vars[tok]; ok {
			if v.Kind == ValNumber {
				out.WriteString(strconv.FormatFloat(v.Num, 'f', -1, 64))
				continue
			}
			// String variable in a math context — not numeric
			return 0, fmt.Errorf("variable %q is not a number", tok)
		}
		out.WriteString(tok)
	}
	substituted := strings.TrimSpace(out.String())
	return parseExpr(substituted)
}

// tokenize splits an expression into operator/paren tokens and identifier/number tokens.
func tokenize(expr string) []string {
	var tokens []string
	cur := strings.Builder{}
	flush := func() {
		if cur.Len() > 0 {
			tokens = append(tokens, cur.String())
			cur.Reset()
		}
	}
	for i := 0; i < len(expr); i++ {
		c := expr[i]
		switch c {
		case '+', '-', '*', '/', '%', '(', ')', '^', ' ', '\t':
			flush()
			if c != ' ' && c != '\t' {
				tokens = append(tokens, string(c))
			}
		default:
			cur.WriteByte(c)
		}
	}
	flush()
	return tokens
}

// Recursive-descent parser for math expressions.
// Grammar:
//   expr   = term   (('+' | '-') term)*
//   term   = factor (('*' | '/' | '%') factor)*
//   factor = unary ('^' factor)?        (right-associative power)
//   unary  = '-' unary | primary
//   primary = number | '(' expr ')'

type parser struct {
	input string
	pos   int
}

func parseExpr(s string) (float64, error) {
	p := &parser{input: strings.TrimSpace(s)}
	v, err := p.expr()
	if err != nil {
		return 0, err
	}
	p.skipSpaces()
	if p.pos != len(p.input) {
		return 0, fmt.Errorf("unexpected character at position %d: %q", p.pos, p.input[p.pos:])
	}
	return v, nil
}

func (p *parser) skipSpaces() {
	for p.pos < len(p.input) && (p.input[p.pos] == ' ' || p.input[p.pos] == '\t') {
		p.pos++
	}
}

func (p *parser) peek() byte {
	p.skipSpaces()
	if p.pos >= len(p.input) {
		return 0
	}
	return p.input[p.pos]
}

func (p *parser) consume() byte {
	c := p.input[p.pos]
	p.pos++
	return c
}

func (p *parser) expr() (float64, error) {
	left, err := p.term()
	if err != nil {
		return 0, err
	}
	for {
		c := p.peek()
		if c != '+' && c != '-' {
			break
		}
		p.consume()
		right, err := p.term()
		if err != nil {
			return 0, err
		}
		if c == '+' {
			left += right
		} else {
			left -= right
		}
	}
	return left, nil
}

func (p *parser) term() (float64, error) {
	left, err := p.factor()
	if err != nil {
		return 0, err
	}
	for {
		c := p.peek()
		if c != '*' && c != '/' && c != '%' {
			break
		}
		p.consume()
		right, err := p.factor()
		if err != nil {
			return 0, err
		}
		switch c {
		case '*':
			left *= right
		case '/':
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			left /= right
		case '%':
			if right == 0 {
				return 0, fmt.Errorf("modulo by zero")
			}
			left = math.Mod(left, right)
		}
	}
	return left, nil
}

func (p *parser) factor() (float64, error) {
	base, err := p.unary()
	if err != nil {
		return 0, err
	}
	if p.peek() == '^' {
		p.consume()
		exp, err := p.factor() // right-associative
		if err != nil {
			return 0, err
		}
		return math.Pow(base, exp), nil
	}
	return base, nil
}

func (p *parser) unary() (float64, error) {
	if p.peek() == '-' {
		p.consume()
		v, err := p.unary()
		return -v, err
	}
	return p.primary()
}

func (p *parser) primary() (float64, error) {
	p.skipSpaces()
	if p.pos >= len(p.input) {
		return 0, fmt.Errorf("unexpected end of expression")
	}
	if p.input[p.pos] == '(' {
		p.consume() // '('
		v, err := p.expr()
		if err != nil {
			return 0, err
		}
		p.skipSpaces()
		if p.pos >= len(p.input) || p.input[p.pos] != ')' {
			return 0, fmt.Errorf("missing closing parenthesis")
		}
		p.consume() // ')'
		return v, nil
	}
	// Read a number
	start := p.pos
	if p.pos < len(p.input) && p.input[p.pos] == '-' {
		p.pos++
	}
	for p.pos < len(p.input) && (p.input[p.pos] >= '0' && p.input[p.pos] <= '9' || p.input[p.pos] == '.') {
		p.pos++
	}
	numStr := p.input[start:p.pos]
	if numStr == "" || numStr == "-" {
		return 0, fmt.Errorf("expected number at position %d", start)
	}
	return strconv.ParseFloat(numStr, 64)
}
