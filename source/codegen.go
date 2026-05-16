package main

import (
        "fmt"
        "strings"
)

type CodeGen struct {
        output   strings.Builder
        vars     map[string]string
        tmpCount int
        strCount int
        fns      map[string]*FnStatement
}

func NewCodeGen() *CodeGen {
        return &CodeGen{
                vars: make(map[string]string),
                fns:  make(map[string]*FnStatement),
        }
}

func (c *CodeGen) tmp() string {
        c.tmpCount++
        return fmt.Sprintf("%%t%d", c.tmpCount)
}

func (c *CodeGen) strLabel() string {
        c.strCount++
        return fmt.Sprintf("@str%d", c.strCount)
}

func (c *CodeGen) emit(s string) {
        c.output.WriteString(s + "\n")
}

func (c *CodeGen) Generate(program *Program) (string, error) {
        var ir strings.Builder

        ir.WriteString("; WhistlerLang LLVM IR\n")
        ir.WriteString("; Target: arm64-linux / arm64-macos\n\n")
        ir.WriteString("target triple = \"aarch64-unknown-linux-gnu\"\n\n")

        ir.WriteString("declare i32 @printf(i8* nocapture, ...)\n")
        ir.WriteString("declare i32 @puts(i8*)\n")
        ir.WriteString("declare double @sin(double)\n")
        ir.WriteString("declare double @cos(double)\n")
        ir.WriteString("declare double @sqrt(double)\n")
        ir.WriteString("declare double @log(double)\n")
        ir.WriteString("declare double @exp(double)\n")
        ir.WriteString("declare double @pow(double, double)\n")
        ir.WriteString("declare double @fabs(double)\n\n")

        ir.WriteString("@fmt_int    = private constant [4 x i8] c\"%ld\\00\"\n")
        ir.WriteString("@fmt_float  = private constant [4 x i8] c\"%g\\0A\\00\"\n")
        ir.WriteString("@fmt_bool_t = private constant [6 x i8] c\"true\\0A\\00\"\n")
        ir.WriteString("@fmt_bool_f = private constant [7 x i8] c\"false\\0A\\00\"\n\n")

        stringLiterals := collectStrings(program)
        for i, s := range stringLiterals {
                escaped := escapeString(s)
                length := len(s) + 2
                ir.WriteString(fmt.Sprintf("@.str%d = private constant [%d x i8] c\"%s\\0A\\00\"\n", i, length, escaped))
        }
        ir.WriteString("\n")

        for _, stmt := range program.Statements {
                if fn, ok := stmt.(*FnStatement); ok {
                        c.fns[fn.Name] = fn
                }
        }

        for _, stmt := range program.Statements {
                if fn, ok := stmt.(*FnStatement); ok {
                        fnIR, err := c.generateFunction(fn)
                        if err != nil {
                                return "", err
                        }
                        ir.WriteString(fnIR)
                }
        }

        ir.WriteString("define i32 @main() {\n")
        ir.WriteString("entry:\n")

        mainIR, err := c.generateStatements(program.Statements, stringLiterals)
        if err != nil {
                return "", err
        }
        ir.WriteString(mainIR)

        ir.WriteString("  ret i32 0\n")
        ir.WriteString("}\n")

        return ir.String(), nil
}

func (c *CodeGen) generateFunction(fn *FnStatement) (string, error) {
        var sb strings.Builder
        params := make([]string, len(fn.Params))
        for i, p := range fn.Params {
                params[i] = fmt.Sprintf("double %%%s", p)
        }
        sb.WriteString(fmt.Sprintf("define double @%s(%s) {\n", fn.Name, strings.Join(params, ", ")))
        sb.WriteString("entry:\n")

        for _, stmt := range fn.Body {
                line, err := c.generateStmt(stmt, nil)
                if err != nil {
                        return "", err
                }
                sb.WriteString(line)
        }

        sb.WriteString("  ret double 0.0\n")
        sb.WriteString("}\n\n")
        return sb.String(), nil
}

func (c *CodeGen) generateStatements(stmts []Node, stringLiterals []string) (string, error) {
        var sb strings.Builder
        for _, stmt := range stmts {
                if _, ok := stmt.(*FnStatement); ok {
                        continue
                }
                line, err := c.generateStmt(stmt, stringLiterals)
                if err != nil {
                        return "", err
                }
                sb.WriteString(line)
        }
        return sb.String(), nil
}

