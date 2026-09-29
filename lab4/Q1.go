package main

import (
	"fmt"
)

func square(ch chan int, num int) {
	fmt.Println("Square goroutine started")
	result := num * num
	ch <- result
	fmt.Println("Square goroutine finished")

}

func cube(ch chan int, num int) {
	fmt.Println("Cube goroutine started")
	result := num * num * num
	ch <- result
	fmt.Println("Cube goroutine finished")
}

func fibonacci(ch chan int, num int) {
	fmt.Println("Fibonacci goroutine started")

	a, b := 0, 1
	for i := 0; i < num; i++ {
		a, b = b, a+b
	}

	ch <- a
	fmt.Println("Fibonacci goroutine finished")
}

func main() {
	ch := make(chan int)

	go square(ch, 5)
	go cube(ch, 5)
	go fibonacci(ch, 5)

	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
