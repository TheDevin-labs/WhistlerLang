package main

import (
	"fmt"
	"strings"
)

type StrictError struct {
	Line    int
	Message string
}

func (e StrictError) Error() string {
	return fmt.Sprintf("strict error line %d: %s", e.Line, e.Message)
}

type StrictChecker struct {
	errors []StrictError
}

func NewStrictChecker() *StrictChecker {
	return &StrictChecker{}
}

func (s *StrictChecker) err(line int, msg string) {
	s.errors = append(s.errors, StrictError{Line: line, Message: msg})
}

func (s *StrictChecker) Check(program *Program) []StrictError {
	for _, stmt := range program.Statements {
		s.checkNode(stmt)
	}
	return s.errors
}

func (s *StrictChecker) checkNode(node Node) {
	switch n := node.(type) {

	case *LetStatement:
		if strings.TrimSpace(n.TypeAnnotation) == "" {
			s.err(0, fmt.Sprintf("variable '%s' must have a type annotation in strict mode\n  hint: let %s: <type> = ...", n.Name, n.Name))
		}

	case *FnStatement:
		for _, p := range n.Params {
			if strings.TrimSpace(p.TypeAnnotation) == "" {
				s.err(0, fmt.Sprintf("parameter '%s' in function '%s' must have a type annotation\n  hint: fn %s(%s: <type>, ...) -> <type> { }", p.Name, n.Name, n.Name, p.Name))
			}
		}
		if strings.TrimSpace(n.ReturnType) == "" {
			s.err(0, fmt.Sprintf("function '%s' must declare a return type in strict mode\n  hint: fn %s(...) -> <type> { }", n.Name, n.Name))
		}
		for _, stmt := range n.Body {
			s.checkNode(stmt)
		}

	case *BlockrockStatement:
		if len(n.PanicBody) == 0 {
			s.err(0, "blockrock without panic handler is not allowed in strict mode\n  hint: blockrock { ... } panic { ... }")
		}
		for _, stmt := range n.Body {
			s.checkNode(stmt)
		}

	case *IfStatement:
		for _, stmt := range n.Then { s.checkNode(stmt) }
		for _, elif := range n.Elifs { for _, stmt := range elif.Body { s.checkNode(stmt) } }
		for _, stmt := range n.Else { s.checkNode(stmt) }

	case *ForStatement:
		for _, stmt := range n.Body { s.checkNode(stmt) }

	case *KnownUseStatement:
		return
	}
}

func FormatStrictErrors(errs []StrictError) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n╔══════════════════════════════════════════════╗\n"))
	sb.WriteString(fmt.Sprintf("║     WhistlerLang Strict Mode Violations      ║\n"))
	sb.WriteString(fmt.Sprintf("╠══════════════════════════════════════════════╣\n"))
	for i, e := range errs {
		sb.WriteString(fmt.Sprintf("║  [%d] %s\n", i+1, e.Message))
		sb.WriteString(fmt.Sprintf("║\n"))
	}
	sb.WriteString(fmt.Sprintf("╠══════════════════════════════════════════════╣\n"))
	sb.WriteString(fmt.Sprintf("║  %d error(s) found. Fix them or use          ║\n", len(errs)))
	sb.WriteString(fmt.Sprintf("║  _knownuse { } to bypass a specific block.   ║\n"))
	sb.WriteString(fmt.Sprintf("╚══════════════════════════════════════════════╝\n"))
	return sb.String()
}
