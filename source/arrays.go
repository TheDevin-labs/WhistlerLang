package main

import "fmt"


func NewArray(elements []WhistlerValue) WhistlerValue {
	return ArrayValue(elements)
}

func NewMatrix(rows [][]WhistlerValue) WhistlerValue {
	return MatrixValue(rows)
}

func ArrayGet(arr WhistlerValue, idx int) (WhistlerValue, error) {
	if arr.Type != TypeArray {
		return WhistlerValue{}, fmt.Errorf("not an array")
	}
	if idx < 0 || idx >= len(arr.ArrayVal) {
		return WhistlerValue{}, fmt.Errorf("index %d out of bounds (length %d)", idx, len(arr.ArrayVal))
	}
	return arr.ArrayVal[idx], nil
}

func MatrixGet(mat WhistlerValue, row, col int) (WhistlerValue, error) {
	if mat.Type != TypeMatrix {
		return WhistlerValue{}, fmt.Errorf("not a matrix")
	}
	if row < 0 || row >= len(mat.MatrixVal) {
		return WhistlerValue{}, fmt.Errorf("row index %d out of bounds", row)
	}
	if col < 0 || col >= len(mat.MatrixVal[row]) {
		return WhistlerValue{}, fmt.Errorf("col index %d out of bounds", col)
	}
	return mat.MatrixVal[row][col], nil
}

func ArrayLen(arr WhistlerValue) (WhistlerValue, error) {
	if arr.Type != TypeArray {
		return WhistlerValue{}, fmt.Errorf("not an array")
	}
	return IntValue(int64(len(arr.ArrayVal))), nil
}

func MatrixShape(mat WhistlerValue) (WhistlerValue, error) {
	if mat.Type != TypeMatrix {
		return WhistlerValue{}, fmt.Errorf("not a matrix")
	}
	rows := int64(len(mat.MatrixVal))
	cols := int64(0)
	if rows > 0 {
		cols = int64(len(mat.MatrixVal[0]))
	}
	return ArrayValue([]WhistlerValue{IntValue(rows), IntValue(cols)}), nil
}

func MatrixAdd(a, b WhistlerValue) (WhistlerValue, error) {
	if err := checkSameShape(a, b); err != nil {
		return WhistlerValue{}, err
	}
	result := make([][]WhistlerValue, len(a.MatrixVal))
	for i, row := range a.MatrixVal {
		result[i] = make([]WhistlerValue, len(row))
		for j := range row {
			val, err := addValues(a.MatrixVal[i][j], b.MatrixVal[i][j])
			if err != nil {
				return WhistlerValue{}, err
			}
			result[i][j] = val
		}
	}
	return MatrixValue(result), nil
}

func MatrixMul(a, b WhistlerValue) (WhistlerValue, error) {
	if a.Type != TypeMatrix || b.Type != TypeMatrix {
		return WhistlerValue{}, fmt.Errorf("both operands must be matrices")
	}
	aRows := len(a.MatrixVal)
	if aRows == 0 {
		return WhistlerValue{}, fmt.Errorf("empty matrix")
	}
	aCols := len(a.MatrixVal[0])
	bRows := len(b.MatrixVal)
	if aCols != bRows {
		return WhistlerValue{}, fmt.Errorf("matrix dimension mismatch: %dx%d * %dx%d", aRows, aCols, bRows, len(b.MatrixVal[0]))
	}
	bCols := len(b.MatrixVal[0])

	result := make([][]WhistlerValue, aRows)
	for i := range result {
		result[i] = make([]WhistlerValue, bCols)
		for j := 0; j < bCols; j++ {
			sum := 0.0
			for k := 0; k < aCols; k++ {
				af, _ := a.MatrixVal[i][k].ToFloat()
				bf, _ := b.MatrixVal[k][j].ToFloat()
				sum += af * bf
			}
			result[i][j] = FloatValue(sum)
		}
	}
	return MatrixValue(result), nil
}

func Transpose(mat WhistlerValue) (WhistlerValue, error) {
	if mat.Type != TypeMatrix {
		return WhistlerValue{}, fmt.Errorf("not a matrix")
	}
	if len(mat.MatrixVal) == 0 {
		return mat, nil
	}
	rows := len(mat.MatrixVal)
	cols := len(mat.MatrixVal[0])
	result := make([][]WhistlerValue, cols)
	for i := range result {
		result[i] = make([]WhistlerValue, rows)
		for j := 0; j < rows; j++ {
			result[i][j] = mat.MatrixVal[j][i]
		}
	}
	return MatrixValue(result), nil
}

func ArrayAppend(arr, val WhistlerValue) (WhistlerValue, error) {
	if arr.Type != TypeArray {
		return WhistlerValue{}, fmt.Errorf("not an array")
	}
	newArr := make([]WhistlerValue, len(arr.ArrayVal)+1)
	copy(newArr, arr.ArrayVal)
	newArr[len(arr.ArrayVal)] = val
	return ArrayValue(newArr), nil
}

func ArraySlice(arr WhistlerValue, start, end int) (WhistlerValue, error) {
	if arr.Type != TypeArray {
		return WhistlerValue{}, fmt.Errorf("not an array")
	}
	if start < 0 || end > len(arr.ArrayVal) || start > end {
		return WhistlerValue{}, fmt.Errorf("slice out of bounds")
	}
	return ArrayValue(arr.ArrayVal[start:end]), nil
}

func checkSameShape(a, b WhistlerValue) error {
	if a.Type != TypeMatrix || b.Type != TypeMatrix {
		return fmt.Errorf("both operands must be matrices")
	}
	if len(a.MatrixVal) != len(b.MatrixVal) {
		return fmt.Errorf("matrix row count mismatch")
	}
	for i := range a.MatrixVal {
		if len(a.MatrixVal[i]) != len(b.MatrixVal[i]) {
			return fmt.Errorf("matrix column count mismatch at row %d", i)
		}
	}
	return nil
}

func addValues(a, b WhistlerValue) (WhistlerValue, error) {
	af, aok := a.ToFloat()
	bf, bok := b.ToFloat()
	if aok && bok {
		return FloatValue(af + bf), nil
	}
	return WhistlerValue{}, fmt.Errorf("cannot add non-numeric values")
}

