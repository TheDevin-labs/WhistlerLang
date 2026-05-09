package main

type NodeKind int

const (
	NodeUnknown NodeKind = iota
	NodeSay
	NodeAssign
	NodeMath
	NodeMathBlock
	NodeRun
	NodeBuild
	NodeExec
	NodeTimePrint
	NodeTimeSet
	NodeIf
)

// Branch holds one conditional branch (if / elif / else).
type Branch struct {
	Cond     string  // raw condition string; empty for "else"
	Children []*Node // body nodes
}

type Node struct {
	Kind     NodeKind
	Raw      string
	Expr     string
	Lines    []string
	Args     []string
	Left     string
	Right    string
	Children []*Node  // kept for NodeIf "if" branch body (first branch)
	Branches []Branch // all branches: [if, elif..., else?]
	IsSafe   bool     // for NodeExec: true when exec-safe
}
