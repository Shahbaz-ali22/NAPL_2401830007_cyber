package main

import "fmt"

type Student struct {
	Name  string
	Marks float32
}

func change(a *int) {
	*a = 20
}

func main() {
	a := 10

	fmt.Println(a)
	fmt.Println(&a)
	fmt.Println(*(&a))

	fmt.Println("Before change:", a)
	change(&a)
	fmt.Println("After change:", a)

	s := new(Student)

	s.Name = "XYZ"
	s.Marks = 85.5

	fmt.Println(s.Name)
	fmt.Println(s.Marks)

	s.Marks = 90
	fmt.Println(s.Marks)
}
