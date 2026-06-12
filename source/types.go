package main

import "fmt"

type ValueType int

const (
	TypeInt     ValueType = iota
	TypeFloat
	TypeComplex
	TypeBool
	TypeString
	TypeByte
	TypeBytes
	TypeArray
	TypeMatrix
	TypeVoid
)

func (t ValueType) String() string {
	switch t {
	case TypeInt:     return "int"
	case TypeFloat:   return "float"
	case TypeComplex: return "complex"
	case TypeBool:    return "bool"
	case TypeString:  return "string"
	case TypeByte:    return "byte"
	case TypeBytes:   return "bytes"
	case TypeArray:   return "array"
	case TypeMatrix:  return "matrix"
	default:          return "void"
	}
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

type WhistlerValue struct {
	Type       ValueType
	IntVal     int64
	FloatVal   float64
	ComplexVal ComplexNum
	BoolVal    bool
	StringVal  string
	ByteVal    uint8
	BytesVal   []uint8
	ArrayVal   []WhistlerValue
	MatrixVal  [][]WhistlerValue
}

func IntValue(v int64) WhistlerValue     { return WhistlerValue{Type: TypeInt, IntVal: v} }
func FloatValue(v float64) WhistlerValue { return WhistlerValue{Type: TypeFloat, FloatVal: v} }
func BoolValue(v bool) WhistlerValue     { return WhistlerValue{Type: TypeBool, BoolVal: v} }
func StringValue(v string) WhistlerValue { return WhistlerValue{Type: TypeString, StringVal: v} }
func ByteValue(v uint8) WhistlerValue    { return WhistlerValue{Type: TypeByte, ByteVal: v} }
func BytesValue(v []uint8) WhistlerValue { return WhistlerValue{Type: TypeBytes, BytesVal: v} }
func ArrayValue(v []WhistlerValue) WhistlerValue   { return WhistlerValue{Type: TypeArray, ArrayVal: v} }
func MatrixValue(v [][]WhistlerValue) WhistlerValue { return WhistlerValue{Type: TypeMatrix, MatrixVal: v} }
func ComplexValue(r, i float64) WhistlerValue {
	return WhistlerValue{Type: TypeComplex, ComplexVal: ComplexNum{Real: r, Imag: i}}
}

func (v WhistlerValue) ToFloat() (float64, bool) {
	switch v.Type {
	case TypeInt:   return float64(v.IntVal), true
	case TypeFloat: return v.FloatVal, true
	case TypeByte:  return float64(v.ByteVal), true
	}
	return 0, false
}

func (v WhistlerValue) String() string {
	switch v.Type {
	case TypeInt:     return fmt.Sprintf("%d", v.IntVal)
	case TypeFloat:   return fmt.Sprintf("%g", v.FloatVal)
	case TypeComplex: return v.ComplexVal.String()
	case TypeBool:
		if v.BoolVal { return "true" }
		return "false"
	case TypeString:  return v.StringVal
	case TypeByte:    return fmt.Sprintf("%d", v.ByteVal)
	case TypeBytes:
		s := "["
		for i, b := range v.BytesVal {
			if i > 0 { s += ", " }
			s += fmt.Sprintf("%d", b)
		}
		return s + "]"
	case TypeArray:
		s := "["
		for i, el := range v.ArrayVal {
			if i > 0 { s += ", " }
			s += el.String()
		}
		return s + "]"
	case TypeMatrix:
		s := "["
		for i, row := range v.MatrixVal {
			if i > 0 { s += ", " }
			s += "["
			for j, el := range row {
				if j > 0 { s += ", " }
				s += el.String()
			}
			s += "]"
		}
		return s + "]"
	}
	return "void"
}
