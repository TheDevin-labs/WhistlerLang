package main

import (
	"fmt"
	"strings"
)

type arrayInfo struct {
	ptr  string
	size int
}

type CodeGen struct {
	tmpCount   int
	labelCount int
	vars       map[string]string
	arrays     map[string]arrayInfo
	fns        map[string]*FnStatement
	strLits    []string
	strSeen    map[string]bool
	triple     string
}

func NewCodeGen(triple string) *CodeGen {
	return &CodeGen{
		vars:    make(map[string]string),
		arrays:  make(map[string]arrayInfo),
		fns:     make(map[string]*FnStatement),
		strSeen: make(map[string]bool),
		triple:  triple,
	}
}

func (c *CodeGen) tmp() string { c.tmpCount++; return fmt.Sprintf("%%t%d", c.tmpCount) }
func (c *CodeGen) lbl() string { c.labelCount++; return fmt.Sprintf("L%d", c.labelCount) }

func (c *CodeGen) addStr(s string) int {
	if !c.strSeen[s] { c.strLits = append(c.strLits, s); c.strSeen[s] = true }
	for i, v := range c.strLits { if v == s { return i } }
	return 0
}

func (c *CodeGen) collectStrings(n Node) {
	switch v := n.(type) {
	case *Program:             for _, s := range v.Statements { c.collectStrings(s) }
	case *FnStatement:         for _, s := range v.Body { c.collectStrings(s) }
	case *LetStatement:        c.collectStrings(v.Value)
	case *SayStatement:        c.collectStrings(v.Value)
	case *StringLiteral:       c.addStr(v.Value)
	case *IfStatement:
		for _, s := range v.Then { c.collectStrings(s) }
		for _, e := range v.Elifs { for _, s := range e.Body { c.collectStrings(s) } }
		for _, s := range v.Else { c.collectStrings(s) }
	case *ForStatement:        for _, s := range v.Body { c.collectStrings(s) }
	case *BlockrockStatement:  for _, s := range v.Body { c.collectStrings(s) }
	case *KnownUseStatement:   for _, s := range v.Body { c.collectStrings(s) }
	case *ExpressionStatement: c.collectStrings(v.Expr)
	case *BinaryExpr:          c.collectStrings(v.Left); c.collectStrings(v.Right)
	case *CallExpr:            for _, a := range v.Args { c.collectStrings(a) }
	}
}

func (c *CodeGen) Generate(program *Program) (string, error) {
	c.collectStrings(program)
	for _, s := range program.Statements {
		if fn, ok := s.(*FnStatement); ok { c.fns[fn.Name] = fn }
	}

	var ir strings.Builder
	ir.WriteString("; WhistlerLang 26.0\n")
	ir.WriteString(fmt.Sprintf("target triple = \"%s\"\n\n", c.triple))
	ir.WriteString("declare i32 @printf(i8* nocapture, ...)\n")
	ir.WriteString("declare i32 @puts(i8* nocapture)\n")
	ir.WriteString("declare double @sin(double)\n")
	ir.WriteString("declare double @cos(double)\n")
	ir.WriteString("declare double @sqrt(double)\n")
	ir.WriteString("declare double @log(double)\n")
	ir.WriteString("declare double @exp(double)\n")
	ir.WriteString("declare double @pow(double, double)\n")
	ir.WriteString("declare double @fabs(double)\n")
	ir.WriteString("declare double @ceil(double)\n")
	ir.WriteString("declare double @floor(double)\n")
	ir.WriteString("declare double @round(double)\n\n")
	ir.WriteString("@fmt_f  = private constant [4 x i8] c\"%g\\0A\\00\"\n")
	ir.WriteString("@fmt_lf = private constant [3 x i8] c\"%g\\00\"\n")
	ir.WriteString("@str_lb = private constant [2 x i8] c\"[\\00\"\n")
	ir.WriteString("@str_rb = private constant [3 x i8] c\"]\\0A\\00\"\n")
	ir.WriteString("@str_cm = private constant [3 x i8] c\", \\00\"\n\n")
	for i, s := range c.strLits {
		escaped := escapeStr(s)
		length := len(s) + 2
		ir.WriteString(fmt.Sprintf("@.s%d = private constant [%d x i8] c\"%s\\0A\\00\"\n", i, length, escaped))
	}
	ir.WriteString("\n")
	for _, s := range program.Statements {
		if fn, ok := s.(*FnStatement); ok {
			fnIR, err := c.genFunction(fn)
			if err != nil { return "", err }
			ir.WriteString(fnIR)
		}
	}
	ir.WriteString("define i32 @main() {\nentry:\n")
	for _, s := range program.Statements {
		if _, ok := s.(*FnStatement); ok { continue }
		code, err := c.genStmt(s)
		if err != nil { return "", err }
		ir.WriteString(code)
	}
	ir.WriteString("  ret i32 0\n}\n")
	return ir.String(), nil
}

