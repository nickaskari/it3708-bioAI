package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
	"sort"
)

// Reads file at filename and returns JSON.
func readJSON(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}
	return data, err
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

// Used for sorting
type Pair struct {
    Value int
    Index int
}

// Sorts an array. Returns array, and original indexes in an other
func sortWithReflection(a []int) ([]int, []int) {
    n := len(a)
    pairs := make([]Pair, n)
    for i, v := range a {
        pairs[i] = Pair{v, i}
    }

    sort.Slice(pairs, func(i, j int) bool {
        return pairs[i].Value > pairs[j].Value
    })

    sortedA := make([]int, n)
    reflectedB := make([]int, n)
    for i, pair := range pairs {
        sortedA[i] = pair.Value
        reflectedB[i] = pair.Index
    }

    return sortedA, reflectedB
}
