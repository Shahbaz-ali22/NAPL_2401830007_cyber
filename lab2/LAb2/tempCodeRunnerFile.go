/ Slice for 5 students
	students := make([]string, 5)

	// Take input for exactly 5 students
	for i := 0; i < 5; i++ {
		fmt.Printf("Enter student %d name: ", i+1)
		fmt.Scan(&stud