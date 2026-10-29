package main

import "fmt"

func PrintName(name string) {
	fmt.Println(name)
}

func main() {
	go PrintName("John")
	PrintName("Doe")
}
