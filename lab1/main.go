package main

import "fmt"
import "example.com/app/mathutil"

func main() {
	sum := mathutil.Add(2, 3)
	fmt.Println("Sum=%d", sum)
}