func (c *CodeGen) genFunction(fn *FnStatement) (string, error) {
	savedVars   := c.vars
	savedArrays := c.arrays
	savedTmp    := c.tmpCount
	c.vars      = make(map[string]string)
	c.arrays    = make(map[string]arrayInfo)
	c.tmpCount  = 0
	for _, p := range fn.Params { c.vars[p.Name] = "%" + p.Name }
	params := make([]string, len(fn.Params))
	for i, p := range fn.Params { params[i] = "double %" + p.Name }
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("define double @%s(%s) {\nentry:\n", fn.Name, strings.Join(params, ", ")))
	var lastReg string
	for _, s := range fn.Body {
		if ret, ok := s.(*ReturnStatement); ok {
			reg, code, err := c.genExpr(ret.Value)
			if err != nil { return "", err }
			sb.WriteString(code)
			sb.WriteString(fmt.Sprintf("  ret double %s\n", reg))
			goto fnDone
		}
		if expr, ok := s.(*ExpressionStatement); ok {
			reg, code, err := c.genExpr(expr.Expr)
			if err != nil { return "", err }
			sb.WriteString(code)
			lastReg = reg
			continue
		}
		code, err := c.genStmt(s)
		if err != nil { return "", err }
		sb.WriteString(code)
	}
	if lastReg != "" {
		sb.WriteString(fmt.Sprintf("  ret double %s\n", lastReg))
	} else {
		sb.WriteString("  ret double 0.0\n")
	}
fnDone:
	sb.WriteString("}\n\n")
	c.vars      = savedVars
	c.arrays    = savedArrays
	c.tmpCount  = savedTmp
	return sb.String(), nil
}

