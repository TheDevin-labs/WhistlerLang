package main

import (
	"fmt"
	"math"
)


func callMathFunction(name string, args []WhistlerValue) (WhistlerValue, error) {
	switch name {
	case "sin":
		return mathUnary(args, math.Sin)
	case "cos":
		return mathUnary(args, math.Cos)
	case "tan":
		return mathUnary(args, math.Tan)
	case "asin":
		return mathUnary(args, math.Asin)
	case "acos":
		return mathUnary(args, math.Acos)
	case "atan":
		return mathUnary(args, math.Atan)
	case "sqrt":
		return mathUnary(args, math.Sqrt)
	case "log":
		return mathUnary(args, math.Log)
	case "log2":
		return mathUnary(args, math.Log2)
	case "log10":
		return mathUnary(args, math.Log10)
	case "exp":
		return mathUnary(args, math.Exp)
	case "abs":
		return mathAbs(args)
	case "ceil":
		return mathUnary(args, math.Ceil)
	case "floor":
		return mathUnary(args, math.Floor)
	case "round":
		return mathUnary(args, math.Round)
	case "pow":
		return mathPow(args)
	case "atan2":
		return mathBinary(args, math.Atan2)
	case "hypot":
		return mathBinary(args, math.Hypot)
	case "pi":
		return FloatValue(math.Pi), nil
	case "e":
		return FloatValue(math.E), nil
	case "inf":
		return FloatValue(math.Inf(1)), nil
	case "nan":
		return FloatValue(math.NaN()), nil
	}
	return WhistlerValue{}, fmt.Errorf("unknown math function: %s", name)
}

func mathUnary(args []WhistlerValue, fn func(float64) float64) (WhistlerValue, error) {
	if len(args) != 1 {
		return WhistlerValue{}, fmt.Errorf("expected 1 argument, got %d", len(args))
	}
	f, ok := args[0].ToFloat()
	if !ok {
		return WhistlerValue{}, fmt.Errorf("argument must be numeric")
	}
	return FloatValue(fn(f)), nil
}

func mathBinary(args []WhistlerValue, fn func(float64, float64) float64) (WhistlerValue, error) {
	if len(args) != 2 {
		return WhistlerValue{}, fmt.Errorf("expected 2 arguments, got %d", len(args))
	}
	a, aok := args[0].ToFloat()
	b, bok := args[1].ToFloat()
	if !aok || !bok {
		return WhistlerValue{}, fmt.Errorf("arguments must be numeric")
	}
	return FloatValue(fn(a, b)), nil
}

func mathAbs(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 1 {
		return WhistlerValue{}, fmt.Errorf("abs() expects 1 argument")
	}
	switch args[0].Type {
	case TypeInt:
		v := args[0].IntVal
		if v < 0 {
			v = -v
		}
		return IntValue(v), nil
	case TypeFloat:
		return FloatValue(math.Abs(args[0].FloatVal)), nil
	case TypeComplex:
		c := args[0].ComplexVal
		return FloatValue(math.Sqrt(c.Real*c.Real + c.Imag*c.Imag)), nil
	}
	return WhistlerValue{}, fmt.Errorf("abs() requires numeric argument")
}

func mathPow(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 2 {
		return WhistlerValue{}, fmt.Errorf("pow() expects 2 arguments")
	}
	base, aok := args[0].ToFloat()
	exp, bok := args[1].ToFloat()
	if !aok || !bok {
		return WhistlerValue{}, fmt.Errorf("pow() requires numeric arguments")
	}
	return FloatValue(math.Pow(base, exp)), nil
}

