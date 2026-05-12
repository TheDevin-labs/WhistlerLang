package main

import (
	"fmt"
	"math"
)

func callLinalgFunction(name string, args []WhistlerValue) (WhistlerValue, error) {
	switch name {
	case "dot":
		return linalgDot(args)
	case "cross":
		return linalgCross(args)
	case "transpose":
		return linalgTranspose(args)
	case "det":
		return linalgDet(args)
	case "inverse":
		return linalgInverse(args)
	case "norm":
		return linalgNorm(args)
	case "rank":
		return linalgRank(args)
	case "zeros":
		return linalgZeros(args)
	case "ones":
		return linalgOnes(args)
	case "identity":
		return linalgIdentity(args)
	}
	return WhistlerValue{}, fmt.Errorf("unknown linalg function: %s", name)
}

func linalgDot(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 2 {
		return WhistlerValue{}, fmt.Errorf("dot() expects 2 arguments")
	}
	a, b := args[0], args[1]

	if a.Type == TypeMatrix && b.Type == TypeMatrix {
		return MatrixMul(a, b)
	}

	if a.Type == TypeArray && b.Type == TypeArray {
		if len(a.ArrayVal) != len(b.ArrayVal) {
			return WhistlerValue{}, fmt.Errorf("dot() arrays must have same length")
		}
		sum := 0.0
		for i := range a.ArrayVal {
			af, _ := a.ArrayVal[i].ToFloat()
			bf, _ := b.ArrayVal[i].ToFloat()
			sum += af * bf
		}
		return FloatValue(sum), nil
	}

	return WhistlerValue{}, fmt.Errorf("dot() expects arrays or matrices")
}

func linalgCross(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 2 {
		return WhistlerValue{}, fmt.Errorf("cross() expects 2 arguments")
	}
	a, b := args[0], args[1]
	if a.Type != TypeArray || b.Type != TypeArray {
		return WhistlerValue{}, fmt.Errorf("cross() expects arrays")
	}
	if len(a.ArrayVal) != 3 || len(b.ArrayVal) != 3 {
		return WhistlerValue{}, fmt.Errorf("cross() expects 3D vectors")
	}
	ax, _ := a.ArrayVal[0].ToFloat()
	ay, _ := a.ArrayVal[1].ToFloat()
	az, _ := a.ArrayVal[2].ToFloat()
	bx, _ := b.ArrayVal[0].ToFloat()
	by, _ := b.ArrayVal[1].ToFloat()
	bz, _ := b.ArrayVal[2].ToFloat()
	return ArrayValue([]WhistlerValue{
		FloatValue(ay*bz - az*by),
		FloatValue(az*bx - ax*bz),
		FloatValue(ax*by - ay*bx),
	}), nil
}

func linalgTranspose(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 1 {
		return WhistlerValue{}, fmt.Errorf("transpose() expects 1 argument")
	}
	return Transpose(args[0])
}

func linalgDet(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 1 || args[0].Type != TypeMatrix {
		return WhistlerValue{}, fmt.Errorf("det() expects a matrix")
	}
	mat := args[0].MatrixVal
	n := len(mat)
	if n == 0 {
		return FloatValue(1), nil
	}
	m := make([][]float64, n)
	for i := range m {
		m[i] = make([]float64, n)
		for j := range m[i] {
			f, ok := mat[i][j].ToFloat()
			if !ok {
				return WhistlerValue{}, fmt.Errorf("matrix must contain numeric values")
			}
			m[i][j] = f
		}
	}
	det := determinant(m, n)
	return FloatValue(det), nil
}

func determinant(m [][]float64, n int) float64 {
	if n == 1 {
		return m[0][0]
	}
	if n == 2 {
		return m[0][0]*m[1][1] - m[0][1]*m[1][0]
	}
	det := 0.0
	for c := 0; c < n; c++ {
		sub := submatrix(m, 0, c, n)
		sign := math.Pow(-1, float64(c))
		det += sign * m[0][c] * determinant(sub, n-1)
	}
	return det
}

func submatrix(m [][]float64, row, col, n int) [][]float64 {
	sub := make([][]float64, n-1)
	si := 0
	for i := 0; i < n; i++ {
		if i == row {
			continue
		}
		sub[si] = make([]float64, n-1)
		sj := 0
		for j := 0; j < n; j++ {
			if j == col {
				continue
			}
			sub[si][sj] = m[i][j]
			sj++
		}
		si++
	}
	return sub
}

