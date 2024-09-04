package main

import (
	"sync"
	"testing"
	"time"
)

func TestProcessData(t *testing.T) {
	// Arrange
	var wg sync.WaitGroup

	// Act
	processData(&wg)

	// Assert
	if wg != (sync.WaitGroup{}) {
		t.Error("WaitGroup not finished")
	}
}

func TestCollectAnimalsData(t *testing.T) {
	// Arrange
	animalChannel := make(chan Animal, 1)
	animal := &Animal{Name: "Lion"}

	// Act
	go collectAnimalsData(animal, animalChannel, time.Duration(0))

	// Assert
	select {
	case receivedAnimal := <-animalChannel:
		if receivedAnimal.Name != "Lion" {
			t.Errorf("Expected animal name 'Lion', but got '%s'", receivedAnimal.Name)
		}
	case <-time.After(time.Second):
		t.Error("Timeout while waiting for collectAnimalsData")
	}
}
