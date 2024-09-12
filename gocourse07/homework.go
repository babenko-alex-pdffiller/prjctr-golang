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
				fmt.Printf("Sensor %d sent data: %d\n", s.ID, s.Value)
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
		select {
		case <-ctx.Done():
			cs.mu.Lock()
			defer cs.mu.Unlock()
			fmt.Println("Central system shutting down, waiting for all records to be saved.")
			cs.wg.Wait()
			return
		case sensorData := <-dataChannel:
			cs.mu.Lock()
			cs.wg.Add(1)
			go func() {
				defer cs.wg.Done()
				if cs.data[sensorData.Type] == nil {
					cs.data[sensorData.Type] = make(map[int]int)
				}
				cs.data[sensorData.Type][int(time.Now().UnixNano())] = sensorData.Value
				fmt.Printf("Central system added to memory %d\n", sensorData.Value)
				time.Sleep(3 * time.Second)
			}()
			cs.mu.Unlock()
		}
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

	for _, sensor := range sensors {
		go sensor.Run(ctx, dataChannel)
	}

	go cs.Run(ctx, dataChannel)

	time.Sleep(3 * time.Second)

	fmt.Println("Started technical maintenance")
	cancel()

	time.Sleep(3 * time.Second)

	fmt.Println("Maintenance completed. System is shut down.")

	for key, item := range cs.Data() {
		for time, data := range item {
			fmt.Printf("Sensor %s, time %d, data %d\n", key, time, data)
		}
	}
}
