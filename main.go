package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Car struct to store car details
type Car struct {
	RegistrationNumber string
}

// ParkingLot struct to manage the parking system
type ParkingLot struct {
	Capacity  int
	Slots     map[int]*Car
	SlotCount int
}

// NewParkingLot creates a new parking lot with given capacity
func NewParkingLot(capacity int) *ParkingLot {
	return &ParkingLot{
		Capacity:  capacity,
		Slots:     make(map[int]*Car),
		SlotCount: 0,
	}
}

// CreateParkingLot initializes the parking lot with given capacity
func (pl *ParkingLot) CreateParkingLot(capacity int) string {
	pl.Capacity = capacity
	pl.Slots = make(map[int]*Car)
	return fmt.Sprintf("Created a parking lot with %d slots", capacity)
}

// GetNextAvailableSlot returns the next available slot number
func (pl *ParkingLot) GetNextAvailableSlot() int {
	for i := 1; i <= pl.Capacity; i++ {
		if _, exists := pl.Slots[i]; !exists {
			return i
		}
	}
	return -1
}

// Park parks a car in the nearest available slot
func (pl *ParkingLot) Park(registrationNumber string) string {
	if pl.SlotCount == pl.Capacity {
		return "Sorry, parking lot is full"
	}

	nextSlot := pl.GetNextAvailableSlot()
	pl.Slots[nextSlot] = &Car{RegistrationNumber: registrationNumber}
	pl.SlotCount++

	return fmt.Sprintf("Allocated slot number: %d", nextSlot)
}

// Leave removes a car from the parking lot
func (pl *ParkingLot) Leave(registrationNumber string, hours int) string {
	for slotNumber, car := range pl.Slots {
		if car.RegistrationNumber == registrationNumber {
			delete(pl.Slots, slotNumber)
			pl.SlotCount--

			// Calculate parking charge
			charge := 10 // Base charge for first 2 hours
			if hours > 2 {
				charge += 10 * (hours - 2) // Additional $10 per hour after first 2 hours
			}

			return fmt.Sprintf("Registration number %s with Slot Number %d is free with Charge $%d", registrationNumber, slotNumber, charge)
		}
	}
	return fmt.Sprintf("Registration number %s not found", registrationNumber)
}

// Status prints the current status of the parking lot
func (pl *ParkingLot) Status() string {
	if pl.SlotCount == 0 {
		return "Parking lot is empty"
	}

	status := "Slot No. Registration No.\n"
	for i := 1; i <= pl.Capacity; i++ {
		if car, exists := pl.Slots[i]; exists {
			status += fmt.Sprintf("%d %s\n", i, car.RegistrationNumber)
		}
	}
	return status
}

// ProcessCommand processes a single command
func (pl *ParkingLot) ProcessCommand(command string) string {
	tokens := strings.Split(command, " ")
	cmd := strings.ToLower(tokens[0]) // Convert command to lowercase for case-insensitive comparison

	switch cmd {
	case "create_parking_lot":
		capacity, _ := strconv.Atoi(tokens[1])
		return pl.CreateParkingLot(capacity)
	case "park":
		return pl.Park(tokens[1])
	case "leave":
		hours, _ := strconv.Atoi(tokens[2])
		return pl.Leave(tokens[1], hours)
	case "status":
		return pl.Status()
	default:
		return fmt.Sprintf("Invalid command: %s", command)
	}
}

func main() {
	// Check if file name is provided
	if len(os.Args) != 2 {
		fmt.Println("Please provide input file name")
		os.Exit(1)
	}

	// Open input file
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Initialize parking lot
	var parkingLot *ParkingLot

	// Read and process commands from file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		command := scanner.Text()
		
		// Skip empty lines
		if len(command) == 0 {
			continue
		}

		// Initialize parking lot if not already done
		if parkingLot == nil {
			tokens := strings.Split(command, " ")
			if tokens[0] == "create_parking_lot" {
				capacity, _ := strconv.Atoi(tokens[1])
				parkingLot = NewParkingLot(capacity)
				fmt.Println(fmt.Sprintf("Created a parking lot with %d slots", capacity))
				continue
			} else {
				fmt.Println("Please create parking lot first")
				continue
			}
		}

		// Process the command and print the result
		result := parkingLot.ProcessCommand(command)
		fmt.Println(result)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}
}