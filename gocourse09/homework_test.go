package main

import (
	"testing"
)

func TestGiveOut(t *testing.T) {
	feeder := makeFeeder()

	foodsNeeded := map[string]int{
		"banana":   3,
		"meat":     2,
		"millet":   4,
		"porridge": 1,
	}

	feeder.giveOut(foodsNeeded)

	expectedFoodAmounts := map[string]int{
		"banana":   7,
		"meat":     8,
		"millet":   6,
		"porridge": 9,
	}

	for foodType, expectedAmount := range expectedFoodAmounts {
		if feeder.Foods[foodType].Amount != expectedAmount {
			t.Errorf("Incorrect food amount for %s: got %d, want %d",
				foodType, feeder.Foods[foodType].Amount, expectedAmount)
		}
	}
}

func TestCalculateFoods(t *testing.T) {
	animals := []Animal{
		{ID: 1, Type: "Gorilla", Position: 1},
		{ID: 2, Type: "Tiger", Position: 5},
		{ID: 3, Type: "Duck", Position: 12},
		{ID: 4, Type: "Lion", Position: 3},
		{ID: 5, Type: "Swan", Position: 16},
	}

	expectedFoods := map[string]int{
		"banana": 1, // Gorilla
		"meat":   2, // Tiger, Lion
		"millet": 2, // Duck, Swan
	}

	foods := calculateFoods(animals)

	for foodType, expectedAmount := range expectedFoods {
		if foods[foodType] != expectedAmount {
			t.Errorf("Incorrect food amount for %s: got %d, want %d",
				foodType, foods[foodType], expectedAmount)
		}
	}
}

func TestFindAnimalsNear(t *testing.T) {
	animals := makeAnimals()

	expectedAnimals := []Animal{
		{ID: 1, Type: "Gorilla", Position: 5},
		{ID: 4, Type: "Duck", Position: 13},
		{ID: 40, Type: "Swan", Position: 5},
	}

	nearAnimals := findAnimalsNear(&animals)

	var expected bool
	for _, expectedAnimal := range expectedAnimals {
		expected = false
		for _, animal := range nearAnimals {
			if animal.ID == expectedAnimal.ID {
				expected = true
				break
			}
		}

		if expected == false {
			t.Errorf("Animal %v not found in nearAnimals", expectedAnimal)
		}
	}
}

func TestMoveAnimalsPosition(t *testing.T) {
	animals := makeAnimals()

	originalPositions := make(map[int]int)
	for _, animal := range animals {
		originalPositions[animal.ID] = animal.Position
	}

	movedAnimals := moveAnimalsPosition(animals)

	for _, animal := range movedAnimals {
		if animal.Position < 1 || animal.Position > 20 {
			t.Errorf("Animal with ID %d moved to invalid position: %d", animal.ID, animal.Position)
		}

		if animal.Position == originalPositions[animal.ID] {
			t.Errorf("Animal with ID %d did not change its position", animal.ID)
		}
	}
}
