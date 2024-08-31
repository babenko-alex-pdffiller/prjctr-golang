package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
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
	var wg sync.WaitGroup
	processData(&wg)
}

func processData(wg *sync.WaitGroup) {
	animals := getAnimals()
	animalsChannel := make(chan Animal)
	wg.Add(len(animals))
	for _, animal := range animals {
		go func(animal Animal) {
			defer wg.Done()
			collectAnimalsData(&animal, animalsChannel)
		}(animal)
	}

	cages := getCages()
	cagesChannel := make(chan Cage)
	wg.Add(len(cages))
	for _, cage := range cages {
		go func(cage Cage) {
			defer wg.Done()
			collectCagesData(&cage, cagesChannel)
		}(cage)
	}

	feeders := getFeeders()
	feedersChannel := make(chan Feeder)
	wg.Add(len(feeders))
	for _, feeder := range feeders {
		go func(feeder Feeder) {
			defer wg.Done()
			collectFeedersData(feeder, feedersChannel)
		}(feeder)
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

	wg.Wait()

	close(animalsChannel)
	close(cagesChannel)
	close(feedersChannel)

	fmt.Println("All goroutines finished")
}

var collectAnimalsData = func(animal *Animal, ch chan<- Animal) {
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

var collectCagesData = func(cage *Cage, ch chan<- Cage) {
	ch <- *cage
	fmt.Printf("Collect Cage data #%d\n", cage.ID)
	time.Sleep(time.Second * time.Duration(rand.IntN(3)+1))
}

func getCages() [5]Cage {
	cages := [5]Cage{}

	for i := range 5 {
		cages[i] = Cage{
			ID:     i + 1,
			IsOpen: rand.IntN(2) == 1,
		}
	}

	return cages
}

var collectFeedersData = func(feeder Feeder, ch chan<- Feeder) {
	ch <- feeder
	fmt.Printf("Collect Feeder data #%d\n", feeder.ID)
	time.Sleep(time.Second * time.Duration(rand.IntN(3)+1))
}

func getFeeders() [5]Feeder {
	feeders := [5]Feeder{}

	for i := range 5 {
		feeders[i] = Feeder{
			ID:     i + 1,
			IsFull: rand.IntN(2) == 1,
		}
	}

	return feeders
}
