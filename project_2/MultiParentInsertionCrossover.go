package main

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strconv"
	"time"
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

		routesAdded := 0
		for _, route := range parent1.Routes {
			randomNum := random.Float64()
			if randomNum < crossoverRate {
				offspring.addRoute(route)
				visitedPatients = registerPatients(route, visitedPatients)
				routesAdded++
			}
		}
		fmt.Println("Routes added from parent1 =", routesAdded)

		routesAddedFromparent2 := 0
		for _, route := range parent2.Routes {
			// check if can add route, based on visited patients
			if !checkAlreadyVisited(route.extractAllVisitedPatients(), visitedPatients) &&
				(len(offspring.Routes) < instance.NbrNurses) {
				offspring.addRoute(route)
				routesAddedFromparent2++
			}
		}
		fmt.Println("Routes added from parent2 =", routesAddedFromparent2)

		for _, pID := range extractUnvisitedPatients(visitedPatients, instance) {

			patientAdded := false
			for _, route := range offspring.Routes { // this is fucked
				feasibleRoute, ok := route.canAddPatient(pID, instance)
				if ok {
					if len(offspring.Routes) < instance.NbrNurses {
						offspring.addRoute(feasibleRoute)
						visitedPatients = registerPatients(route, visitedPatients)
						patientAdded = true
					}
				}
			}
			//fmt.Println("NUMBER OF ROUTES =", len(offspring.Routes))

			if !patientAdded && (len(offspring.Routes) < instance.NbrNurses) {
				newRoute := initalizeOneRoute(instance)
				newRoute.visitPatient(instance.Patients[strconv.Itoa(pID)], instance)
				offspring.addRoute(newRoute)
				visitedPatients = registerPatients(newRoute, visitedPatients)
			}
		}

		parent1 = offspring
		iteration++
		fmt.Println("LENGTH OF TOTAL OFFSPRING ROUTE ", len(offspring.Routes))
		fmt.Println("LENGTH OF TOTAL unvisit ", len(offspring.Routes))
		fmt.Println("LENTH OF VISITED PATIENTS", len(visitedPatients))
		fmt.Println("\n\niteration =", iteration, "AND numPArents =", numParents, "\n")
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
