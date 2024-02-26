package main

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"time"
)

/*
Things that work:
- Adding routes from parent1, registering them as visited.
- Adding routes from parent2 with patients that are not visited in any of the routes from parent1.


Potential issues:
 - does not enforce return time constriant
 - sometimes a good portion of nurse routes are empty
*/

// Multi Parent Insertion Crossover operator
func mpic(allParents []Individual, numParents int, instance Instance, crossoverRate float64) Individual {
	var iteration int = 0

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	randomIndex := random.Intn(len(allParents))
	parent1 := allParents[randomIndex]

	for iteration < numParents {
		offspring := Individual{
			Fitness: math.Inf(1),
			Age:     0,
			Routes:  make([]Route, 0),
		}

		parent2index := getOtherParentIndex(randomIndex, len(allParents))
		parent2 := allParents[parent2index]

		visitedPatients := []int{}

		for _, route := range parent1.Routes {
			randomNum := random.Float64()
			if randomNum < crossoverRate {
				offspring.Routes = append(offspring.Routes, route)
				visitedPatients = registerPatients(route, visitedPatients)
			}
		}

		for _, route := range parent2.Routes {
			alreadyVisited := false
			for _, p := range route.Patients {
				if slices.Contains(visitedPatients, p.ID) {
					alreadyVisited = true
				}
			}
			if !alreadyVisited {
				if len(offspring.Routes) < instance.NbrNurses {
					offspring.Routes = append(offspring.Routes, route)
					visitedPatients = registerPatients(route, visitedPatients)
				}
			}
		}

		for pID := 1; pID < 101; pID++ {

			if !slices.Contains(visitedPatients, pID) {
				patientAdded := false
				for index, route := range offspring.Routes {
					feasibleRoute, ok := route.canAddPatient(pID, instance)

					if ok {
						if !slices.Contains(visitedPatients, pID) {

							offspring.Routes[index] = feasibleRoute
							visitedPatients = append(visitedPatients, pID)
							patientAdded = true

						}
					} 
				}

				if !patientAdded {
					if len(offspring.Routes) < instance.NbrNurses {
						newRoute := initalizeOneRoute(instance)
						newRoute.visitPatient(instance.getPatientAtID(pID), instance)
						offspring.addRoute(newRoute)
						visitedPatients = append(visitedPatients, pID)
					} 
				}
			}
		}

		parent1 = offspring
		iteration++

	}
	// Intialize empty routes for the ones who did not get any patients.
	for len(parent1.Routes) < instance.NbrNurses {
		notOnDuty := initalizeOneRoute(instance)
		parent1.addRoute(notOnDuty)
	}
	parent1.calculateFitness(instance)

	return parent1
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

func checkForDuplicates(individual Individual) bool {
	allPatients := make([]int, 0)
	for _, route := range individual.Routes {
		for _, patient := range route.Patients {
			if slices.Contains(allPatients, patient.ID) {
				fmt.Println("duplicateid", patient.ID)
				return true
			}
			allPatients = append(allPatients, patient.ID)
		}
	}

	return false

}

func hasDuplicates(slice []int) bool {
	occurrences := make(map[int]bool)
	for _, value := range slice {
		if _, exists := occurrences[value]; exists {
			// Duplicate found
			return true
		}
		occurrences[value] = true
	}
	// No duplicates found
	return false
}

func getOtherParentIndex(index int, n int) int {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	for {
		num := random.Intn(n)
		if num != index {
			return num
		}
	}
	return 0
}
