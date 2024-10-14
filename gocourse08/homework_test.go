package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestAnimalInitialization(t *testing.T) {
	//Arrange & Act
	bear := NewAnimal("bear", 48, 37.5)
	cat := NewAnimal("cat", 160, 38.5)
	gorilla := NewAnimal("primat", 50, 36.5)

	//Assert
	if bear.Collar.AnimalKind != "bear" {
		t.Errorf("Expected bear, got %s", bear.Collar.AnimalKind)
	}

	if cat.Collar.AnimalKind != "cat" {
		t.Errorf("Expected cat, got %s", cat.Collar.AnimalKind)
	}

	if gorilla.Collar.AnimalKind != "primat" {
		t.Errorf("Expected primat, got %s", gorilla.Collar.AnimalKind)
	}
}

func TestSender_AddData(t *testing.T) {
	//Arrange
	channelGprs := make(chan bool)
	sender := Sender{activeGprs: false, channel: channelGprs}

	cat := NewAnimal("cat", 160, 38.5)
	catData := makeAnimalData(cat)
	//Act
	sender.AddData(catData)
	//Assert
	if len(sender.data) != 1 {
		t.Errorf("Expected 1 data entry, got %d", len(sender.data))
	}

	if sender.data[0].Kind != "cat" {
		t.Errorf("Expected cat data, got %s", sender.data[0].Kind)
	}
}

func TestSender_ActivateGprs(t *testing.T) {
	//Arrange
	channelGprs := make(chan bool, 1)
	sender := Sender{activeGprs: false, channel: channelGprs}
	//Act
	sender.ActivateGprs()

	isActivated := <-channelGprs
	//Assert
	if !isActivated {
		t.Errorf("Expected GPRS to be activated, but it was not")
	}
}

func TestGoroutine_SendData(t *testing.T) {
	//Arrange
	channelAnimals := make(chan AnimalData[any], 3)
	channelGprs := make(chan bool, 1)
	sender := Sender{activeGprs: false, channel: channelGprs}

	bear := NewAnimal("bear", 48, 37.5)
	cat := NewAnimal("cat", 160, 38.5)
	gorilla := NewAnimal("primat", 50, 36.5)

	channelAnimals <- makeAnimalData(bear)
	channelAnimals <- makeAnimalData(cat)
	channelAnimals <- makeAnimalData(gorilla)
	//Act
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		activateGprs := false
		for {
			select {
			case <-ctx.Done():
				return
			case animalData := <-channelAnimals:
				if activateGprs {
					sender.Send(animalData)
				} else {
					sender.AddData(animalData)
				}
			case isActivateGprs := <-channelGprs:
				activateGprs = isActivateGprs
				for _, data := range sender.LocalData() {
					sender.Send(data)
				}
			}
		}
	}()

	sender.ActivateGprs()

	time.Sleep(time.Second)
	cancel()
	wg.Wait()
	//Assert
	if len(sender.data) != 0 {
		t.Errorf("Expected no local data, but got %d entries", len(sender.data))
	}
}