func linalgInverse(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 1 || args[0].Type != TypeMatrix {
		return WhistlerValue{}, fmt.Errorf("inverse() expects a matrix")
	}
	mat := args[0].MatrixVal
	n := len(mat)
	m := make([][]float64, n)
	for i := range m {
		m[i] = make([]float64, n)
		for j := range m[i] {
			f, _ := mat[i][j].ToFloat()
			m[i][j] = f
		}
	}
	det := determinant(m, n)
	if math.Abs(det) < 1e-10 {
		return WhistlerValue{}, fmt.Errorf("matrix is singular, cannot invert")
	}
	aug := make([][]float64, n)
	for i := range aug {
		aug[i] = make([]float64, 2*n)
		for j := 0; j < n; j++ {
			aug[i][j] = m[i][j]
		}
		aug[i][n+i] = 1
	}
	for col := 0; col < n; col++ {
		pivot := col
		for row := col + 1; row < n; row++ {
			if math.Abs(aug[row][col]) > math.Abs(aug[pivot][col]) {
				pivot = row
			}
		}
		aug[col], aug[pivot] = aug[pivot], aug[col]
		scale := aug[col][col]
		for j := range aug[col] {
			aug[col][j] /= scale
		}
		for row := 0; row < n; row++ {
			if row != col {
				factor := aug[row][col]
				for j := range aug[row] {
					aug[row][j] -= factor * aug[col][j]
				}
			}
		}
	}
	result := make([][]WhistlerValue, n)
	for i := range result {
		result[i] = make([]WhistlerValue, n)
		for j := 0; j < n; j++ {
			result[i][j] = FloatValue(aug[i][n+j])
		}
	}
	return MatrixValue(result), nil
}

func linalgNorm(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 1 || args[0].Type != TypeArray {
		return WhistlerValue{}, fmt.Errorf("norm() expects an array")
	}
	sum := 0.0
	for _, el := range args[0].ArrayVal {
		f, _ := el.ToFloat()
		sum += f * f
	}
	return FloatValue(math.Sqrt(sum)), nil
}

func linalgRank(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 1 || args[0].Type != TypeMatrix {
		return WhistlerValue{}, fmt.Errorf("rank() expects a matrix")
	}
	mat := args[0].MatrixVal
	rows := len(mat)
	if rows == 0 {
		return IntValue(0), nil
	}
	cols := len(mat[0])
	m := make([][]float64, rows)
	for i := range m {
		m[i] = make([]float64, cols)
		for j := range m[i] {
			f, _ := mat[i][j].ToFloat()
			m[i][j] = f
		}
	}
	rank := 0
	rowUsed := make([]bool, rows)
	for col := 0; col < cols; col++ {
		pivot := -1
		for row := 0; row < rows; row++ {
			if !rowUsed[row] && math.Abs(m[row][col]) > 1e-10 {
				pivot = row
				break
			}
		}
		if pivot == -1 {
			continue
		}
		rowUsed[pivot] = true
		rank++
		for row := 0; row < rows; row++ {
			if row != pivot && math.Abs(m[row][col]) > 1e-10 {
				factor := m[row][col] / m[pivot][col]
				for j := 0; j < cols; j++ {
					m[row][j] -= factor * m[pivot][j]
				}
			}
		}
	}
	return IntValue(int64(rank)), nil
}

func linalgZeros(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 2 {
		return WhistlerValue{}, fmt.Errorf("zeros(rows, cols) expects 2 arguments")
	}
	rows := int(args[0].IntVal)
	cols := int(args[1].IntVal)
	mat := make([][]WhistlerValue, rows)
	for i := range mat {
		mat[i] = make([]WhistlerValue, cols)
		for j := range mat[i] {
			mat[i][j] = FloatValue(0)
		}
	}
	return MatrixValue(mat), nil
}

func linalgOnes(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 2 {
		return WhistlerValue{}, fmt.Errorf("ones(rows, cols) expects 2 arguments")
	}
	rows := int(args[0].IntVal)
	cols := int(args[1].IntVal)
	mat := make([][]WhistlerValue, rows)
	for i := range mat {
		mat[i] = make([]WhistlerValue, cols)
		for j := range mat[i] {
			mat[i][j] = FloatValue(1)
		}
	}
	return MatrixValue(mat), nil
}

func linalgIdentity(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 1 {
		return WhistlerValue{}, fmt.Errorf("identity(n) expects 1 argument")
	}
	n := int(args[0].IntVal)
	mat := make([][]WhistlerValue, n)
	for i := range mat {
		mat[i] = make([]WhistlerValue, n)
		for j := range mat[i] {
			if i == j {
				mat[i][j] = FloatValue(1)
			} else {
				mat[i][j] = FloatValue(0)
			}
		}
	}
	return MatrixValue(mat), nil
}

