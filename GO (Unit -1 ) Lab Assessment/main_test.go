package main

import "testing"

func TestAddItem(t *testing.T) {

	items := []Item{}

	item := Item{
		ID:       1,
		Name:     "Keyboard",
		Quantity: 5,
		Price:    800,
	}

	items = addItem(items, item)

	if len(items) != 1 {
		t.Error("Item was not added")
	}
}

func TestUpdateItem(t *testing.T) {

	items := []Item{
		{
			ID:       1,
			Name:     "Keyboard",
			Quantity: 5,
			Price:    800,
		},
	}

	updateItem(items, 1, 10, 700)

	if items[0].Quantity != 10 {
		t.Error("Quantity was not updated")
	}
}

func TestRemoveItem(t *testing.T) {

	items := []Item{
		{
			ID:       1,
			Name:     "Keyboard",
			Quantity: 5,
			Price:    800,
		},
		{
			ID:       2,
			Name:     "Mouse",
			Quantity: 10,
			Price:    500,
		},
	}

	items = removeItem(items, 1)

	if len(items) != 1 {
		t.Error("Item was not removed")
	}
}
