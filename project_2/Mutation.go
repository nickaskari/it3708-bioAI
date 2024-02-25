package main

import (
	"fmt"
	"math/rand"
)


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
