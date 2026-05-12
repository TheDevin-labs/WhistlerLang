package main

type Node interface {
	nodeType() string
}

type Program struct {
	Statements []Node
}

func (p *Program) nodeType() string { return "Program" }

type LetStatement struct {
	Name  string
	Value Node
}

func (l *LetStatement) nodeType() string { return "LetStatement" }

type SayStatement struct {
	Value Node
}

func (s *SayStatement) nodeType() string { return "SayStatement" }

type FnStatement struct {
	Name   string
	Params []string
	Body   []Node
}

func (f *FnStatement) nodeType() string { return "FnStatement" }

type ReturnStatement struct {
	Value Node
}

func (r *ReturnStatement) nodeType() string { return "ReturnStatement" }

type IfStatement struct {
	Condition Node
	Then      []Node
	Elifs     []ElifClause
	Else      []Node
}

type ElifClause struct {
	Condition Node
	Body      []Node
}

func (i *IfStatement) nodeType() string { return "IfStatement" }

type ForStatement struct {
	Variable string
	Iterable Node
	Body     []Node
}

func (f *ForStatement) nodeType() string { return "ForStatement" }

type ExpressionStatement struct {
	Expr Node
}

func (e *ExpressionStatement) nodeType() string { return "ExpressionStatement" }

type IntLiteral struct {
	Value int64
}

func (i *IntLiteral) nodeType() string { return "IntLiteral" }

type FloatLiteral struct {
	Value float64
}

func (f *FloatLiteral) nodeType() string { return "FloatLiteral" }

type ComplexLiteral struct {
	Real float64
	Imag float64
}

func (c *ComplexLiteral) nodeType() string { return "ComplexLiteral" }

type BoolLiteral struct {
	Value bool
}

func (b *BoolLiteral) nodeType() string { return "BoolLiteral" }

type StringLiteral struct {
	Value string
}

func (s *StringLiteral) nodeType() string { return "StringLiteral" }

type Identifier struct {
	Name string
}

func (i *Identifier) nodeType() string { return "Identifier" }

type BinaryExpr struct {
	Left     Node
	Operator string
	Right    Node
}

func (b *BinaryExpr) nodeType() string { return "BinaryExpr" }

type UnaryExpr struct {
	Operator string
	Operand  Node
}

func (u *UnaryExpr) nodeType() string { return "UnaryExpr" }

type ArrayLiteral struct {
	Elements []Node
}

func (a *ArrayLiteral) nodeType() string { return "ArrayLiteral" }

type MatrixLiteral struct {
	Rows [][]Node
}

func (m *MatrixLiteral) nodeType() string { return "MatrixLiteral" }

type IndexExpr struct {
	Object Node
	Index  Node
}

func (i *IndexExpr) nodeType() string { return "IndexExpr" }

type CallExpr struct {
	Function string
	Args     []Node
}

func (c *CallExpr) nodeType() string { return "CallExpr" }

type RangeExpr struct {
	End Node
}

func (r *RangeExpr) nodeType() string { return "RangeExpr" }