func (c *CodeGen) genStmt(node Node) (string, error) {
	var sb strings.Builder
	switch n := node.(type) {

	case *LetStatement:
		if arr, ok := n.Value.(*ArrayLiteral); ok {
			size := len(arr.Elements)
			ptr := c.tmp()
			sb.WriteString(fmt.Sprintf("  %s = alloca [%d x double]\n", ptr, size))
			for i, el := range arr.Elements {
				reg, code, err := c.genExpr(el)
				if err != nil { return "", err }
				sb.WriteString(code)
				gep := c.tmp()
				sb.WriteString(fmt.Sprintf("  %s = getelementptr [%d x double], [%d x double]* %s, i32 0, i32 %d\n", gep, size, size, ptr, i))
				sb.WriteString(fmt.Sprintf("  store double %s, double* %s\n", reg, gep))
			}
			c.arrays[n.Name] = arrayInfo{ptr: ptr, size: size}
			c.vars[n.Name] = ptr
			return sb.String(), nil
		}
		if mat, ok := n.Value.(*MatrixLiteral); ok {
			total := 0
			for _, row := range mat.Rows { total += len(row) }
			ptr := c.tmp()
			sb.WriteString(fmt.Sprintf("  %s = alloca [%d x double]\n", ptr, total))
			idx := 0
			for _, row := range mat.Rows {
				for _, el := range row {
					reg, code, err := c.genExpr(el)
					if err != nil { return "", err }
					sb.WriteString(code)
					gep := c.tmp()
					sb.WriteString(fmt.Sprintf("  %s = getelementptr [%d x double], [%d x double]* %s, i32 0, i32 %d\n", gep, total, total, ptr, idx))
					sb.WriteString(fmt.Sprintf("  store double %s, double* %s\n", reg, gep))
					idx++
				}
			}
			c.arrays[n.Name] = arrayInfo{ptr: ptr, size: total}
			c.vars[n.Name] = ptr
			return sb.String(), nil
		}
		reg, code, err := c.genExpr(n.Value)
		if err != nil { return "", err }
		sb.WriteString(code)
		c.vars[n.Name] = reg

	case *SayStatement:
		if v, ok := n.Value.(*StringLiteral); ok {
			idx := c.addStr(v.Value)
			length := len(v.Value) + 2
			ptr := c.tmp()
			sb.WriteString(fmt.Sprintf("  %s = getelementptr [%d x i8], [%d x i8]* @.s%d, i32 0, i32 0\n", ptr, length, length, idx))
			sb.WriteString(fmt.Sprintf("  call i32 @puts(i8* %s)\n", ptr))
			return sb.String(), nil
		}
		if ident, ok := n.Value.(*Identifier); ok {
			if info, isArr := c.arrays[ident.Name]; isArr {
				lbp := c.tmp()
				sb.WriteString(fmt.Sprintf("  %s = getelementptr [2 x i8], [2 x i8]* @str_lb, i32 0, i32 0\n", lbp))
				sb.WriteString(fmt.Sprintf("  call i32 (i8*, ...) @printf(i8* %s)\n", lbp))
				for i := 0; i < info.size; i++ {
					gep := c.tmp(); val := c.tmp()
					sb.WriteString(fmt.Sprintf("  %s = getelementptr [%d x double], [%d x double]* %s, i32 0, i32 %d\n", gep, info.size, info.size, info.ptr, i))
					sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", val, gep))
					fmtp := c.tmp()
					sb.WriteString(fmt.Sprintf("  %s = getelementptr [3 x i8], [3 x i8]* @fmt_lf, i32 0, i32 0\n", fmtp))
					sb.WriteString(fmt.Sprintf("  call i32 (i8*, ...) @printf(i8* %s, double %s)\n", fmtp, val))
					if i < info.size-1 {
						cmp := c.tmp()
						sb.WriteString(fmt.Sprintf("  %s = getelementptr [3 x i8], [3 x i8]* @str_cm, i32 0, i32 0\n", cmp))
						sb.WriteString(fmt.Sprintf("  call i32 (i8*, ...) @printf(i8* %s)\n", cmp))
					}
				}
				rbp := c.tmp()
				sb.WriteString(fmt.Sprintf("  %s = getelementptr [3 x i8], [3 x i8]* @str_rb, i32 0, i32 0\n", rbp))
				sb.WriteString(fmt.Sprintf("  call i32 (i8*, ...) @printf(i8* %s)\n", rbp))
				return sb.String(), nil
			}
		}
		reg, code, err := c.genExpr(n.Value)
		if err != nil { return "", err }
		sb.WriteString(code)
		fptr := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = getelementptr [4 x i8], [4 x i8]* @fmt_f, i32 0, i32 0\n", fptr))
		sb.WriteString(fmt.Sprintf("  call i32 (i8*, ...) @printf(i8* %s, double %s)\n", fptr, reg))

	case *ExpressionStatement:
		_, code, err := c.genExpr(n.Expr)
		if err != nil { return "", err }
		sb.WriteString(code)

	case *ReturnStatement:
		reg, code, err := c.genExpr(n.Value)
		if err != nil { return "", err }
		sb.WriteString(code)
		sb.WriteString(fmt.Sprintf("  ret double %s\n", reg))

	case *KnownUseStatement:
		for _, s := range n.Body {
			code, err := c.genStmt(s)
			if err != nil { return "", err }
			sb.WriteString(code)
		}

	case *BlockrockStatement:
		for _, s := range n.Body {
			code, err := c.genStmt(s)
			if err != nil {
				for _, ps := range n.PanicBody {
					pc, _ := c.genStmt(ps)
					sb.WriteString(pc)
				}
				return sb.String(), nil
			}
			sb.WriteString(code)
		}

	case *IfStatement:
		condReg, condCode, err := c.genExpr(n.Condition)
		if err != nil { return "", err }
		sb.WriteString(condCode)
		thenL := c.lbl(); endL := c.lbl()
		elifCondL := make([]string, len(n.Elifs))
		elifBodyL := make([]string, len(n.Elifs))
		for i := range n.Elifs { elifCondL[i] = c.lbl(); elifBodyL[i] = c.lbl() }
		elseL := c.lbl()
		firstFalse := elseL
		if len(n.Elifs) > 0 { firstFalse = elifCondL[0] }
		cmp := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = fcmp one double %s, 0.0\n", cmp, condReg))
		sb.WriteString(fmt.Sprintf("  br i1 %s, label %%%s, label %%%s\n", cmp, thenL, firstFalse))
		sb.WriteString(thenL + ":\n")
		for _, s := range n.Then {
			code, err := c.genStmt(s)
			if err != nil { return "", err }
			sb.WriteString(code)
		}
		sb.WriteString(fmt.Sprintf("  br label %%%s\n", endL))
		for i, elif := range n.Elifs {
			sb.WriteString(elifCondL[i] + ":\n")
			er, ec, err := c.genExpr(elif.Condition)
			if err != nil { return "", err }
			sb.WriteString(ec)
			nextFalse := elseL
			if i+1 < len(elifCondL) { nextFalse = elifCondL[i+1] }
			ecmp := c.tmp()
			sb.WriteString(fmt.Sprintf("  %s = fcmp one double %s, 0.0\n", ecmp, er))
			sb.WriteString(fmt.Sprintf("  br i1 %s, label %%%s, label %%%s\n", ecmp, elifBodyL[i], nextFalse))
			sb.WriteString(elifBodyL[i] + ":\n")
			for _, s := range elif.Body {
				code, err := c.genStmt(s)
				if err != nil { return "", err }
				sb.WriteString(code)
			}
			sb.WriteString(fmt.Sprintf("  br label %%%s\n", endL))
		}
		sb.WriteString(elseL + ":\n")
		for _, s := range n.Else {
			code, err := c.genStmt(s)
			if err != nil { return "", err }
			sb.WriteString(code)
		}
		sb.WriteString(fmt.Sprintf("  br label %%%s\n", endL))
		sb.WriteString(endL + ":\n")

	case *ForStatement:
		if rng, ok := n.Iterable.(*RangeExpr); ok {
			endReg, endCode, err := c.genExpr(rng.End)
			if err != nil { return "", err }
			sb.WriteString(endCode)
			iPtr := c.tmp()
			sb.WriteString(fmt.Sprintf("  %s = alloca double\n", iPtr))
			sb.WriteString(fmt.Sprintf("  store double 0.0, double* %s\n", iPtr))
			loopL := c.lbl(); bodyL := c.lbl(); afterL := c.lbl()
			sb.WriteString(fmt.Sprintf("  br label %%%s\n", loopL))
			sb.WriteString(loopL + ":\n")
			iVal := c.tmp()
			sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", iVal, iPtr))
			loopcmp := c.tmp()
			sb.WriteString(fmt.Sprintf("  %s = fcmp olt double %s, %s\n", loopcmp, iVal, endReg))
			sb.WriteString(fmt.Sprintf("  br i1 %s, label %%%s, label %%%s\n", loopcmp, bodyL, afterL))
			sb.WriteString(bodyL + ":\n")
			c.vars[n.Variable] = iVal
			for _, s := range n.Body {
				code, err := c.genStmt(s)
				if err != nil { return "", err }
				sb.WriteString(code)
			}
			iVal2 := c.tmp(); inc := c.tmp()
			sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", iVal2, iPtr))
			sb.WriteString(fmt.Sprintf("  %s = fadd double %s, 1.0\n", inc, iVal2))
			sb.WriteString(fmt.Sprintf("  store double %s, double* %s\n", inc, iPtr))
			sb.WriteString(fmt.Sprintf("  br label %%%s\n", loopL))
			sb.WriteString(afterL + ":\n")
		}
	}

	return sb.String(), nil
}

