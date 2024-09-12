package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

// Тест роботи сенсорів та центральної системи
func TestSensorAndCentralSystem(t *testing.T) {
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
	cancel()

	time.Sleep(2 * time.Second)

	if len(cs.data) == 0 {
		t.Errorf("Expected central system to have some data, but found none")
	}

	for _, sensor := range sensors {
		if len(cs.data[sensor.Type]) == 0 {
			t.Errorf("Expected data for sensor type %s, but found none", sensor.Type)
		}
	}
}