func (c *CodeGen) generateStmt(node Node, stringLiterals []string) (string, error) {
        var sb strings.Builder

        switch n := node.(type) {
        case *LetStatement:
                reg, code, err := c.generateExpr(n.Value)
                if err != nil {
                        return "", err
                }
                sb.WriteString(code)
                c.vars[n.Name] = reg

        case *SayStatement:
                switch val := n.Value.(type) {
                case *StringLiteral:
                        idx := findStringIndex(stringLiterals, val.Value)
                        length := len(val.Value) + 2
                        ptr := c.tmp()
                        sb.WriteString(fmt.Sprintf("  %s = getelementptr [%d x i8], [%d x i8]* @.str%d, i32 0, i32 0\n", ptr, length, length, idx))
                        sb.WriteString(fmt.Sprintf("  call i32 @puts(i8* %s)\n", ptr))
                default:
                        reg, code, err := c.generateExpr(n.Value)
                        if err != nil {
                                return "", err
                        }
                        sb.WriteString(code)
                        ptr := c.tmp()
                        sb.WriteString(fmt.Sprintf("  %s = getelementptr [4 x i8], [4 x i8]* @fmt_float, i32 0, i32 0\n", ptr))
                        sb.WriteString(fmt.Sprintf("  call i32 (i8*, ...) @printf(i8* %s, double %s)\n", ptr, reg))
                }

        case *ExpressionStatement:
                _, code, err := c.generateExpr(n.Expr)
                if err != nil {
                        return "", err
                }
                sb.WriteString(code)

        case *IfStatement:
                condReg, condCode, err := c.generateExpr(n.Condition)
                if err != nil {
                        return "", err
                }
                sb.WriteString(condCode)
                thenLabel := fmt.Sprintf("then%d", c.tmpCount)
                elseLabel := fmt.Sprintf("else%d", c.tmpCount)
                endLabel := fmt.Sprintf("end%d", c.tmpCount)
                c.tmpCount++
                cmpReg := c.tmp()
                sb.WriteString(fmt.Sprintf("  %s = fcmp one double %s, 0.0\n", cmpReg, condReg))
                sb.WriteString(fmt.Sprintf("  br i1 %s, label %%%s, label %%%s\n", cmpReg, thenLabel, elseLabel))
                sb.WriteString(fmt.Sprintf("%s:\n", thenLabel))
                for _, stmt := range n.Then {
                        code, err := c.generateStmt(stmt, stringLiterals)
                        if err != nil {
                                return "", err
                        }
                        sb.WriteString(code)
                }
                sb.WriteString(fmt.Sprintf("  br label %%%s\n", endLabel))
                sb.WriteString(fmt.Sprintf("%s:\n", elseLabel))
                for _, stmt := range n.Else {
                        code, err := c.generateStmt(stmt, stringLiterals)
                        if err != nil {
                                return "", err
                        }
                        sb.WriteString(code)
                }
                sb.WriteString(fmt.Sprintf("  br label %%%s\n", endLabel))
                sb.WriteString(fmt.Sprintf("%s:\n", endLabel))
        }

        return sb.String(), nil
}

