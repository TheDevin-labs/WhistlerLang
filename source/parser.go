package main

import (
	"os"
	"strings"
)

func ParseFileToNodes(path string) ([]*Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	return parseLineSlice(lines), nil
}

func ParseLine(line string) ([]*Node, error) {
	return ParseLines([]string{line}), nil
}

// ParseLines parses a slice of raw lines into AST nodes.
// Used for if-body parsing (does not support nested if blocks).
func ParseLines(lines []string) []*Node {
	out := []*Node{}
	inMath := false
	var mathLines []string

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if inMath {
			if line == "end" {
				out = append(out, &Node{Kind: NodeMathBlock, Lines: mathLines})
				inMath = false
				mathLines = nil
				continue
			}
			mathLines = append(mathLines, line)
			continue
		}
		if line == "math;" {
			inMath = true
			mathLines = []string{}
			continue
		}
		out = append(out, parseSingleLine(raw, line))
	}
	return out
}

// parseLineSlice is the full top-level parser that handles if/elif/else/end blocks.
func parseLineSlice(lines []string) []*Node {
	out := []*Node{}
	inMath := false
	var mathLines []string

	// if-block state
	type branchState struct {
		cond  string
		lines []string
	}
	inIf := false
	var branches []branchState
	depth := 0 // nesting depth for future nested-if support

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// --- inside a math block ---
		if inMath {
			if line == "end" {
				out = append(out, &Node{Kind: NodeMathBlock, Lines: mathLines})
				inMath = false
				mathLines = nil
				continue
			}
			mathLines = append(mathLines, line)
			continue
		}

		// --- inside an if/elif/else block ---
		if inIf {
			// Track nested if depth so inner "end" doesn't close the outer block
			if strings.HasPrefix(line, "if ") {
				depth++
				branches[len(branches)-1].lines = append(branches[len(branches)-1].lines, raw)
				continue
			}
			if depth > 0 {
				if line == "end" {
					depth--
				}
				branches[len(branches)-1].lines = append(branches[len(branches)-1].lines, raw)
				continue
			}
			if line == "end" {
				// Build the NodeIf from collected branches
				nodeBranches := make([]Branch, 0, len(branches))
				for _, b := range branches {
					children := ParseLines(b.lines)
					nodeBranches = append(nodeBranches, Branch{Cond: b.cond, Children: children})
				}
				n := &Node{
					Kind:     NodeIf,
					Raw:      raw,
					Branches: nodeBranches,
					// Populate legacy Children with the first branch body for any old callers
					Children: nodeBranches[0].Children,
				}
				out = append(out, n)
				inIf = false
				branches = nil
				depth = 0
				continue
			}
			if strings.HasPrefix(line, "elif ") {
				cond := strings.TrimSpace(strings.TrimPrefix(line, "elif"))
				branches = append(branches, branchState{cond: cond})
				continue
			}
			if line == "else" {
				branches = append(branches, branchState{cond: ""}) // empty cond = else
				continue
			}
			// Accumulate into the current (last) branch
			branches[len(branches)-1].lines = append(branches[len(branches)-1].lines, raw)
			continue
		}

		// --- top-level dispatch ---
		if line == "math;" {
			inMath = true
			mathLines = []string{}
			continue
		}
		if strings.HasPrefix(line, "if ") {
			cond := strings.TrimSpace(strings.TrimPrefix(line, "if"))
			inIf = true
			branches = []branchState{{cond: cond}}
			depth = 0
			continue
		}

		out = append(out, parseSingleLine(raw, line))
	}
	return out
}

// parseSingleLine converts one trimmed line into an AST node.
func parseSingleLine(raw, line string) *Node {
	switch {
	case strings.HasPrefix(line, "say "):
		arg := strings.TrimSpace(line[len("say "):])
		return &Node{Kind: NodeSay, Raw: raw, Expr: arg}

	case strings.HasPrefix(line, "let "):
		return &Node{Kind: NodeAssign, Raw: raw, Expr: strings.TrimSpace(line[len("let "):])}

	case strings.HasPrefix(line, "run "):
		arg := trimQuotes(strings.TrimSpace(line[len("run "):]))
		return &Node{Kind: NodeRun, Raw: raw, Args: []string{arg}}

	case strings.HasPrefix(line, "build "):
		parts := splitArgs(line)
		args := []string{}
		for _, p := range parts[1:] {
			args = append(args, trimQuotes(p))
		}
		return &Node{Kind: NodeBuild, Raw: raw, Args: args}

	case strings.HasPrefix(line, "exec-safe "):
		parts := splitArgs(line)
		return &Node{Kind: NodeExec, Raw: raw, Args: parts[1:], IsSafe: true}

	case strings.HasPrefix(line, "exec "):
		parts := splitArgs(line)
		return &Node{Kind: NodeExec, Raw: raw, Args: parts[1:], IsSafe: false}

	case line == "time.print":
		return &Node{Kind: NodeTimePrint, Raw: raw}

	case strings.HasPrefix(line, "time.set"):
		parts := splitArgs(line)
		args := []string{}
		if len(parts) >= 2 {
			args = parts[1:]
		}
		return &Node{Kind: NodeTimeSet, Raw: raw, Args: args}

	case strings.HasPrefix(line, "math() "):
		expr := strings.TrimSpace(line[len("math() "):])
		return &Node{Kind: NodeMath, Raw: raw, Expr: expr}

	default:
		// fallback: treat as math/expression
		return &Node{Kind: NodeMath, Raw: raw, Expr: line}
	}
}
