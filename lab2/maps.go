package main

import "fmt"

func main() {
	marks := []int{40, 50, 60, 70, 80}
	var avg int
	for _, mark := range marks {
		avg += mark
	}
	avg /= len(marks)
	fmt.Println(avg)
}
