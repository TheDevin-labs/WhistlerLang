package main

import "fmt"

func ArrayGet(arr WhistlerValue, idx int) (WhistlerValue, error) {
	if arr.Type != TypeArray {
		return WhistlerValue{}, fmt.Errorf("not an array")
	}
	if idx < 0 || idx >= len(arr.ArrayVal) {
		return WhistlerValue{}, fmt.Errorf("index %d out of bounds (len %d)", idx, len(arr.ArrayVal))
	}
	return arr.ArrayVal[idx], nil
}

func MatrixGet(mat WhistlerValue, row, col int) (WhistlerValue, error) {
	if mat.Type != TypeMatrix {
		return WhistlerValue{}, fmt.Errorf("not a matrix")
	}
	if row < 0 || row >= len(mat.MatrixVal) {
		return WhistlerValue{}, fmt.Errorf("row %d out of bounds", row)
	}
	if col < 0 || col >= len(mat.MatrixVal[row]) {
		return WhistlerValue{}, fmt.Errorf("col %d out of bounds", col)
	}
	return mat.MatrixVal[row][col], nil
}

func MatrixAdd(a, b WhistlerValue) (WhistlerValue, error) {
	if a.Type != TypeMatrix || b.Type != TypeMatrix {
		return WhistlerValue{}, fmt.Errorf("both must be matrices")
	}
	if len(a.MatrixVal) != len(b.MatrixVal) {
		return WhistlerValue{}, fmt.Errorf("matrix size mismatch")
	}
	result := make([][]WhistlerValue, len(a.MatrixVal))
	for i, row := range a.MatrixVal {
		result[i] = make([]WhistlerValue, len(row))
		for j := range row {
			af, _ := a.MatrixVal[i][j].ToFloat()
			bf, _ := b.MatrixVal[i][j].ToFloat()
			result[i][j] = FloatValue(af + bf)
		}
	}
	return MatrixValue(result), nil
}

func MatrixMul(a, b WhistlerValue) (WhistlerValue, error) {
	if a.Type != TypeMatrix || b.Type != TypeMatrix {
		return WhistlerValue{}, fmt.Errorf("both must be matrices")
	}
	aRows, aCols := len(a.MatrixVal), len(a.MatrixVal[0])
	bRows, bCols := len(b.MatrixVal), len(b.MatrixVal[0])
	if aCols != bRows {
		return WhistlerValue{}, fmt.Errorf("dimension mismatch: %dx%d * %dx%d", aRows, aCols, bRows, bCols)
	}
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
	rows, cols := len(mat.MatrixVal), len(mat.MatrixVal[0])
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
