package main

import "fmt"

type ValueType int

const (
	TypeInt     ValueType = iota
	TypeFloat
	TypeComplex
	TypeBool
	TypeString
	TypeArray
	TypeMatrix
	TypeVoid
)

func (t ValueType) String() string {
	switch t {
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeComplex:
		return "complex"
	case TypeBool:
		return "bool"
	case TypeString:
		return "string"
	case TypeArray:
		return "array"
	case TypeMatrix:
		return "matrix"
	case TypeVoid:
		return "void"
	default:
		return "unknown"
	}
}

type WhistlerValue struct {
	Type    ValueType
	IntVal  int64
	FloatVal float64
	ComplexVal ComplexNum
	BoolVal bool
	StringVal string
	ArrayVal  []WhistlerValue
	MatrixVal [][]WhistlerValue
}

type ComplexNum struct {
	Real float64
	Imag float64
}

func (c ComplexNum) String() string {
	if c.Imag >= 0 {
		return fmt.Sprintf("%g+%gi", c.Real, c.Imag)
	}
	return fmt.Sprintf("%g%gi", c.Real, c.Imag)
}

func IntValue(v int64) WhistlerValue {
	return WhistlerValue{Type: TypeInt, IntVal: v}
}

func FloatValue(v float64) WhistlerValue {
	return WhistlerValue{Type: TypeFloat, FloatVal: v}
}

func ComplexValue(r, i float64) WhistlerValue {
	return WhistlerValue{Type: TypeComplex, ComplexVal: ComplexNum{Real: r, Imag: i}}
}

func BoolValue(v bool) WhistlerValue {
	return WhistlerValue{Type: TypeBool, BoolVal: v}
}

func StringValue(v string) WhistlerValue {
	return WhistlerValue{Type: TypeString, StringVal: v}
}

func ArrayValue(v []WhistlerValue) WhistlerValue {
	return WhistlerValue{Type: TypeArray, ArrayVal: v}
}

func MatrixValue(v [][]WhistlerValue) WhistlerValue {
	return WhistlerValue{Type: TypeMatrix, MatrixVal: v}
}

func (v WhistlerValue) String() string {
	switch v.Type {
	case TypeInt:
		return fmt.Sprintf("%d", v.IntVal)
	case TypeFloat:
		return fmt.Sprintf("%g", v.FloatVal)
	case TypeComplex:
		return v.ComplexVal.String()
	case TypeBool:
		if v.BoolVal {
			return "true"
		}
		return "false"
	case TypeString:
		return v.StringVal
	case TypeArray:
		parts := make([]string, len(v.ArrayVal))
		for i, el := range v.ArrayVal {
			parts[i] = el.String()
		}
		return "[" + joinStrings(parts, ", ") + "]"
	case TypeMatrix:
		rows := make([]string, len(v.MatrixVal))
		for i, row := range v.MatrixVal {
			cols := make([]string, len(row))
			for j, el := range row {
				cols[j] = el.String()
			}
			rows[i] = "[" + joinStrings(cols, ", ") + "]"
		}
		return "[" + joinStrings(rows, ", ") + "]"
	default:
		return "void"
	}
}

func (v WhistlerValue) ToFloat() (float64, bool) {
	switch v.Type {
	case TypeInt:
		return float64(v.IntVal), true
	case TypeFloat:
		return v.FloatVal, true
	default:
		return 0, false
	}
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

