package main

import (
	"math"
	"math/rand"
	"slices"
	"time"
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
		child = swapMutation(child, instance)
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