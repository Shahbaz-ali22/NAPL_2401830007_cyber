package lab2
package main

import "fmt"

func main() {

	s := []int{10, 20}

	s = append(s, 30)

	fmt.Println(s)

	nums2 := make([]int, 5, 10)

	fmt.Println(nums2)
}