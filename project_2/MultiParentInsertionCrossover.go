package main

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strconv"
	"time"
)
/*
Things that work:
- Adding routes from parent1, registering them as visited. 
- Adding routes from parent2 with patients that are not visited in any of the routes from parent1.


Potential issues:
- Potensielle feil med måten registerPatients brukes på
- Visited (og følgelig unvisited) patients blir ikke oppdatert riktig etter første iterasjon. Tror dette har noe å gjøre med at visitedPatients settes til 0 for hver iterasjon.
- Det er duplikater av pasienter i individet mpic returnerer.
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
		fmt.Println("visitedPatients: ", visitedPatients) 
		printSolution(offspring, instance)

		//fmt.Println("Routes added from parent1 =", routesAdded)

		routesAddedFromparent2 := 0
		for _, route := range parent2.Routes {
			// check if can add route, based on visited patients
			if !checkAlreadyVisited(route.extractAllVisitedPatients(), visitedPatients) &&
				(len(offspring.Routes) < instance.NbrNurses) {
				offspring.addRoute(route)
				visitedPatients = registerPatients(route, visitedPatients)
				routesAddedFromparent2++
			}
			fmt.Println("route.extractAllVisitedPatients()", route.extractAllVisitedPatients()) 
		}
		//fmt.Println("Routes added from parent2 =", routesAddedFromparent2)
		printSolution(offspring, instance)

		fmt.Println("visitedPatients: ", len(visitedPatients))
		for _, pID := range extractUnvisitedPatients(visitedPatients, instance) {
			fmt.Println("extractUnvisitedPatients(visitedPatients, instance) ", len(extractUnvisitedPatients(visitedPatients, instance))) 
			patientAdded := false
			
			for _, route := range offspring.Routes { 
				feasibleRoute, ok := route.canAddPatient(pID, instance)
				if ok {
					if len(offspring.Routes) < instance.NbrNurses {
						offspring.addRoute(feasibleRoute)
						visitedPatients = registerPatients(route, visitedPatients)
						patientAdded = true
					}
				}
			}
			
			// CORRECT FOR FIRST ITERATION, NOT AFTERWARDS. something must be wrong in the update of these
			fmt.Println("visitedPatients: ", len(visitedPatients))
			fmt.Println("extractUnvisitedPatients(visitedPatients, instance) ", len(extractUnvisitedPatients(visitedPatients, instance))) 

			if !patientAdded && (len(offspring.Routes) < instance.NbrNurses) {
				newRoute := initalizeOneRoute(instance)
				newRoute.visitPatient(instance.Patients[strconv.Itoa(pID)], instance)
				offspring.addRoute(newRoute)
				visitedPatients = registerPatients(newRoute, visitedPatients)
			}
		}

		parent1 = offspring
		iteration++
		fmt.Println("NUMBER OF OFFSPRING ROUTES ", len(offspring.Routes))
		fmt.Println("LENGTH OF TOTAL UNVISITED PATIENTS ", len(extractUnvisitedPatients(visitedPatients, instance)))
		fmt.Println("LENTH OF VISITED PATIENTS", len(visitedPatients))
		fmt.Println("\n\niteration =", iteration, "AND numPArents =", numParents)
	}
	parent1.calculateFitness(instance)

	if checkForDuplicates(parent1) {
		fmt.Println("Duplicates found")
	} else {
		fmt.Println("No duplicates found")
	}

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
				return true
			}
			allPatients = append(allPatients, patient.ID)
		}
	}

	return false
	
}