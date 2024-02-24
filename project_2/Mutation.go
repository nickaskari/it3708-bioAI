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

func swapMutationRoute(route Route, instance Instance) Route {
	r := deepCopyRoute(route)

	for i := 0; i < len(r.Patients)-1; i++ {
		for j := i + 1; j < len(r.Patients); j++ {
			if !notViolatesTimeWindowConstraints(r, r.Patients[i], instance) && !notViolatesTimeWindowConstraints(r, r.Patients[j], instance) {
				r.Patients[i], r.Patients[j] = r.Patients[j], r.Patients[i]
				newRoute := createRouteFromPatientsVisited(r.Patients, instance)

				// Assuming there is logic here to calculate the fitness of the route
				oldFitness := calculateRouteFitness(route, instance)    // Placeholder function
				newFitness := calculateRouteFitness(newRoute, instance) // Placeholder function

				if newFitness < oldFitness {
					fmt.Println("Route got mutated")
					return newRoute
				} else {
					// Revert the swap if no improvement
					r.Patients[j], r.Patients[i] = r.Patients[i], r.Patients[j]
				}
			}
		}
	}

	return route
}

func swapMutationIndividual(individual Individual, instance Instance) Individual {
	for routeIndex, originalRoute := range individual.Routes {
		newRoute := swapMutationRoute(originalRoute, instance)
		if mutated {
			fmt.Println("Nurse", routeIndex+1, "route got mutated")
			individual.Routes[routeIndex] = newRoute

			individual.Fitness
		}
	}

	return individual
}

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
