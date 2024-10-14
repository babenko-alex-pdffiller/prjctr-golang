package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

type Sensor struct {
	ID       int
	Type     string
	Value    int
	IsActive bool
}

func (s *Sensor) Run(ctx context.Context, dataChannel chan<- Sensor) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Sensor #%d %s stopped\n", s.ID, s.Type)
			return
		default:
			if s.IsActive {
				s.Value = rand.IntN(90) + 10
				dataChannel <- *s
				fmt.Printf("Sensor %s sent data: %d\n", s.Type, s.Value)
			}
			time.Sleep(time.Second)
		}
	}
}

type CentralSystem struct {
	data map[string]map[int]int
	mu   sync.Mutex
	wg   *sync.WaitGroup
}

func (cs *CentralSystem) Run(ctx context.Context, dataChannel <-chan Sensor) {
	for {
		cs.mu.Lock()
		select {
		case <-ctx.Done():
			cs.wg.Wait()
			fmt.Println("Central system shutting down, waiting for all records to be saved.")
			return
		case sensor := <-dataChannel:
			if cs.data[sensor.Type] == nil {
				cs.data[sensor.Type] = make(map[int]int)
			}
			cs.data[sensor.Type][int(time.Now().UnixNano())] = sensor.Value
			fmt.Printf("Central system added to memory %d\n", sensor.Value)
		}
		cs.mu.Unlock()
	}
}

func (cs *CentralSystem) Data() map[string]map[int]int {
	return cs.data
}

func main() {
	dataChannel := make(chan Sensor, 3)
	wg := sync.WaitGroup{}

	sensors := []Sensor{
		{ID: 1, Type: "Temperature", IsActive: true},
		{ID: 2, Type: "Brightness", IsActive: true},
		{ID: 3, Type: "Humidity", IsActive: true},
	}

	cs := CentralSystem{
		data: make(map[string]map[int]int),
		wg:   &wg,
	}

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(len(sensors))
	for i, _ := range sensors {
		go func() {
			defer wg.Done()
			sensors[i].Run(ctx, dataChannel)
		}()
	}

	go cs.Run(ctx, dataChannel)

	time.Sleep(3 * time.Second)

	fmt.Println("Started technical maintenance")
	cancel()

	fmt.Println("Maintenance completed. System is shut down.")

	time.Sleep(3 * time.Second)

	for key, item := range cs.Data() {
		for time, data := range item {
			fmt.Printf("Sensor %s, time %d, data %d\n", key, time, data)
		}
	}
}
