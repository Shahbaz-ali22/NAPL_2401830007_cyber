package main

import "fmt"

func main() {
	var a, b int
	var x, y float64

	for {
		fmt.Print("Enter two integers: ")
		fmt.Scan(&a, &b)
		if a > 0 && b > 0 {
			break
		}
		fmt.Println("Invalid input! Please enter positive numbers.")
	}

	fmt.Println("Integer Add:", a+b)
	fmt.Println("Integer Sub:", a-b)
	fmt.Println("Integer Mul:", a*b)

	for {
		fmt.Print("\nEnter two floats: ")
		fmt.Scan(&x, &y)
		if x > 0 && y > 0 {
			break
		}
		fmt.Println("Invalid input! Please enter positive numbers.")
	}

	fmt.Println("Float Add:", x+y)
	fmt.Println("Float Sub:", x-y)
	fmt.Println("Float Mul:", x*y)
}