func (c *CodeGen) genExpr(node Node) (string, string, error) {
	var sb strings.Builder
	switch n := node.(type) {

	case *IntLiteral:
		r := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = sitofp i64 %d to double\n", r, n.Value))
		return r, sb.String(), nil

	case *FloatLiteral:
		r := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = fadd double 0.0, %f\n", r, n.Value))
		return r, sb.String(), nil

	case *BoolLiteral:
		r := c.tmp()
		v := 0.0
		if n.Value { v = 1.0 }
		sb.WriteString(fmt.Sprintf("  %s = fadd double 0.0, %f\n", r, v))
		return r, sb.String(), nil

	case *ByteLiteral:
		r := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = fadd double 0.0, %f\n", r, float64(n.Value)))
		return r, sb.String(), nil

	case *BytesLiteral:
		r := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = fadd double 0.0, 0.0\n", r))
		return r, sb.String(), nil

	case *StringLiteral:
		r := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = fadd double 0.0, 0.0\n", r))
		return r, sb.String(), nil

	case *Identifier:
		if reg, ok := c.vars[n.Name]; ok { return reg, "", nil }
		return "", "", fmt.Errorf("undefined variable: %s", n.Name)

	case *UnaryExpr:
		reg, code, err := c.genExpr(n.Operand)
		if err != nil { return "", "", err }
		sb.WriteString(code)
		r := c.tmp()
		if n.Operator == "-" {
			sb.WriteString(fmt.Sprintf("  %s = fneg double %s\n", r, reg))
		} else {
			cmp := c.tmp()
			sb.WriteString(fmt.Sprintf("  %s = fcmp oeq double %s, 0.0\n", cmp, reg))
			sb.WriteString(fmt.Sprintf("  %s = uitofp i1 %s to double\n", r, cmp))
		}
		return r, sb.String(), nil

	case *BinaryExpr:
		lr, lc, err := c.genExpr(n.Left)
		if err != nil { return "", "", err }
		rr, rc, err := c.genExpr(n.Right)
		if err != nil { return "", "", err }
		sb.WriteString(lc); sb.WriteString(rc)
		r := c.tmp()
		switch n.Operator {
		case "+":  sb.WriteString(fmt.Sprintf("  %s = fadd double %s, %s\n", r, lr, rr))
		case "-":  sb.WriteString(fmt.Sprintf("  %s = fsub double %s, %s\n", r, lr, rr))
		case "*":  sb.WriteString(fmt.Sprintf("  %s = fmul double %s, %s\n", r, lr, rr))
		case "/":  sb.WriteString(fmt.Sprintf("  %s = fdiv double %s, %s\n", r, lr, rr))
		case "==": cmp := c.tmp(); sb.WriteString(fmt.Sprintf("  %s = fcmp oeq double %s, %s\n", cmp, lr, rr)); sb.WriteString(fmt.Sprintf("  %s = uitofp i1 %s to double\n", r, cmp))
		case "!=": cmp := c.tmp(); sb.WriteString(fmt.Sprintf("  %s = fcmp one double %s, %s\n", cmp, lr, rr)); sb.WriteString(fmt.Sprintf("  %s = uitofp i1 %s to double\n", r, cmp))
		case "<":  cmp := c.tmp(); sb.WriteString(fmt.Sprintf("  %s = fcmp olt double %s, %s\n", cmp, lr, rr)); sb.WriteString(fmt.Sprintf("  %s = uitofp i1 %s to double\n", r, cmp))
		case ">":  cmp := c.tmp(); sb.WriteString(fmt.Sprintf("  %s = fcmp ogt double %s, %s\n", cmp, lr, rr)); sb.WriteString(fmt.Sprintf("  %s = uitofp i1 %s to double\n", r, cmp))
		case "<=": cmp := c.tmp(); sb.WriteString(fmt.Sprintf("  %s = fcmp ole double %s, %s\n", cmp, lr, rr)); sb.WriteString(fmt.Sprintf("  %s = uitofp i1 %s to double\n", r, cmp))
		case ">=": cmp := c.tmp(); sb.WriteString(fmt.Sprintf("  %s = fcmp oge double %s, %s\n", cmp, lr, rr)); sb.WriteString(fmt.Sprintf("  %s = uitofp i1 %s to double\n", r, cmp))
		default:   return "", "", fmt.Errorf("unknown operator: %s", n.Operator)
		}
		return r, sb.String(), nil

	case *IndexExpr:
		if ident, ok := n.Object.(*Identifier); ok {
			if info, isArr := c.arrays[ident.Name]; isArr {
				idxReg, idxCode, err := c.genExpr(n.Index)
				if err != nil { return "", "", err }
				sb.WriteString(idxCode)
				i32reg := c.tmp()
				sb.WriteString(fmt.Sprintf("  %s = fptosi double %s to i32\n", i32reg, idxReg))
				gep := c.tmp()
				sb.WriteString(fmt.Sprintf("  %s = getelementptr [%d x double], [%d x double]* %s, i32 0, i32 %s\n", gep, info.size, info.size, info.ptr, i32reg))
				val := c.tmp()
				sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", val, gep))
				return val, sb.String(), nil
			}
		}
		return "", "", fmt.Errorf("index on non-array")

	case *CallExpr:
		switch n.Function {
		case "sin", "cos", "sqrt", "log", "exp", "ceil", "floor", "round":
			if len(n.Args) != 1 { return "", "", fmt.Errorf("%s() expects 1 argument", n.Function) }
			ar, ac, err := c.genExpr(n.Args[0])
			if err != nil { return "", "", err }
			sb.WriteString(ac)
			r := c.tmp()
			sb.WriteString(fmt.Sprintf("  %s = call double @%s(double %s)\n", r, n.Function, ar))
			return r, sb.String(), nil
		case "abs":
			if len(n.Args) != 1 { return "", "", fmt.Errorf("abs() expects 1 argument") }
			ar, ac, err := c.genExpr(n.Args[0])
			if err != nil { return "", "", err }
			sb.WriteString(ac)
			r := c.tmp()
			sb.WriteString(fmt.Sprintf("  %s = call double @fabs(double %s)\n", r, ar))
			return r, sb.String(), nil
		case "pow":
			if len(n.Args) != 2 { return "", "", fmt.Errorf("pow() expects 2 arguments") }
			ar, ac, err := c.genExpr(n.Args[0])
			if err != nil { return "", "", err }
			br, bc, err := c.genExpr(n.Args[1])
			if err != nil { return "", "", err }
			sb.WriteString(ac); sb.WriteString(bc)
			r := c.tmp()
			sb.WriteString(fmt.Sprintf("  %s = call double @pow(double %s, double %s)\n", r, ar, br))
			return r, sb.String(), nil
		case "sum", "mean", "min", "max", "std", "variance", "len":
			if len(n.Args) != 1 { return "", "", fmt.Errorf("%s() expects 1 argument", n.Function) }
			info, err := c.requireArray(n.Args[0])
			if err != nil { return "", "", err }
			return c.genStatsIR(&sb, n.Function, info)
		default:
			var argRegs []string
			for _, a := range n.Args {
				reg, code, err := c.genExpr(a)
				if err != nil { return "", "", err }
				sb.WriteString(code)
				argRegs = append(argRegs, "double "+reg)
			}
			r := c.tmp()
			sb.WriteString(fmt.Sprintf("  %s = call double @%s(%s)\n", r, n.Function, strings.Join(argRegs, ", ")))
			return r, sb.String(), nil
		}
	}

	r := c.tmp()
	sb.WriteString(fmt.Sprintf("  %s = fadd double 0.0, 0.0\n", r))
	return r, sb.String(), nil
}

