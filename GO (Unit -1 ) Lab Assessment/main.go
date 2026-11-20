package main

import "fmt"

type Item struct {
	ID       int
	Name     string
	Quantity int
	Price    float64
}

// Add item
func addItem(items []Item, item Item) []Item {

	if item.ID == 0 {
		fmt.Println("Invalid ID")
		return items
	}

	if item.Name == "" {
		fmt.Println("Name is empty")
		return items
	}

	items = append(items, item)

	return items
}

// Update item
func updateItem(items []Item, id int, quantity int, price float64) {

	for i := 0; i < len(items); i++ {

		if items[i].ID == id {
			items[i].Quantity = quantity
			items[i].Price = price

			fmt.Println("Item updated")
			return
		}
	}

	fmt.Println("Item not found")
}

// Remove item
func removeItem(items []Item, id int) []Item {

	for i := 0; i < len(items); i++ {

		if items[i].ID == id {

			items = append(items[:i], items[i+1:]...)

			fmt.Println("Item removed")
			return items
		}
	}

	fmt.Println("Item not found")

	return items
}

// Calculate total
func calculateTotal(items []Item) float64 {

	total := 0.0

	for i := 0; i < len(items); i++ {

		total = total + items[i].Quantity*items[i].Price
	}

	return total
}

// Goroutine function
func summary(items []Item, ch chan float64) {

	total := calculateTotal(items)

	ch <- total
}

func main() {

	// Empty slice
	items := []Item{}

	// Add first item
	item1 := Item{
		ID:       1,
		Name:     "Keyboard",
		Quantity: 5,
		Price:    800,
	}

	items = addItem(items, item1)

	// Add second item
	item2 := Item{
		ID:       2,
		Name:     "Mouse",
		Quantity: 10,
		Price:    500,
	}

	items = addItem(items, item2)

	// Display items
	fmt.Println("Inventory:")

	for i := 0; i < len(items); i++ {

		fmt.Println(
			items[i].ID,
			items[i].Name,
			items[i].Quantity,
			items[i].Price,
		)
	}

	// Update item
	updateItem(items, 1, 8, 750)

	// Remove item
	items = removeItem(items, 2)

	// Show remaining items
	fmt.Println("After update and remove:")

	for i := 0; i < len(items); i++ {

		fmt.Println(
			items[i].ID,
			items[i].Name,
			items[i].Quantity,
			items[i].Price,
		)
	}

	// Channel
	ch := make(chan float64)

	// Goroutine
	go summary(items, ch)

	// Receive value from channel
	total := <-ch

	fmt.Println("Total items:", len(items))
	fmt.Println("Total value:", total)
}
