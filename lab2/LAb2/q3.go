package main

import "fmt"

func main() {

	marks := map[string]int{
		"Math":    85,
		"English": 78,
		"Science": 90,
	}

	fmt.Println("Initial map:", marks)

	// Insert
	var subject string
	var mark int

	fmt.Print("Enter subject to insert: ")
	fmt.Scan(&subject)

	fmt.Print("Enter marks: ")
	fmt.Scan(&mark)

	marks[subject] = mark

	fmt.Println("After insertion:", marks)

	// Delete
	fmt.Print("Enter subject to delete: ")
	fmt.Scan(&subject)

	delete(marks, subject)

	fmt.Println("After deletion:", marks)

	// Lookup
	fmt.Print("Enter subject to lookup: ")
	fmt.Scan(&subject)

	value, exists := marks[subject]

	if exists {
		fmt.Println("Marks:", value)
	} else {
		fmt.Println("Subject not found")
	}

	fmt.Println("Final map:", marks)
}
