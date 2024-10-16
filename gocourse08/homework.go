package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
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
	data       []AnimalData[any]
	activeGprs bool
	channel    chan<- bool
}

func (s *Sender) LocalData() []AnimalData[any] {
	data := s.data
	s.data = nil
	return data
}

func (s *Sender) Send(animalData AnimalData[any]) {
	fmt.Printf("Send %s data to server\n", animalData.Kind)
}

func (s *Sender) ActivateGprs() {
	s.activeGprs = true
	s.channel <- s.activeGprs
}

func (s *Sender) IsActivateGprs() bool {
	return s.activeGprs
}

func (s *Sender) AddData(animalData AnimalData[any]) {
	s.data = append(s.data, animalData)
}

func main() {
	bear := NewAnimal("bear", 48, 37.5)
	cat := NewAnimal("cat", 160, 38.5)
	gorilla := NewAnimal("primat", 50, 36.5)

	bearData := makeAnimalData(bear)
	catData := makeAnimalData(cat)
	gorillaData := makeAnimalData(gorilla)

	channelAnimals := make(chan AnimalData[any])

	channelGprs := make(chan bool)
	sender := Sender{activeGprs: false, channel: channelGprs}
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go func(channelAnimals <-chan AnimalData[any], channelGprs <-chan bool) {
		defer wg.Done()
		//		mu := sync.Mutex{}
		activateGprs := false
		for {
			//			mu.Lock()
			select {
			case <-ctx.Done():
				return
			case animalData := <-channelAnimals:
				if activateGprs == true {
					fmt.Printf("Gprs active, Sending data for %s: %+v\n", animalData.Kind, animalData)
					sender.Send(animalData)
				} else {
					fmt.Printf("Gprs unactive, save data to local: %+v\n", animalData)
					sender.AddData(animalData)
				}
			case isActivateGprs := <-channelGprs:
				activateGprs = isActivateGprs
				for _, data := range sender.LocalData() {
					sender.Send(data)
				}
			}
			time.Sleep(time.Second)
			//			mu.Unlock()
		}
	}(channelAnimals, channelGprs)

	channelAnimals <- bearData
	channelAnimals <- catData
	channelAnimals <- gorillaData

	time.Sleep(time.Second)
	// GPRS signal is activated
	sender.ActivateGprs()

	bear.Pulse = 50
	cat.Temperature = 38.8
	gorilla.Pulse = 48

	bearData = makeAnimalData(bear)
	catData = makeAnimalData(cat)
	gorillaData = makeAnimalData(gorilla)

	channelAnimals <- bearData
	channelAnimals <- catData
	channelAnimals <- gorillaData

	cancel()
	wg.Wait()
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
