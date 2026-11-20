package main

import "testing"

func TestDivide(t *testing.T) {
	val, err := Divide(10, 2)
	if err != nil || val != 5 {
		t.Errorf("Expected 5, got %v", val)
	}

	_, err = Divide(10, 0)
	if err == nil {
		t.Errorf("Expected division by zero error")
	}
}

func TestSqrt(t *testing.T) {
	val, err := Sqrt(16)
	if err != nil || val != 4 {
		t.Errorf("Expected 4, got %v", val)
	}

	_, err = Sqrt(-4)
	if err == nil {
		t.Errorf("Expected negative square root error")
	}
}

func TestCheckEmpty(t *testing.T) {
	val, err := CheckEmpty("hello")
	if err != nil || val != "hello" {
		t.Errorf("Expected hello, got %v", val)
	}

	_, err = CheckEmpty("")
	if err == nil {
		t.Errorf("Expected empty input error")
	}
}

