package main

import (
	"math"
	"math/rand"
	"time"
	"slices"
)

// Multi Parent Insertion Crossover operator
func mpic(allParents []Individual, numParents int, instance Instance, crossoverRate float64) Individual {
	var iteration int = 0

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	// Choose one random parent
	randomIndex := random.Intn(len(allParents))
	parent1 := allParents[randomIndex]

	for iteration < numParents {
		offspring := Individual{
			Fitness: math.Inf(1),
			Age:     0,
			Routes:  make([]Route, 0),
		}
		parent2 := allParents[random.Intn(len(allParents))]

		visitedPatients := make([]int, 0)

		for _, route := range parent1.Routes {
			randomNum := random.Float64()
			if randomNum < crossoverRate {
				offspring.addRoute(route)
				visitedPatients = registerPatients(route, visitedPatients)
			}
		}

		for _, route := range parent2.Routes {
			// check if can add route, based on visited patients			
			if !checkAlreadyVisited(route.extractAllVisitedPatients(), visitedPatients) {
				offspring.addRoute(route)
			}
		}
		

		for _, pID := range extractUnvisitedPatients(visitedPatients, instance) {
			
			patientAdded := false
			for _, route := range offspring.Routes {
				//if pId.CanBeFeasiblyAdded(route) { 
		}
	}
}



// Adds patients from route to patient (ID's only) array. Returns patients array.
func registerPatients(route Route, patients []int) []int {
	routePatients := route.Patients 
	for _, p := range routePatients {
		patients = append(patients, p.ID)
	}
	return patients
}

// Checks if the patient ID of one route is already visited in another patient ID array
func checkAlreadyVisited(routePatients []int, visitedPatients []int) bool {
	for _, pID := range routePatients {
		if slices.Contains(visitedPatients, pID) {
			return true
		}
	}
	return false
}

// Extract unvisited patients from visitied patients. Returns []int of patient ID's that are unvisited.
func extractUnvisitedPatients(visitedPatients []int, instance Instance) []int {
	allPatients := instance.getPatients()
	unvistedPatients := make([]int, 0)

	for _, patient := range allPatients {
		if !slices.Contains(visitedPatients, patient.ID) {
			unvistedPatients = append(unvistedPatients, patient.ID)
		}
	}

	return unvistedPatients
}


