package main

import (
	"fmt"
	"strconv"
	"strings"
)

func EvaluateNodes(nodes []*Node, rt *RuntimeEnv) error {
	for _, n := range nodes {
		switch n.Kind {
		case NodeSay:
			res := evalSayConcat(n.Expr, rt)
			fmt.Println(res)
			rt.mu.Lock()
			rt.Variables["last_say"] = Value{Kind: ValString, Str: res}
			rt.mu.Unlock()

		case NodeAssign:
			handleLet(n.Expr, rt)

		case NodeMath:
			v, err := EvalExpressionWithVars(n.Expr, rt.Variables)
			if err != nil {
				fmt.Println("Math error:", err)
			} else {
				printNumber(v)
			}

		case NodeMathBlock:
			for _, l := range n.Lines {
				if strings.TrimSpace(l) == "" {
					continue
				}
				v, err := EvalExpressionWithVars(l, rt.Variables)
				if err != nil {
					fmt.Println("Math error:", err)
				} else {
					printNumber(v)
				}
			}

		case NodeRun:
			if len(n.Args) >= 1 {
				_ = rt.RunScript(n.Args[0])
			}

		case NodeBuild:
			if len(n.Args) >= 1 {
				outPath, arch := "", ""
				if len(n.Args) >= 2 {
					outPath = n.Args[1]
				}
				if len(n.Args) >= 3 {
					arch = n.Args[2]
				}
				res, err := rt.BuildScript(n.Args[0], outPath, arch)
				if err != nil {
					fmt.Println("build error:", err)
				} else {
					fmt.Println(res)
				}
			}

		case NodeExec:
			if len(n.Args) >= 1 {
				_, _ = rt.ExecShell(strings.Join(n.Args, " "), n.IsSafe)
			}

		case NodeTimePrint:
			fmt.Println(PrintTime())

		case NodeTimeSet:
			if len(n.Args) >= 2 {
				SetGlobalTime(trimQuotes(n.Args[0]), trimQuotes(n.Args[1]))
			}

		case NodeIf:
			// Walk branches: if, elif..., else
			for _, branch := range n.Branches {
				// Empty cond = else clause — always runs
				if branch.Cond == "" {
					_ = EvaluateNodes(branch.Children, rt)
					break
				}
				if evalCondition(branch.Cond, rt) {
					_ = EvaluateNodes(branch.Children, rt)
					break
				}
			}

		default:
			line := strings.TrimSpace(n.Raw)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			fmt.Println("Unknown command:", n.Raw)
		}
	}
	return nil
}

func printNumber(v float64) {
	if v == float64(int64(v)) {
		fmt.Printf("%.0f\n", v)
	} else {
		fmt.Printf("%v\n", v)
	}
}

// evalSayConcat resolves a say argument that may contain string literals and
// variables joined by '+'.
func evalSayConcat(arg string, rt *RuntimeEnv) string {
	parts := splitByPlus(arg)
	out := strings.Builder{}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		// Quoted string literal
		if len(p) >= 2 && ((p[0] == '"' && p[len(p)-1] == '"') || (p[0] == '\'' && p[len(p)-1] == '\'')) {
			out.WriteString(trimQuotes(p))
			continue
		}
		// Variable lookup
		rt.mu.RLock()
		val, ok := rt.Variables[p]
		rt.mu.RUnlock()
		if ok {
			switch val.Kind {
			case ValString:
				out.WriteString(val.Str)
			case ValNumber:
				out.WriteString(strconv.FormatFloat(val.Num, 'f', -1, 64))
			}
			continue
		}
		// Try as math expression
		v, err := EvalExpressionWithVars(p, rt.Variables)
		if err == nil {
			out.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
			continue
		}
		// Fallback: emit as-is
		out.WriteString(p)
	}
	return out.String()
}

func splitByPlus(s string) []string {
	parts := []string{}
	cur := strings.Builder{}
	inQuotes := false
	var q rune
	for _, r := range s {
		if inQuotes {
			cur.WriteRune(r)
			if r == q {
				inQuotes = false
			}
			continue
		}
		if r == '"' || r == '\'' {
			inQuotes = true
			q = r
			cur.WriteRune(r)
			continue
		}
		if r == '+' {
			parts = append(parts, cur.String())
			cur.Reset()
			continue
		}
		cur.WriteRune(r)
	}
	if cur.Len() > 0 {
		parts = append(parts, cur.String())
	}
	return parts
}

func handleLet(expr string, rt *RuntimeEnv) {
	rest := strings.TrimSpace(expr)
	idx := strings.Index(rest, "=")
	if idx < 0 {
		return
	}
	name := strings.TrimSpace(rest[:idx])
	right := strings.TrimSpace(rest[idx+1:])
	if name == "" || right == "" {
		return
	}
	// String literal
	if len(right) >= 2 && ((right[0] == '"' && right[len(right)-1] == '"') || (right[0] == '\'' && right[len(right)-1] == '\'')) {
		rt.mu.Lock()
		rt.Variables[name] = Value{Kind: ValString, Str: trimQuotes(right)}
		rt.mu.Unlock()
		return
	}
	// Numeric / expression
	if v, err := EvalExpressionWithVars(right, rt.Variables); err == nil {
		rt.mu.Lock()
		rt.Variables[name] = Value{Kind: ValNumber, Num: v}
		rt.mu.Unlock()
		return
	}
	// Bare word → store as string
	rt.mu.Lock()
	rt.Variables[name] = Value{Kind: ValString, Str: right}
	rt.mu.Unlock()
}

// evalCondition evaluates a condition string like "score >= 90" or "user == \"Alice\"".
// It tries operators longest-first to avoid ">=" being matched as ">".
func evalCondition(cond string, rt *RuntimeEnv) bool {
	ops := []string{">=", "<=", "==", "!=", ">", "<", "="}
	for _, op := range ops {
		idx := strings.Index(cond, op)
		if idx < 0 {
			continue
		}
		left := strings.TrimSpace(cond[:idx])
		right := strings.TrimSpace(cond[idx+len(op):])

		// String comparison
		leftIsStr := len(left) > 0 && (left[0] == '"' || left[0] == '\'')
		rightIsStr := len(right) > 0 && (right[0] == '"' || right[0] == '\'')

		// Resolve left as variable string if not quoted
		if !leftIsStr {
			rt.mu.RLock()
			if v, ok := rt.Variables[left]; ok && v.Kind == ValString {
				rt.mu.RUnlock()
				l := v.Str
				r := trimQuotes(right)
				switch op {
				case "==", "=":
					return l == r
				case "!=":
					return l != r
				}
				return false
			}
			rt.mu.RUnlock()
		}

		if leftIsStr || rightIsStr {
			l := trimQuotes(left)
			r := trimQuotes(right)
			switch op {
			case "==", "=":
				return l == r
			case "!=":
				return l != r
			}
			return false
		}

		// Numeric comparison
		ln, lerr := EvalExpressionWithVars(left, rt.Variables)
		rn, rerr := EvalExpressionWithVars(right, rt.Variables)
		if lerr != nil || rerr != nil {
			return false
		}
		switch op {
		case ">=":
			return ln >= rn
		case "<=":
			return ln <= rn
		case ">":
			return ln > rn
		case "<":
			return ln < rn
		case "==", "=":
			return ln == rn
		case "!=":
			return ln != rn
		}
	}
	return false
}
