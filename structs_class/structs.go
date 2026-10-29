package main

import "fmt"

type Student struct {
	name   string
	rollNo int
	marks  int
}

func main() {

	students := []Student{
		{"Shahbaz", 1, 85},
		{"Rahul", 2, 78},
	}
	for _, student := range students {
		fmt.Println(student.name, student.rollNo, student.marks)
	}
}
