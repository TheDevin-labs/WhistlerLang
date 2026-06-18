package main

import (
	"fmt"
	"math"
)

func callLinalg(name string, args []WhistlerValue) (WhistlerValue, error) {
	switch name {
	case "dot":       return linDot(args)
	case "cross":     return linCross(args)
	case "transpose": return linTranspose(args)
	case "det":       return linDet(args)
	case "inverse":   return linInverse(args)
	case "norm":      return linNorm(args)
	case "rank":      return linRank(args)
	case "zeros":     return linZeros(args)
	case "ones":      return linOnes(args)
	case "identity":  return linIdentity(args)
	}
	return WhistlerValue{}, fmt.Errorf("unknown linalg function: %s", name)
}

func linDot(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 2 { return WhistlerValue{}, fmt.Errorf("dot() expects 2 arguments") }
	a, b := args[0], args[1]
	if a.Type == TypeMatrix && b.Type == TypeMatrix { return MatrixMul(a, b) }
	if a.Type == TypeArray && b.Type == TypeArray {
		if len(a.ArrayVal) != len(b.ArrayVal) { return WhistlerValue{}, fmt.Errorf("dot() arrays must match in length") }
		s := 0.0
		for i := range a.ArrayVal { af, _ := a.ArrayVal[i].ToFloat(); bf, _ := b.ArrayVal[i].ToFloat(); s += af * bf }
		return FloatValue(s), nil
	}
	return WhistlerValue{}, fmt.Errorf("dot() expects arrays or matrices")
}

func linCross(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 2 { return WhistlerValue{}, fmt.Errorf("cross() expects 2 arguments") }
	a, b := args[0], args[1]
	if a.Type != TypeArray || b.Type != TypeArray || len(a.ArrayVal) != 3 || len(b.ArrayVal) != 3 {
		return WhistlerValue{}, fmt.Errorf("cross() expects two 3D vectors")
	}
	ax, _ := a.ArrayVal[0].ToFloat(); ay, _ := a.ArrayVal[1].ToFloat(); az, _ := a.ArrayVal[2].ToFloat()
	bx, _ := b.ArrayVal[0].ToFloat(); by, _ := b.ArrayVal[1].ToFloat(); bz, _ := b.ArrayVal[2].ToFloat()
	return ArrayValue([]WhistlerValue{
		FloatValue(ay*bz - az*by),
		FloatValue(az*bx - ax*bz),
		FloatValue(ax*by - ay*bx),
	}), nil
}

func linTranspose(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 1 { return WhistlerValue{}, fmt.Errorf("transpose() expects 1 argument") }
	return Transpose(args[0])
}

func matToFloat(mat WhistlerValue) ([][]float64, int, error) {
	n := len(mat.MatrixVal)
	m := make([][]float64, n)
	for i := range m {
		m[i] = make([]float64, len(mat.MatrixVal[i]))
		for j := range m[i] {
			f, ok := mat.MatrixVal[i][j].ToFloat()
			if !ok { return nil, 0, fmt.Errorf("matrix must contain numeric values") }
			m[i][j] = f
		}
	}
	return m, n, nil
}

func det(m [][]float64, n int) float64 {
	if n == 1 { return m[0][0] }
	if n == 2 { return m[0][0]*m[1][1] - m[0][1]*m[1][0] }
	d := 0.0
	for c := 0; c < n; c++ {
		sub := make([][]float64, n-1)
		for i := 1; i < n; i++ {
			sub[i-1] = make([]float64, 0)
			for j := 0; j < n; j++ { if j != c { sub[i-1] = append(sub[i-1], m[i][j]) } }
		}
		sign := 1.0; if c%2 != 0 { sign = -1.0 }
		d += sign * m[0][c] * det(sub, n-1)
	}
	return d
}

func linDet(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 1 || args[0].Type != TypeMatrix { return WhistlerValue{}, fmt.Errorf("det() expects a matrix") }
	m, n, err := matToFloat(args[0]); if err != nil { return WhistlerValue{}, err }
	return FloatValue(det(m, n)), nil
}

