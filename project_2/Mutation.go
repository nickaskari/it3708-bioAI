package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Performs swap mutation. Return mutated individual. Mutates in only one route.
func swapMutation(individual Individual, instance Instance) Individual {
	//newRoutes := make(map[int]Route)

	for routeIndex, originalRoute := range individual.Routes {
		// Traverse the route
		r := deepCopyRoute(originalRoute)
		for i := 0; i < len(r.Patients)-1; i++ {
			for j := i + 1; j < len(r.Patients); j++ {
				if !notViolatesTimeWindowConstraints(r, r.Patients[i], instance) && !notViolatesTimeWindowConstraints(r, r.Patients[j], instance) {
					r.Patients[i], r.Patients[j] = r.Patients[j], r.Patients[i]
					newRoute := createRouteFromPatientsVisited(r.Patients, instance)

					// Problem --> individual changes
					newIndividual := createAlteredIndivual(individual, routeIndex, newRoute, instance)

					oldFitness := individual.Fitness
					newFitness := newIndividual.Fitness

					//fmt.Println("DIFFERENCE =", newFitness-oldFitness)
					if newFitness < oldFitness {
						fmt.Println("Nurse", routeIndex+1, "got mutated")
						fmt.Println("DIFFERENCE =", newFitness-oldFitness)
						return newIndividual
					}
				}
			}
		}
	}
	fmt.Println("DID NOT MUTATE")
	return individual
}

func createAlteredIndivual(individual Individual, routeIndex int, route Route, instance Instance) Individual {
	newIndividual := deepCopyIndividual(individual)

	newIndividual.Routes[routeIndex] = route
	newIndividual.calculateFitness(instance)

	return newIndividual
}

/*

func swapRoutes(individual *Individual, index int, newRoute Route) {
	// delete the old route and replace it with the new one
	individual.Routes = append(individual.Routes[:index], individual.Routes[index+1:]...)
	individual.Routes = append(individual.Routes, newRoute)

	//individual.Routes[index] = newRoute
	//return individual
}

// Swap Mutation that tries to mitigate time window violations. Only accepts if it improves fitness.
func swapMutation(individual *Individual, instance Instance) {

	// iterate through routes, when a time window violation is found for two patients in the same route,
	// swap them if it improves the fitness for the individual

	outerLoop:
	for i := range individual.Routes {
		route := &individual.Routes[i]
		patientsCopy := make([]Patient, len(route.Patients))
		copy(patientsCopy, route.Patients)

		fmt.Println("first loop")

		for j := 0; j < len(route.Patients)-1; j++ {
			fmt.Println("second loop")
			for k := j + 1; k < len(route.Patients); k++ {
				fmt.Println("third loop")
				// Only consider swapping if one of the patients violates time window constraints
				if violatesTimeWindowConstraints(*route, route.Patients[j], instance) || violatesTimeWindowConstraints(*route, route.Patients[k], instance) {
					patientsCopy[j], patientsCopy[k] = patientsCopy[k], patientsCopy[j]

					// Create a new route based on this swapped configuration to assess the impact
					newRoute := createRouteFromPatientsVisited(patientsCopy, instance)

					// Temporarily update the route with the new configuration to calculate fitness
					originalRoute := *route // Remember the original route
					*route = newRoute

					oldFitness := individual.Fitness
					fmt.Println("old fitness: ", oldFitness)
					swapRoutes(individual, i, newRoute)

					// Recalculate the fitness with the updated route configuration

					newFitness := individual.Fitness
					fmt.Println("new fitness: ", newFitness)

					if newFitness <= oldFitness {
						*route = originalRoute
						individual.Fitness = oldFitness
						patientsCopy[j], patientsCopy[k] = patientsCopy[k], patientsCopy[j]
					} else {
						// If there's an improvement, update the original route to reflect the new patient order
						copy(route.Patients, patientsCopy)
						break outerLoop
					}
				}
			}
		}
	}
}

*/
// Selects a random route and inverts the order of patients between two random points in the route. Accepts regardless of fitness.
func inversionMutation(individual *Individual, instance Instance) {

	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)

	// Randomly select a route
	routeIndex := rand.Intn(len(individual.Routes))
	route := &individual.Routes[routeIndex]

	// Selected route need more than two patients for inversion to make sense
	if len(route.Patients) < 3 {
		return
	}

	// Randomly select two distinct points within the route for inversion
	point1 := rand.Intn(len(route.Patients) - 1)
	point2 := rand.Intn(len(route.Patients)-point1-1) + point1 + 1 // Ensure point2 is after point1

	// Invert the order of patients between point1 and point2
	for i, j := point1, point2; i < j; i, j = i+1, j-1 {
		route.Patients[i], route.Patients[j] = route.Patients[j], route.Patients[i]
	}

	// Can edit this to only accept if it improves fitness
}
