package main

import (
	"math"
	"math/rand"
	"slices"
	"time"
	"fmt"
)

/*
	TODO
	Implement a crossover function + mutation
*/

/*
Performs crossover with edge recombination and a possible swap mutation (performed randomly with mutationRate).
Returns one child.
*/
func crossover(parent1 Individual, parent2 Individual, instance Instance, mutationRate float64) Individual {
	matchedRoutes := make([]int, 0)
	childRoutes := make([]Route, 0)

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	for i, r1 := range parent1.Routes {
		i += len(parent2.Routes)
		mostSimularScore := int(math.Inf(-1))
		var mostSimularRoute Route
		var mostSimularRouteIndex int
		for j, r2 := range parent2.Routes {
			score := calculateSimularityScore(r1, r2)
			if score > mostSimularScore {
				mostSimularScore = score
				mostSimularRoute = r2
				mostSimularRouteIndex = j
			}
		}

		if !(slices.Contains(matchedRoutes, mostSimularRouteIndex)) {
			matchedRoutes = append(matchedRoutes, mostSimularRouteIndex)
			child := edgeRecombination(r1, mostSimularRoute, instance)
			childRoutes = append(childRoutes, child)
		}
	}

	child := Individual{
		Fitness: 0.0,
		Age: 0,
		Routes: childRoutes,
	}

	child.calculateFitness(instance)

	if random.Float64() < mutationRate {
		swapMutation(&child, instance)
	}
	return child
}

// Helper function for crossover. Checks simularity score between two routes. I.E. how many patients in common
func calculateSimularityScore(r1 Route, r2 Route) int {
	simularityScore := 0
	for _, p1 := range r1.Patients {
		for _, p2 := range r2.Patients {
			if p1.ID == p2.ID {
				simularityScore++
			}
		}
	}
	return simularityScore
}

func swapRoutes(individual *Individual, index int, newRoute Route) *Individual {
	individual.Routes[index] = newRoute
	return individual
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

		for j := 0; j < len(route.Patients)-1; j++ {
			for k := j + 1; k < len(route.Patients); k++ {
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
					newIndividual := swapRoutes(individual, i, newRoute)

					// Recalculate the fitness with the updated route configuration
					
					newFitness := newIndividual.Fitness
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