func (c *CodeGen) requireArray(node Node) (arrayInfo, error) {
	if ident, ok := node.(*Identifier); ok {
		if info, ok := c.arrays[ident.Name]; ok { return info, nil }
		return arrayInfo{}, fmt.Errorf("'%s' is not an array", ident.Name)
	}
	return arrayInfo{}, fmt.Errorf("expected an array variable")
}

func (c *CodeGen) genSumIR(sb *strings.Builder, ptr string, size int) (string, error) {
	sumPtr := c.tmp(); iPtr := c.tmp()
	sb.WriteString(fmt.Sprintf("  %s = alloca double\n", sumPtr))
	sb.WriteString(fmt.Sprintf("  store double 0.0, double* %s\n", sumPtr))
	sb.WriteString(fmt.Sprintf("  %s = alloca double\n", iPtr))
	sb.WriteString(fmt.Sprintf("  store double 0.0, double* %s\n", iPtr))
	loopL := c.lbl(); bodyL := c.lbl(); afterL := c.lbl()
	sb.WriteString(fmt.Sprintf("  br label %%%s\n", loopL))
	sb.WriteString(loopL + ":\n")
	iVal := c.tmp(); cmp := c.tmp()
	sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", iVal, iPtr))
	sb.WriteString(fmt.Sprintf("  %s = fcmp olt double %s, %f\n", cmp, iVal, float64(size)))
	sb.WriteString(fmt.Sprintf("  br i1 %s, label %%%s, label %%%s\n", cmp, bodyL, afterL))
	sb.WriteString(bodyL + ":\n")
	i32 := c.tmp(); gep := c.tmp(); el := c.tmp(); cur := c.tmp(); newS := c.tmp()
	sb.WriteString(fmt.Sprintf("  %s = fptosi double %s to i32\n", i32, iVal))
	sb.WriteString(fmt.Sprintf("  %s = getelementptr [%d x double], [%d x double]* %s, i32 0, i32 %s\n", gep, size, size, ptr, i32))
	sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", el, gep))
	sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", cur, sumPtr))
	sb.WriteString(fmt.Sprintf("  %s = fadd double %s, %s\n", newS, cur, el))
	sb.WriteString(fmt.Sprintf("  store double %s, double* %s\n", newS, sumPtr))
	iVal2 := c.tmp(); inc := c.tmp()
	sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", iVal2, iPtr))
	sb.WriteString(fmt.Sprintf("  %s = fadd double %s, 1.0\n", inc, iVal2))
	sb.WriteString(fmt.Sprintf("  store double %s, double* %s\n", inc, iPtr))
	sb.WriteString(fmt.Sprintf("  br label %%%s\n", loopL))
	sb.WriteString(afterL + ":\n")
	r := c.tmp()
	sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", r, sumPtr))
	return r, nil
}

