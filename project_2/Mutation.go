package main

import (
	"fmt"
	"math/rand"
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





func swapMutationRoute(originalRoute Route, instance Instance) (Route, bool) {
    r := deepCopyRoute(originalRoute)
    originalFitness := calculateRouteFitness(originalRoute, instance)

    for i := 0; i < len(r.Patients)-1; i++ {
        for j := i + 1; j < len(r.Patients); j++ {
            r.Patients[i], r.Patients[j] = r.Patients[j], r.Patients[i]
         
            if !notViolatesTimeWindowConstraints(r, r.Patients[i], instance) && !notViolatesTimeWindowConstraints(r, r.Patients[j], instance) {

                newRoute := createRouteFromPatientsVisited(r.Patients, instance)
                newFitness := calculateRouteFitness(newRoute, instance)

                if newFitness < originalFitness {
                    return newRoute, true
                }
            }
            r.Patients[i], r.Patients[j] = r.Patients[j], r.Patients[i]
        }
    }
    return originalRoute, false 
}


func swapMutationIndividual(individual Individual, instance Instance) Individual {
    mutatedIndividual := deepCopyIndividual(individual) 
    anyMutationOccurred := false

    for routeIndex, originalRoute := range mutatedIndividual.Routes {
        mutatedRoute, mutated := swapMutationRoute(originalRoute, instance)
        if mutated {
            mutatedIndividual.Routes[routeIndex] = mutatedRoute
            anyMutationOccurred = true 
			fmt.Println("Nurse", routeIndex+1, "got mutated")
        }
    }

    if anyMutationOccurred {
        mutatedIndividual.calculateFitness(instance)
    } else {
        fmt.Println("DID NOT MUTATE")
    }
	oldFitness := individual.Fitness
	newFitness := mutatedIndividual.Fitness
	fmt.Println("DIFFERENCE =", newFitness-oldFitness)
    return mutatedIndividual
}





// Selects a random route and inverts the order of patients between two random points in the route. Accepts regardless of fitness.
func inversionMutation(originalRoute Route, instance Instance) (Route, bool) {
    if len(originalRoute.Patients) < 2 {
        return originalRoute, false
    }

    r := deepCopyRoute(originalRoute)
    originalFitness := calculateRouteFitness(originalRoute, instance)

    start, end := rand.Intn(len(r.Patients)-1), rand.Intn(len(r.Patients)-1)
    if start > end {
        start, end = end, start
    }

    for i, j := start, end; i < j; i, j = i+1, j-1 {
        r.Patients[i], r.Patients[j] = r.Patients[j], r.Patients[i]
    }

    newRoute := createRouteFromPatientsVisited(r.Patients, instance)
    newFitness := calculateRouteFitness(newRoute, instance)

    if newFitness < originalFitness {
		return newRoute, true
    }

    return originalRoute, false
}

func inversionMutationIndividual(individual Individual, instance Instance) Individual {
    mutatedIndividual := deepCopyIndividual(individual) // Start with a deep copy to apply mutations progressively
    anyMutationOccurred := false

    for routeIndex, originalRoute := range mutatedIndividual.Routes {
        mutatedRoute, mutated := inversionMutation(originalRoute, instance)
        if mutated {
            mutatedIndividual.Routes[routeIndex] = mutatedRoute
            anyMutationOccurred = true 
			fmt.Println("Nurse", routeIndex+1, "got mutated")
        }
    }

    if anyMutationOccurred {
        mutatedIndividual.calculateFitness(instance)
    } else {
        fmt.Println("NO INVERSION MUTATION OCCURRED")
    }
	oldFitness := individual.Fitness
	newFitness := mutatedIndividual.Fitness
	fmt.Println("DIFFERENCE =", newFitness-oldFitness)
    return mutatedIndividual
}
