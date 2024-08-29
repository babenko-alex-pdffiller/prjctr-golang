package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type Animal struct {
	Name   string
	Health int
	Hunger int
	State  string
}

type Cage struct {
	ID     int
	IsOpen bool
}

type Feeder struct {
	ID     int
	IsFull bool
}

func main() {
	animals := getAnimals()
	animalsChannel := make(chan Animal)
	for _, animal := range animals {
		go collectAnimalsData(&animal, animalsChannel)
	}

	cages := getCages()
	cagesChannel := make(chan Cage)
	for _, cage := range cages {
		go collectCagesData(&cage, cagesChannel)
	}

	feeders := getFeeders()
	feedersChannel := make(chan Feeder)
	for _, feeder := range feeders {
		go collectFeedersData(&feeder, feedersChannel)
	}

	go func() {
		for status := range animalsChannel {
			fmt.Printf("Received Aminal data: %+v\n", status)
		}
	}()

	go func() {
		for cage := range cagesChannel {
			fmt.Printf("Cage #%d is open %t\n", cage.ID, cage.IsOpen)
		}
	}()

	go func() {
		for feeder := range feedersChannel {
			fmt.Printf("Feeder #%d is full %t\n", feeder.ID, feeder.IsFull)
		}
	}()

	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}

func collectAnimalsData(animal *Animal, ch chan<- Animal) {
	ch <- *animal
	fmt.Printf("Collect data for %s\n", animal.Name)
	time.Sleep(time.Second * time.Duration(rand.IntN(5)+1))
}

func getAnimals() [5]Animal {
	animalNames := []string{"Lion", "Elephant", "Giraffe", "Zebra", "Monkey"}
	animals := [5]Animal{}

	for i, name := range animalNames {
		animals[i] = Animal{
			Name:   name,
			Health: rand.IntN(100),
			Hunger: rand.IntN(100),
			State:  []string{"Happy", "Sad", "Angry"}[rand.IntN(3)],
		}
	}

	return animals
}

func collectCagesData(cage *Cage, ch chan<- Cage) {
	ch <- *cage
	fmt.Printf("Collect Cage data #%d\n", cage.ID)
	time.Sleep(time.Second * time.Duration(rand.IntN(3)+1))
}

func getCages() [5]Cage {
	cages := [5]Cage{}

	for i := range [5]int{} {
		cages[i] = Cage{
			ID:     i + 1,
			IsOpen: rand.IntN(2) == 1,
		}
	}

	return cages
}

func collectFeedersData(feeder *Feeder, ch chan<- Feeder) {
	ch <- *feeder
	fmt.Printf("Collect Feeder data #%d\n", feeder.ID)
	time.Sleep(time.Second * time.Duration(rand.IntN(3)+1))
}

func getFeeders() [5]Feeder {
	feeders := [5]Feeder{}

	for i := range [5]int{} {
		feeders[i] = Feeder{
			ID:     i + 1,
			IsFull: rand.IntN(2) == 1,
		}
	}

	return feeders
}
