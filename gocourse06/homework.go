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
	chanel chan<- bool
	wg     *sync.WaitGroup
}

func (f *Feeder) ToFill() {
	f.IsFull = true
}

func (f *Feeder) ToEmpty() {
	f.IsFull = false
}

func (f *Feeder) Feed() {
	f.wg.Add(1)
	go func() {
		defer f.wg.Done()
		f.chanel <- f.IsFull
	}()
}

func main() {
	var wg sync.WaitGroup
	processData(&wg)
}

func processData(wg *sync.WaitGroup) {
	animals := generateAnimals()
	animalsChannel := make(chan Animal)
	wg.Add(len(animals))
	for _, animal := range animals {
		go func() {
			defer wg.Done()
			collectAnimalsData(animal, animalsChannel, time.Duration(rand.IntN(5))*time.Second)
		}()
	}

	cagesChannel := make(chan Cage)
	cage := Cage{
		ID:     111,
		IsOpen: rand.IntN(2) == 1,
	}

	go toggleCageDoor(cagesChannel)

	cagesChannel <- cage

	feederChannel := make(chan bool)
	feeder := Feeder{
		ID:     11,
		IsFull: false,
		chanel: feederChannel,
		wg:     wg,
	}

	feeder.Feed()
	feederState := <-feederChannel
	fmt.Printf("Recived Feeder status: is full %t\n", feederState)

	time.Sleep(time.Second)

	feeder.ToFill()
	feeder.Feed()
	feederState = <-feederChannel
	fmt.Printf("Recived Feeder status: is full %t\n", feederState)

	go func() {
		for status := range animalsChannel {
			fmt.Printf("Received Animal data: %+v\n", status)
		}
	}()

	go func() {
		for cage := range cagesChannel {
			fmt.Printf("Cage #%d is open %t\n", cage.ID, cage.IsOpen)
		}
	}()

	cagesChannel <- cage

	wg.Wait()

	fmt.Println("All goroutines finished")
}

func collectAnimalsData(animal Animal, ch chan<- Animal, timeout time.Duration) {
	fmt.Printf("Collect data for %s\n", animal.Name)
	ch <- animal
	time.Sleep(timeout)
}

func generateAnimals() [5]Animal {
	animalNames := []string{"Lion", "Elephant", "Giraffe", "Zebra", "Monkey"}
	animals := [5]Animal{}

	for i, name := range animalNames {
		animals[i] = Animal{
			Name:   name,
			Health: rand.IntN(100),
			Hunger: rand.IntN(100),

			State: []string{"Happy", "Sad", "Angry"}[rand.IntN(3)],
		}
	}

	return animals
}

func toggleCageDoor(ch <-chan Cage) {
	cage := <-ch
	if cage.IsOpen == true {
		cage.IsOpen = false
		fmt.Printf("Cage #%d is closed\n", cage.ID)
	} else {
		cage.IsOpen = true
		fmt.Printf("Cage #%d is opened\n", cage.ID)
	}
}
