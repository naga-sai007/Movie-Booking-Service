package main

import (
	"log"
	"os"
	"os/exec"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runUserMicroservice(); err != nil {
			log.Fatalf("User Microservice error: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runMovieMicroservice(); err != nil {
			log.Fatalf("Movie Microservice error: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runTheatreMicroservice(); err != nil {
			log.Fatalf("Movie Microservice error: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runBookingMicroservice(); err != nil {
			log.Fatalf("Movie Microservice error: %v", err)
		}
	}()
	// Wait for a termination signal
	//waitForTerminationSignal()

	// Wait for goroutines to finish
	wg.Wait()
}

func runUserMicroservice() error {
	// Code to initialize user microservice, if needed
	log.Println("Starting User Microservice...")
	cmd := exec.Command("go", "run", "user-service/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runMovieMicroservice() error {
	// Code to initialize product microservice
	log.Println("Starting Movie Microservice...")
	cmd := exec.Command("go", "run", "movie-service/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runTheatreMicroservice() error {
	// Code to initialize product microservice
	log.Println("Starting Theatre Microservice...")
	cmd := exec.Command("go", "run", "theatre-service/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runBookingMicroservice() error {
	// Code to initialize product microservice
	log.Println("Starting Booking Microservice...")
	cmd := exec.Command("go", "run", "booking-service/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
