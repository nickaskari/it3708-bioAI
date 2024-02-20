package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
    "time"
)

// Reads file at filename and returns JSON.
func readJSON(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}
	return data, err
}

// Reads problem instance from JSON and returns an Instance object.
func getProblemInstance(filename string) Instance {
	data, err := readJSON(filename)
	if err != nil {
		fmt.Println("failed to read JSON data: %w", err)
		os.Exit(1)
	}

	if data == nil {
		// Check for empty data
		fmt.Println("empty JSON data")
		os.Exit(1)
	}

	var instance Instance
	err = json.Unmarshal(data, &instance)
	if err != nil {
		fmt.Println("error unmarshaling JSON: %w", err)
		os.Exit(1)
	}
	return instance
}

// A struct to perform random math operations
type Randomizer struct {
	*rand.Rand
}

// NewRandomizer creates a new Randomizer (constructor)
func NewRandomizer() *Randomizer {
	return &Randomizer{
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Returns a random route from an array of routes, as well as the route's index in the array.
func (r *Randomizer) getRandomRoute(routes []Route) (Route, int) {

	randomIndex := r.Intn(len(routes))

	return routes[randomIndex], randomIndex
}

// Pares a string to an int and returns the int.
func strToInt(str string) int {

	n, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println(err)
	}
	return n
}
