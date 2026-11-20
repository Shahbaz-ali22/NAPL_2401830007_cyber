package main

import (
	"fmt"
	"math"
)

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	return a / b, nil
}

func Sqrt(n float64) (float64, error) {
	if n < 0 {
		return 0, fmt.Errorf("cannot calculate square root of negative number: %v", n)
	}
	return math.Sqrt(n), nil
}

func CheckEmpty(s string) (string, error) {
	if s == "" {
		return "", fmt.Errorf("input string cannot be empty")
	}
	return s, nil
}
