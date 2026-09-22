package main

import "fmt"

type Person struct {
	Name   string
	Age    int
	Job    string
	Salary float64
}

// Method to take input
func (p *Person) readData() {
	fmt.Print("Enter Name: ")
	fmt.Scan(&p.Name)

	fmt.Print("Enter Age: ")
	fmt.Scan(&p.Age)

	fmt.Print("Enter Job: ")
	fmt.Scan(&p.Job)

	fmt.Print("Enter Salary: ")
	fmt.Scan(&p.Salary)
}

// Method to print data
func (p Person) printData() {
	fmt.Println("Name:", p.Name)
	fmt.Println("Age:", p.Age)
	fmt.Println("Job:", p.Job)
	fmt.Println("Salary:", p.Salary)
}

func main() {

	var person1 Person
	var person2 Person

	fmt.Println("Enter Person 1 details:")
	person1.readData()

	fmt.Println("\nEnter Person 2 details:")
	person2.readData()

	fmt.Println("\nPerson 1:")
	person1.printData()

	fmt.Println("\nPerson 2:")
	person2.printData()
}