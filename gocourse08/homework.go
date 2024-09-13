package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type AnimalData[T any] struct {
	Kind        string
	Pulse       int
	Temperature float64
	Breathing   T
	Sound       T
}

type Collar struct {
	Breathing  int
	Sound      string
	AnimalKind string
}

func (c *Collar) Init(animal Animal) {
	if animal.Pulse > 140 && animal.Pulse < 220 && animal.Temperature >= 38 && animal.Temperature <= 39 {
		c.AnimalKind = "cat"
		c.Breathing = rand.IntN(50) + 50
		c.Sound = "meow"
	} else if animal.Pulse >= 40 && animal.Pulse <= 50 && animal.Temperature >= 37 && animal.Temperature <= 38 {
		c.AnimalKind = "bear"
		c.Breathing = rand.IntN(40) + 10
		c.Sound = "roar"
	} else {
		c.Breathing = rand.IntN(50) + 20
		c.Sound = "scream"
		c.AnimalKind = "primat"
	}

	fmt.Println(c.AnimalKind)
}

type Animal struct {
	kind        string
	Pulse       int
	Temperature float64
	Collar      Collar
}

func NewAnimal(kind string, pulse int, temp float64) Animal {
	animal := Animal{
		kind:        kind,
		Pulse:       pulse,
		Temperature: temp,
		Collar:      Collar{},
	}

	animal.Collar.Init(animal)

	return animal
}

func (a *Animal) SetCollar(c Collar) {
	a.Collar = c
}

type Sender struct {
	data         []AnimalData[any]
	isActiveGprs bool
}

func (s *Sender) Send(channel <-chan AnimalData[any]) {
	for {
		animalData := <-channel
		if s.isActiveGprs {
			fmt.Printf("Sending data for %s: %+v\n", animalData.Kind, animalData)
		} else {
			fmt.Printf("Gprs unactive, save data to local: %+v\n", animalData)
			s.data = append(s.data, animalData)
		}
	}
}

func (s *Sender) SendAllData(channel chan<- AnimalData[any]) {
	if s.isActiveGprs {
		for _, data := range s.data {
			fmt.Printf("Send local data: %+v\n", data)
			channel <- data
		}
		s.data = []AnimalData[any]{}
	}
}

func (s *Sender) ActivateGprs() {
	s.isActiveGprs = true
}

func main() {
	bear := NewAnimal("bear", 48, 37.5)
	cat := NewAnimal("cat", 160, 38.5)
	gorilla := NewAnimal("primat", 50, 36.5)

	bearData := makeAnimalData(bear)
	catData := makeAnimalData(cat)
	gorillaData := makeAnimalData(gorilla)
	channel := make(chan AnimalData[any])

	sender := Sender{isActiveGprs: false}

	go sender.Send(channel)

	channel <- bearData
	channel <- catData
	channel <- gorillaData

	time.Sleep(time.Second)
	// GPRS signal is activated
	sender.ActivateGprs()
	sender.SendAllData(channel)

	time.Sleep(time.Second)
}

func makeAnimalData(animal Animal) AnimalData[any] {
	animalData := AnimalData[any]{
		Kind:        animal.kind,
		Pulse:       animal.Pulse,
		Temperature: animal.Temperature,
	}

	animalData = collectBreath[any](animalData, animal.Collar.Breathing)
	animalData = collectSound[any](animalData, animal.Collar.Sound)

	return animalData
}

func collectSound[T any](animalData AnimalData[T], soundData T) AnimalData[T] {
	animalData.Sound = soundData

	return animalData
}

func collectBreath[T any](animalData AnimalData[T], breathingData T) AnimalData[T] {
	animalData.Breathing = breathingData

	return animalData
}
