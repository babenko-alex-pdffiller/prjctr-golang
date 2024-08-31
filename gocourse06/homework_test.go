package main

import (
	"sync"
	"testing"
	"time"
)

func mockCollectAnimalsData(a *Animal, ch chan<- Animal) {
	ch <- *a
}

func mockCollectCagesData(c *Cage, ch chan<- Cage) {
	ch <- *c
}

func mockCollectFeedersData(f Feeder, ch chan<- Feeder) {
	ch <- f
}

func TestProcessData(t *testing.T) {
	// Arrange
	collectAnimalsData = mockCollectAnimalsData
	collectCagesData = mockCollectCagesData
	collectFeedersData = mockCollectFeedersData

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
	go collectAnimalsData(animal, animalChannel)

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
