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

func main() {
	res1, err1 := Divide(10, 0)
	if err1 != nil {
		fmt.Println("Error:", err1)
	} else {
		fmt.Println("Result:", res1)
	}

	res2, err2 := Sqrt(-9)
	if err2 != nil {
		fmt.Println("Error:", err2)
	} else {
		fmt.Println("Result:", res2)
	}

	res3, err3 := CheckEmpty("")
	if err3 != nil {
		fmt.Println("Error:", err3)
	} else {
		fmt.Println("Result:", res3)
	}
}