func linInverse(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 1 || args[0].Type != TypeMatrix { return WhistlerValue{}, fmt.Errorf("inverse() expects a matrix") }
	m, n, err := matToFloat(args[0]); if err != nil { return WhistlerValue{}, err }
	if math.Abs(det(m, n)) < 1e-10 { return WhistlerValue{}, fmt.Errorf("matrix is singular") }
	aug := make([][]float64, n)
	for i := range aug {
		aug[i] = make([]float64, 2*n)
		copy(aug[i], m[i])
		aug[i][n+i] = 1
	}
	for col := 0; col < n; col++ {
		pivot := col
		for row := col + 1; row < n; row++ { if math.Abs(aug[row][col]) > math.Abs(aug[pivot][col]) { pivot = row } }
		aug[col], aug[pivot] = aug[pivot], aug[col]
		sc := aug[col][col]
		for j := range aug[col] { aug[col][j] /= sc }
		for row := 0; row < n; row++ {
			if row != col { f := aug[row][col]; for j := range aug[row] { aug[row][j] -= f * aug[col][j] } }
		}
	}
	result := make([][]WhistlerValue, n)
	for i := range result {
		result[i] = make([]WhistlerValue, n)
		for j := 0; j < n; j++ { result[i][j] = FloatValue(aug[i][n+j]) }
	}
	return MatrixValue(result), nil
}

func linNorm(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 1 || args[0].Type != TypeArray { return WhistlerValue{}, fmt.Errorf("norm() expects an array") }
	s := 0.0; for _, el := range args[0].ArrayVal { f, _ := el.ToFloat(); s += f * f }
	return FloatValue(math.Sqrt(s)), nil
}

func linRank(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 1 || args[0].Type != TypeMatrix { return WhistlerValue{}, fmt.Errorf("rank() expects a matrix") }
	m, n, err := matToFloat(args[0]); if err != nil { return WhistlerValue{}, err }
	cols := len(m[0]); rank := 0; used := make([]bool, n)
	for col := 0; col < cols; col++ {
		pivot := -1
		for row := 0; row < n; row++ { if !used[row] && math.Abs(m[row][col]) > 1e-10 { pivot = row; break } }
		if pivot == -1 { continue }
		used[pivot] = true; rank++
		for row := 0; row < n; row++ {
			if row != pivot && math.Abs(m[row][col]) > 1e-10 {
				f := m[row][col] / m[pivot][col]
				for j := 0; j < cols; j++ { m[row][j] -= f * m[pivot][j] }
			}
		}
	}
	return IntValue(int64(rank)), nil
}

func linZeros(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 2 { return WhistlerValue{}, fmt.Errorf("zeros(rows, cols) expects 2 arguments") }
	rows, cols := int(args[0].IntVal), int(args[1].IntVal)
	mat := make([][]WhistlerValue, rows)
	for i := range mat { mat[i] = make([]WhistlerValue, cols); for j := range mat[i] { mat[i][j] = FloatValue(0) } }
	return MatrixValue(mat), nil
}

func linOnes(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 2 { return WhistlerValue{}, fmt.Errorf("ones(rows, cols) expects 2 arguments") }
	rows, cols := int(args[0].IntVal), int(args[1].IntVal)
	mat := make([][]WhistlerValue, rows)
	for i := range mat { mat[i] = make([]WhistlerValue, cols); for j := range mat[i] { mat[i][j] = FloatValue(1) } }
	return MatrixValue(mat), nil
}

func linIdentity(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 1 { return WhistlerValue{}, fmt.Errorf("identity(n) expects 1 argument") }
	n := int(args[0].IntVal)
	mat := make([][]WhistlerValue, n)
	for i := range mat { mat[i] = make([]WhistlerValue, n); for j := range mat[i] { if i == j { mat[i][j] = FloatValue(1) } else { mat[i][j] = FloatValue(0) } } }
	return MatrixValue(mat), nil
}
