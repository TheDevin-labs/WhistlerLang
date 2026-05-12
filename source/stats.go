package main

import (
	"fmt"
	"math"
	"sort"
)

func callStatsFunction(name string, args []WhistlerValue) (WhistlerValue, error) {
	switch name {
	case "mean":
		return statsMean(args)
	case "sum":
		return statsSum(args)
	case "min":
		return statsMin(args)
	case "max":
		return statsMax(args)
	case "std":
		return statsStd(args)
	case "variance":
		return statsVariance(args)
	case "median":
		return statsMedian(args)
	case "len":
		return statsLen(args)
	}
	return WhistlerValue{}, fmt.Errorf("unknown stats function: %s", name)
}

func extractFloats(args []WhistlerValue) ([]float64, error) {
	if len(args) == 1 && args[0].Type == TypeArray {
		var result []float64
		for _, el := range args[0].ArrayVal {
			f, ok := el.ToFloat()
			if !ok {
				return nil, fmt.Errorf("array elements must be numeric")
			}
			result = append(result, f)
		}
		return result, nil
	}
	var result []float64
	for _, a := range args {
		f, ok := a.ToFloat()
		if !ok {
			return nil, fmt.Errorf("arguments must be numeric")
		}
		result = append(result, f)
	}
	return result, nil
}

func statsMean(args []WhistlerValue) (WhistlerValue, error) {
	nums, err := extractFloats(args)
	if err != nil {
		return WhistlerValue{}, err
	}
	if len(nums) == 0 {
		return WhistlerValue{}, fmt.Errorf("mean() requires at least one value")
	}
	sum := 0.0
	for _, n := range nums {
		sum += n
	}
	return FloatValue(sum / float64(len(nums))), nil
}

func statsSum(args []WhistlerValue) (WhistlerValue, error) {
	nums, err := extractFloats(args)
	if err != nil {
		return WhistlerValue{}, err
	}
	sum := 0.0
	for _, n := range nums {
		sum += n
	}
	return FloatValue(sum), nil
}

func statsMin(args []WhistlerValue) (WhistlerValue, error) {
	nums, err := extractFloats(args)
	if err != nil {
		return WhistlerValue{}, err
	}
	if len(nums) == 0 {
		return WhistlerValue{}, fmt.Errorf("min() requires at least one value")
	}
	m := nums[0]
	for _, n := range nums[1:] {
		if n < m {
			m = n
		}
	}
	return FloatValue(m), nil
}

func statsMax(args []WhistlerValue) (WhistlerValue, error) {
	nums, err := extractFloats(args)
	if err != nil {
		return WhistlerValue{}, err
	}
	if len(nums) == 0 {
		return WhistlerValue{}, fmt.Errorf("max() requires at least one value")
	}
	m := nums[0]
	for _, n := range nums[1:] {
		if n > m {
			m = n
		}
	}
	return FloatValue(m), nil
}

func statsVariance(args []WhistlerValue) (WhistlerValue, error) {
	nums, err := extractFloats(args)
	if err != nil {
		return WhistlerValue{}, err
	}
	if len(nums) == 0 {
		return WhistlerValue{}, fmt.Errorf("variance() requires at least one value")
	}
	mean := 0.0
	for _, n := range nums {
		mean += n
	}
	mean /= float64(len(nums))
	variance := 0.0
	for _, n := range nums {
		diff := n - mean
		variance += diff * diff
	}
	variance /= float64(len(nums))
	return FloatValue(variance), nil
}

func statsStd(args []WhistlerValue) (WhistlerValue, error) {
	v, err := statsVariance(args)
	if err != nil {
		return WhistlerValue{}, err
	}
	return FloatValue(math.Sqrt(v.FloatVal)), nil
}

func statsMedian(args []WhistlerValue) (WhistlerValue, error) {
	nums, err := extractFloats(args)
	if err != nil {
		return WhistlerValue{}, err
	}
	if len(nums) == 0 {
		return WhistlerValue{}, fmt.Errorf("median() requires at least one value")
	}
	sorted := make([]float64, len(nums))
	copy(sorted, nums)
	sort.Float64s(sorted)
	n := len(sorted)
	if n%2 == 0 {
		return FloatValue((sorted[n/2-1] + sorted[n/2]) / 2), nil
	}
	return FloatValue(sorted[n/2]), nil
}

func statsLen(args []WhistlerValue) (WhistlerValue, error) {
	if len(args) != 1 {
		return WhistlerValue{}, fmt.Errorf("len() expects 1 argument")
	}
	switch args[0].Type {
	case TypeArray:
		return IntValue(int64(len(args[0].ArrayVal))), nil
	case TypeString:
		return IntValue(int64(len(args[0].StringVal))), nil
	case TypeMatrix:
		return IntValue(int64(len(args[0].MatrixVal))), nil
	}
	return WhistlerValue{}, fmt.Errorf("len() requires array, string, or matrix")
}

