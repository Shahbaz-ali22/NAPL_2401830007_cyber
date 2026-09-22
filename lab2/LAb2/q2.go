package main

import "fmt"

func main() {

	// Slice for 5 students
	students := make([]string, 5)

	// Take input for exactly 5 students
	for i := 0; i < 5; i++ {
		fmt.Printf("Enter student %d name: ", i+1)
		fmt.Scan(&students[i])
	}

	fmt.Println("\nInitial slice:", students)

	// Add
	var name string

	fmt.Print("\nEnter student name to add: ")
	fmt.Scan(&name)

	students = append(students, name)

	fmt.Println("After adding:", students)

	// Remove by index
	var index int

	fmt.Print("\nEnter index to remove: ")
	fmt.Scan(&index)

	if index >= 0 && index < len(students) {
		students = append(students[:index], students[index+1:]...)
		fmt.Println("After removing:", students)
	} else {
		fmt.Println("Invalid index!")
	}

	// Update
	fmt.Print("\nEnter index to update: ")
	fmt.Scan(&index)

	if index >= 0 && index < len(students) {
		fmt.Print("Enter new student name: ")
		fmt.Scan(&name)

		students[index] = name

		fmt.Println("After updating:", students)
	} else {
		fmt.Println("Invalid index!")
	}
}
