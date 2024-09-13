package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSenderActivateGprs(t *testing.T) {
	sender := Sender{isActiveGprs: false}
	sender.ActivateGprs()

	if !sender.isActiveGprs {
		t.Error("ActivateGprs() did not set isActiveGprs to true")
	}
}

func TestMakeAnimalData(t *testing.T) {
	animal := NewAnimal("primat", 30, 36.0)
	data := makeAnimalData(animal)

	if data.Kind != "primat" {
		t.Errorf("Expected kind 'primat', got '%s'", data.Kind)
	}

	if data.Pulse != 30 {
		t.Errorf("Expected heart rate 30, got %d", data.Pulse)
	}

	if data.Temperature != 36.0 {
		t.Errorf("Expected body temperature 36.0, got %f", data.Temperature)
	}
}

func TestSendAllData(t *testing.T) {
	// Arrange
	channel := make(chan AnimalData[any])
	sender := Sender{isActiveGprs: false}
	go sender.Send(channel)

	channel <- AnimalData[any]{Kind: "cat", Pulse: 40, Temperature: 37.0}
	channel <- AnimalData[any]{Kind: "bear", Pulse: 45, Temperature: 37.5}
	channel <- AnimalData[any]{Kind: "primat", Pulse: 50, Temperature: 38.0}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Act
	time.Sleep(time.Second)
	sender.ActivateGprs()
	sender.SendAllData(channel)

	time.Sleep(time.Second)
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Assert
	expectedAnimals := []string{"cat", "bear", "primat"}
	for _, animal := range expectedAnimals {
		if !strings.Contains(output, fmt.Sprintf("Sending data for %s", animal)) {
			fmt.Println("WTF: " + output)
			t.Errorf("Expected output to contain 'Sending data for %s', but it didn't", animal)
		}
	}

	// Check if the channel is empty after sending all data
	select {
	case <-channel:
		t.Error("Channel should be empty after SendAllData")
	default:
		// Channel is empty, which is expected
	}
}
