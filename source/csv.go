package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func csvAutoDetect(cell string) WhistlerValue {
	cell = strings.TrimSpace(cell)
	if v, err := strconv.ParseInt(cell, 10, 64); err == nil {
		return IntValue(v)
	}
	if v, err := strconv.ParseFloat(cell, 64); err == nil {
		return FloatValue(v)
	}
	return StringValue(cell)
}

func csvParseLine(line string) []WhistlerValue {
	parts := strings.Split(line, ",")
	row := make([]WhistlerValue, len(parts))
	for i, p := range parts {
		row[i] = csvAutoDetect(p)
	}
	return row
}

func CsvOpen(path string) (WhistlerValue, error) {
	f, err := os.Open(path)
	if err != nil {
		return WhistlerValue{}, fmt.Errorf("csv.open: cannot open file %q: %w", path, err)
	}
	defer f.Close()

	var rows [][]WhistlerValue
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		rows = append(rows, csvParseLine(line))
	}
	if err := scanner.Err(); err != nil {
		return WhistlerValue{}, fmt.Errorf("csv.open: read error: %w", err)
	}

	result := make([][]WhistlerValue, len(rows))
	for i, row := range rows {
		vals := make([]WhistlerValue, len(row))
		copy(vals, row)
		result[i] = vals
	}
	return MatrixValue(result), nil
}

func CsvLineReader(path string) ([]WhistlerValue, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("csv.line: cannot open file %q: %w", path, err)
	}
	defer f.Close()

	var lines []WhistlerValue
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		row := csvParseLine(line)
		lines = append(lines, ArrayValue(row))
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("csv.line: read error: %w", err)
	}
	return lines, nil
}