func (c *CodeGen) generateExpr(node Node) (string, string, error) {
        var sb strings.Builder

        switch n := node.(type) {
        case *IntLiteral:
                reg := c.tmp()
                sb.WriteString(fmt.Sprintf("  %s = sitofp i64 %d to double\n", reg, n.Value))
                return reg, sb.String(), nil

        case *FloatLiteral:
                reg := c.tmp()
                sb.WriteString(fmt.Sprintf("  %s = fadd double 0.0, %f\n", reg, n.Value))
                return reg, sb.String(), nil

        case *BoolLiteral:
                reg := c.tmp()
                if n.Value {
                        sb.WriteString(fmt.Sprintf("  %s = fadd double 0.0, 1.0\n", reg))
                } else {
                        sb.WriteString(fmt.Sprintf("  %s = fadd double 0.0, 0.0\n", reg))
                }
                return reg, sb.String(), nil

        case *Identifier:
                if reg, ok := c.vars[n.Name]; ok {
                        return reg, "", nil
                }
                return "", "", fmt.Errorf("undefined variable: %s", n.Name)

        case *BinaryExpr:
                leftReg, leftCode, err := c.generateExpr(n.Left)
                if err != nil {
                        return "", "", err
                }
                rightReg, rightCode, err := c.generateExpr(n.Right)
                if err != nil {
                        return "", "", err
                }
                sb.WriteString(leftCode)
                sb.WriteString(rightCode)
                reg := c.tmp()
                switch n.Operator {
                case "+":
                        sb.WriteString(fmt.Sprintf("  %s = fadd double %s, %s\n", reg, leftReg, rightReg))
                case "-":
                        sb.WriteString(fmt.Sprintf("  %s = fsub double %s, %s\n", reg, leftReg, rightReg))
                case "*":
                        sb.WriteString(fmt.Sprintf("  %s = fmul double %s, %s\n", reg, leftReg, rightReg))
                case "/":
                        sb.WriteString(fmt.Sprintf("  %s = fdiv double %s, %s\n", reg, leftReg, rightReg))
                case "==":
                        cmp := c.tmp()
                        sb.WriteString(fmt.Sprintf("  %s = fcmp oeq double %s, %s\n", cmp, leftReg, rightReg))
                        sb.WriteString(fmt.Sprintf("  %s = uitofp i1 %s to double\n", reg, cmp))
                case "!=":
                        cmp := c.tmp()
                        sb.WriteString(fmt.Sprintf("  %s = fcmp one double %s, %s\n", cmp, leftReg, rightReg))
                        sb.WriteString(fmt.Sprintf("  %s = uitofp i1 %s to double\n", reg, cmp))
                case "<":
                        cmp := c.tmp()
                        sb.WriteString(fmt.Sprintf("  %s = fcmp olt double %s, %s\n", cmp, leftReg, rightReg))
                        sb.WriteString(fmt.Sprintf("  %s = uitofp i1 %s to double\n", reg, cmp))
                case ">":
                        cmp := c.tmp()
                        sb.WriteString(fmt.Sprintf("  %s = fcmp ogt double %s, %s\n", cmp, leftReg, rightReg))
                        sb.WriteString(fmt.Sprintf("  %s = uitofp i1 %s to double\n", reg, cmp))
                case "<=":
                        cmp := c.tmp()
                        sb.WriteString(fmt.Sprintf("  %s = fcmp ole double %s, %s\n", cmp, leftReg, rightReg))
                        sb.WriteString(fmt.Sprintf("  %s = uitofp i1 %s to double\n", reg, cmp))
                case ">=":
                        cmp := c.tmp()
                        sb.WriteString(fmt.Sprintf("  %s = fcmp oge double %s, %s\n", cmp, leftReg, rightReg))
                        sb.WriteString(fmt.Sprintf("  %s = uitofp i1 %s to double\n", reg, cmp))
                default:
                        return "", "", fmt.Errorf("unknown operator: %s", n.Operator)
                }
                return reg, sb.String(), nil

        case *CallExpr:
                var argRegs []string
                for _, arg := range n.Args {
                        reg, code, err := c.generateExpr(arg)
                        if err != nil {
                                return "", "", err
                        }
                        sb.WriteString(code)
                        argRegs = append(argRegs, "double "+reg)
                }
                reg := c.tmp()
                switch n.Function {
                case "sin", "cos", "sqrt", "log", "exp":
                        sb.WriteString(fmt.Sprintf("  %s = call double @%s(double %s)\n", reg, n.Function, argRegs[0][7:]))
                case "pow":
                        parts := strings.Split(argRegs[0], " ")
                        parts2 := strings.Split(argRegs[1], " ")
                        sb.WriteString(fmt.Sprintf("  %s = call double @pow(double %s, double %s)\n", reg, parts[1], parts2[1]))
                case "abs":
                        parts := strings.Split(argRegs[0], " ")
                        sb.WriteString(fmt.Sprintf("  %s = call double @fabs(double %s)\n", reg, parts[1]))
                default:
                        argStr := strings.Join(argRegs, ", ")
                        sb.WriteString(fmt.Sprintf("  %s = call double @%s(%s)\n", reg, n.Function, argStr))
                }
                return reg, sb.String(), nil

        case *StringLiteral:
                return "0.0", "", nil
        }

        return "0.0", sb.String(), nil
}

func collectStrings(program *Program) []string {
        var result []string
        seen := map[string]bool{}
        var walk func(node Node)
        walk = func(node Node) {
                switch n := node.(type) {
                case *Program:
                        for _, s := range n.Statements {
                                walk(s)
                        }
                case *SayStatement:
                        walk(n.Value)
                case *StringLiteral:
                        if !seen[n.Value] {
                                result = append(result, n.Value)
                                seen[n.Value] = true
                        }
                case *LetStatement:
                        walk(n.Value)
                case *BinaryExpr:
                        walk(n.Left)
                        walk(n.Right)
                }
        }
        walk(program)
        return result
}

func findStringIndex(strs []string, s string) int {
        for i, v := range strs {
                if v == s {
                        return i
                }
        }
        return 0
}

func escapeString(s string) string {
        s = strings.ReplaceAll(s, "\\", "\\\\")
        s = strings.ReplaceAll(s, "\n", "\\0A")
        s = strings.ReplaceAll(s, "\t", "\\09")
        s = strings.ReplaceAll(s, "\"", "\\22")
        return s
}
