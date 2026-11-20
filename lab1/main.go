package main

import "fmt"

func main() {
	fmt.Println("Go Environment Verified!")

	a := 10
	b := 4
	x := 12.5
	y := 2.5

	fmt.Println("\nInteger Operations:")
	fmt.Println("Addition:", a+b)
	fmt.Println("Subtraction:", a-b)
	fmt.Println("Multiplication:", a*b)

	fmt.Println("\nFloat Operations:")
	fmt.Println("Addition:", x+y)
	fmt.Println("Subtraction:", x-y)
	fmt.Println("Multiplication:", x*y)
}