func (c *CodeGen) genStatsIR(sb *strings.Builder, fn string, info arrayInfo) (string, string, error) {
	size := info.size; ptr := info.ptr
	switch fn {
	case "len":
		r := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = fadd double 0.0, %f\n", r, float64(size)))
		return r, sb.String(), nil
	case "sum":
		r, err := c.genSumIR(sb, ptr, size)
		return r, sb.String(), err
	case "mean":
		sumReg, err := c.genSumIR(sb, ptr, size)
		if err != nil { return "", "", err }
		r := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = fdiv double %s, %f\n", r, sumReg, float64(size)))
		return r, sb.String(), nil
	case "min", "max":
		gep0 := c.tmp(); first := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = getelementptr [%d x double], [%d x double]* %s, i32 0, i32 0\n", gep0, size, size, ptr))
		sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", first, gep0))
		accPtr := c.tmp(); iPtr := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = alloca double\n", accPtr))
		sb.WriteString(fmt.Sprintf("  store double %s, double* %s\n", first, accPtr))
		sb.WriteString(fmt.Sprintf("  %s = alloca double\n", iPtr))
		sb.WriteString(fmt.Sprintf("  store double 1.0, double* %s\n", iPtr))
		loopL := c.lbl(); bodyL := c.lbl(); afterL := c.lbl()
		sb.WriteString(fmt.Sprintf("  br label %%%s\n", loopL))
		sb.WriteString(loopL + ":\n")
		iVal := c.tmp(); cmp := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", iVal, iPtr))
		sb.WriteString(fmt.Sprintf("  %s = fcmp olt double %s, %f\n", cmp, iVal, float64(size)))
		sb.WriteString(fmt.Sprintf("  br i1 %s, label %%%s, label %%%s\n", cmp, bodyL, afterL))
		sb.WriteString(bodyL + ":\n")
		i32 := c.tmp(); gep := c.tmp(); el := c.tmp(); cur := c.tmp(); sel := c.tmp(); pred := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = fptosi double %s to i32\n", i32, iVal))
		sb.WriteString(fmt.Sprintf("  %s = getelementptr [%d x double], [%d x double]* %s, i32 0, i32 %s\n", gep, size, size, ptr, i32))
		sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", el, gep))
		sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", cur, accPtr))
		if fn == "min" {
			sb.WriteString(fmt.Sprintf("  %s = fcmp olt double %s, %s\n", pred, el, cur))
		} else {
			sb.WriteString(fmt.Sprintf("  %s = fcmp ogt double %s, %s\n", pred, el, cur))
		}
		sb.WriteString(fmt.Sprintf("  %s = select i1 %s, double %s, double %s\n", sel, pred, el, cur))
		sb.WriteString(fmt.Sprintf("  store double %s, double* %s\n", sel, accPtr))
		iVal2 := c.tmp(); inc := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", iVal2, iPtr))
		sb.WriteString(fmt.Sprintf("  %s = fadd double %s, 1.0\n", inc, iVal2))
		sb.WriteString(fmt.Sprintf("  store double %s, double* %s\n", inc, iPtr))
		sb.WriteString(fmt.Sprintf("  br label %%%s\n", loopL))
		sb.WriteString(afterL + ":\n")
		r := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", r, accPtr))
		return r, sb.String(), nil
	case "variance":
		sumReg, err := c.genSumIR(sb, ptr, size)
		if err != nil { return "", "", err }
		meanReg := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = fdiv double %s, %f\n", meanReg, sumReg, float64(size)))
		varPtr := c.tmp(); iPtr := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = alloca double\n", varPtr))
		sb.WriteString(fmt.Sprintf("  store double 0.0, double* %s\n", varPtr))
		sb.WriteString(fmt.Sprintf("  %s = alloca double\n", iPtr))
		sb.WriteString(fmt.Sprintf("  store double 0.0, double* %s\n", iPtr))
		loopL := c.lbl(); bodyL := c.lbl(); afterL := c.lbl()
		sb.WriteString(fmt.Sprintf("  br label %%%s\n", loopL))
		sb.WriteString(loopL + ":\n")
		iVal := c.tmp(); cmp := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", iVal, iPtr))
		sb.WriteString(fmt.Sprintf("  %s = fcmp olt double %s, %f\n", cmp, iVal, float64(size)))
		sb.WriteString(fmt.Sprintf("  br i1 %s, label %%%s, label %%%s\n", cmp, bodyL, afterL))
		sb.WriteString(bodyL + ":\n")
		i32 := c.tmp(); gep := c.tmp(); el := c.tmp(); diff := c.tmp(); sq := c.tmp(); cur := c.tmp(); newV := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = fptosi double %s to i32\n", i32, iVal))
		sb.WriteString(fmt.Sprintf("  %s = getelementptr [%d x double], [%d x double]* %s, i32 0, i32 %s\n", gep, size, size, ptr, i32))
		sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", el, gep))
		sb.WriteString(fmt.Sprintf("  %s = fsub double %s, %s\n", diff, el, meanReg))
		sb.WriteString(fmt.Sprintf("  %s = fmul double %s, %s\n", sq, diff, diff))
		sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", cur, varPtr))
		sb.WriteString(fmt.Sprintf("  %s = fadd double %s, %s\n", newV, cur, sq))
		sb.WriteString(fmt.Sprintf("  store double %s, double* %s\n", newV, varPtr))
		iVal2 := c.tmp(); inc := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", iVal2, iPtr))
		sb.WriteString(fmt.Sprintf("  %s = fadd double %s, 1.0\n", inc, iVal2))
		sb.WriteString(fmt.Sprintf("  store double %s, double* %s\n", inc, iPtr))
		sb.WriteString(fmt.Sprintf("  br label %%%s\n", loopL))
		sb.WriteString(afterL + ":\n")
		varSum := c.tmp(); r := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = load double, double* %s\n", varSum, varPtr))
		sb.WriteString(fmt.Sprintf("  %s = fdiv double %s, %f\n", r, varSum, float64(size)))
		return r, sb.String(), nil
	case "std":
		varReg, _, err := c.genStatsIR(sb, "variance", info)
		if err != nil { return "", "", err }
		r := c.tmp()
		sb.WriteString(fmt.Sprintf("  %s = call double @sqrt(double %s)\n", r, varReg))
		return r, sb.String(), nil
	}
	r := c.tmp()
	sb.WriteString(fmt.Sprintf("  %s = fadd double 0.0, 0.0\n", r))
	return r, sb.String(), nil
}

func escapeStr(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\5C")
	s = strings.ReplaceAll(s, "\n", "\\0A")
	s = strings.ReplaceAll(s, "\t", "\\09")
	s = strings.ReplaceAll(s, "\"", "\\22")
	return s
}
