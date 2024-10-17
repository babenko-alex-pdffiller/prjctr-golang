package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type Animal struct {
	ID       int
	Type     string
	Position int
}

type Feeder struct {
	Foods map[string]Food
}

func (f *Feeder) giveOut(foodsAmount map[string]int) {
	for foodType, amount := range foodsAmount {
		food := f.Foods[foodType]
		if food.Amount < amount {
			fmt.Printf("Not enough food")
			continue
		}

		food.Amount -= amount
		f.Foods[foodType] = food
		fmt.Printf("Issued out %s in %d portions\n", foodType, amount)
	}
}

type Food struct {
	Amount int
	Type   string
}

func main() {
	animals := makeAnimals()

	feeder := makeFeeder()
	nearAnimals := findAnimalsNear(&animals)

	feeder.giveOut(calculateFoods(nearAnimals))

	stop := make(chan bool)

	// animals moving
	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("Animals stopped")
				return
			default:
				animals = moveAnimalsPosition(animals)
				time.Sleep(time.Second)
			}
		}
	}()

	// waiting for push the button
	fmt.Scanln()
	stop <- true

	nearAnimals = findAnimalsNear(&animals)
	feeder.giveOut(calculateFoods(nearAnimals))
}

func findAnimalsNear(animals *[]Animal) []Animal {
	nearAnimals := []Animal{}
	for _, animal := range *animals {
		if animal.Position < 10 {
			fmt.Printf("Found near %s\n", animal.Type)
			nearAnimals = append(nearAnimals, animal)
		}
	}

	return nearAnimals
}

func calculateFoods(animals []Animal) map[string]int {
	foodsAmount := make(map[string]int)
	for _, animal := range animals {
		switch animal.Type {
		case "Gorilla", "Chimpanzee":
			foodsAmount["banana"] += 1
		case "Tiger", "Lion":
			foodsAmount["meat"] += 1
		case "Duck", "Swan":
			foodsAmount["millet"] += 1
		default:
			foodsAmount["porridge"] += 1
		}
	}

	return foodsAmount
}

func makeFeeder() Feeder {
	return Feeder{
		Foods: map[string]Food{
			"banana": {
				Amount: 10,
				Type:   "banana",
			},
			"meat": {
				Amount: 10,
				Type:   "meat",
			},
			"millet": {
				Amount: 10,
				Type:   "millet",
			},
			"porridge": {
				Amount: 10,
				Type:   "porridge",
			},
		},
	}
}

func makeAnimals() []Animal {
	return []Animal{
		{
			ID:       1,
			Type:     "Gorilla",
			Position: 1,
		},
		{
			ID:       10,
			Type:     "Chimpanzee",
			Position: 15,
		},
		{
			ID:       2,
			Type:     "Tiger",
			Position: 11,
		},
		{
			ID:       20,
			Type:     "Lion",
			Position: 12,
		},
		{
			ID:       4,
			Type:     "Duck",
			Position: 3,
		},
		{
			ID:       40,
			Type:     "Swan",
			Position: 5,
		},
	}
}

func moveAnimalsPosition(animals []Animal) []Animal {
	for i := range animals {
		animals[i].Position = rand.IntN(20) + 1
	}

	fmt.Println("Animals moved")
	return animals
}
