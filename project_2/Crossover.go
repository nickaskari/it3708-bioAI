package main

import (
	"math/rand"
	"slices"
	"time"
)

//"math"
//"math/rand"
//"slices"
//"time"

/*
	TODO
	Implement a crossover function + mutation
*/

/*
Performs crossover with edge recombination and a possible swap mutation (performed randomly with mutationRate).
Returns one child. MAKE SURE UNIQUE PATIOENTS IN CHILD
*/
func crossover(parent1 Individual, parent2 Individual, instance Instance, mutationRate float64) Individual {
	matchedRoutes := make(map[int][]int, 0)
	childRoutes := make([]Route, 0)

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	for i, r1 := range parent1.Routes {

		scoreTable := make([]int, len(parent2.Routes))
		for j, r2 := range parent2.Routes {
			score := calculateSimularityScore(r1, r2)
			scoreTable[j] = score
		}

		_, rankedBySimularity := sortWithReflection(scoreTable)
		matchedRoutes[i] = rankedBySimularity
	}

	usedRoutesFromParent2 := make([]int, 0)
	for parent1RouteIndex, preferedRoutesIndexesFromParent2 := range matchedRoutes {
		for _, i := range preferedRoutesIndexesFromParent2 {
			if !(slices.Contains(usedRoutesFromParent2, i)) {
				r1 := parent1.Routes[parent1RouteIndex]
				r2 := parent2.Routes[i]
				child := edgeRecombination(r1, r2, instance)
				childRoutes = append(childRoutes, child)
				usedRoutesFromParent2 = append(usedRoutesFromParent2, i)
			}
		}
	}

	child := Individual{
		Fitness: 0.0,
		Age:     0,
		Routes:  childRoutes,
	}

	child.calculateFitness(instance)

	if random.Float64() < mutationRate {
		child = inversionMutationIndividual(child, instance)
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
