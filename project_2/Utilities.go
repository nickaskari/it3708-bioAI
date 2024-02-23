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

// getRandomRoute picks a random route from the map and returns the route along with its index.
func getRandomRoute(routesMap map[int]Route) (Route, int) {
	keys := make([]int, 0, len(routesMap))
	for k := range routesMap {
		keys = append(keys, k)
	}

	// Pick a random key.
	randomIndex := rand.Intn(len(keys))
	randomKey := keys[randomIndex]

	return routesMap[randomKey], randomKey
}

// Pares a string to an int and returns the int.
func strToInt(str string) int {

	n, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println(err)
	}
	return n
}

// Prints a solution, or an individual
func printSolution(individual Individual, instance Instance) {
	nurseCapacity := instance.CapacityNurse
	depotReturnTime := instance.Depot.ReturnTime
	objectiveValue := individual.Fitness

	fmt.Println("Age of individual:", individual.Age)
	fmt.Println("Nurse capacity:", nurseCapacity)
	fmt.Println("Depot return time:", depotReturnTime)
	printDivider(150, "-")

	const maxSequenceLength = 10000

	for i, route := range individual.Routes {
		nurseIdentifier := fmt.Sprintf("Nurse %-3d", i+1)
		routeDuration := fmt.Sprintf("%-6.2f", route.CurrentTime)
		coveredDemand := 0
		patientSequence := ""
		if len(route.Patients) > 0 {
			patientSequence += "D (0)"
			for _, patient := range route.Patients {
				sequencePart := fmt.Sprintf(" -> %d (%.2f-%.2f) [%d-%d]",
					patient.ID, float64(patient.VisitTime), float64(patient.LeavingTime), patient.StartTime, patient.EndTime)
				if len(patientSequence)+len(sequencePart) > maxSequenceLength {
					patientSequence += " ..."
					break
				}
				patientSequence += sequencePart
				coveredDemand += patient.Demand
			}
			patientSequence += fmt.Sprintf(" -> D (%.2f)", route.CurrentTime)
		} else {
			patientSequence = "NOT ON DUTY"
		}

		coveredDemandStr := fmt.Sprintf("%-4d", coveredDemand)

		fmt.Printf("%-10s %-10s %-5s %-s\n", nurseIdentifier, routeDuration, coveredDemandStr, patientSequence)
	}

	printDivider(150, "-")
	fmt.Println("Objective value (total duration):", objectiveValue)
}

// Prints out a divider (for example: "-----") of desired length
func printDivider(length int, dividerChar string) {
	for i := 0; i < length; i++ {
		fmt.Print(dividerChar)
	}
	fmt.Println()
}

// takes in a route array, and outputs a dictionary with the keys being the indexes of the array.
func sliceToMap(routes []Route) map[int]Route {
	availableRoutes := make(map[int]Route)
	for index, route := range routes {
		availableRoutes[index] = route
	}
	return availableRoutes
}

// Deletes key from route dictionary. Note the keys can be though of as indexes.
func removeRouteFromMap(routesMap map[int]Route, index int) map[int]Route {
	delete(routesMap, index)
	return routesMap
}

// Chooses "size" amount of elements from "slice" randomly.
func chooseRandomUnique[T any](slice []T, size int) []T {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	r.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})

	if size > len(slice) {
		size = len(slice) 
	}
	return slice[:size]
}

// A struct to register constriant violations.
type Violation struct {
	Count int
	Example string
}

// Counts a violation.
func (v *Violation) countViolation() {
	v.Count++
}

// Registers an example of the violation.
func (v *Violation) registerExample(example string) {
	v.Example = example
}
