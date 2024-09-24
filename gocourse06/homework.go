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
	ID            int
	IsOpen        bool
	openCloseChan chan<- map[int]bool
	wg            *sync.WaitGroup
}

func (c *Cage) ToggleCageDoor() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		if c.IsOpen {
			c.IsOpen = false
		} else {
			c.IsOpen = true
		}
		c.openCloseChan <- map[int]bool{c.ID: c.IsOpen}
		time.Sleep(time.Second)
	}()
}

type Feeder struct {
	ID          int
	IsFull      bool
	fullingChan chan<- bool
	wg          *sync.WaitGroup
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
		f.fullingChan <- f.IsFull
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

	cageChannel := make(chan map[int]bool)
	cage := Cage{
		ID:            111,
		IsOpen:        false,
		openCloseChan: cageChannel,
		wg:            wg,
	}

	cage.ToggleCageDoor()

	feederChannel := make(chan bool)
	feeder := Feeder{
		ID:          11,
		IsFull:      false,
		fullingChan: feederChannel,
		wg:          wg,
	}

	feeder.Feed()
	feederState := <-feederChannel
	fmt.Printf("Received Feeder status: is full %t\n", feederState)

	time.Sleep(time.Second)

	feeder.ToFill()
	feeder.Feed()
	feederState = <-feederChannel
	fmt.Printf("Received Feeder status: is full %t\n", feederState)

	cage.ToggleCageDoor()

	go func() {
		for status := range animalsChannel {
			fmt.Printf("Received Animal data: %+v\n", status)
		}
	}()

	go func() {
		for cage := range cageChannel {
			for id, isOpen := range cage {
				fmt.Printf("Cage #%d is open %t\n", id, isOpen)
			}
		}
	}()

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
