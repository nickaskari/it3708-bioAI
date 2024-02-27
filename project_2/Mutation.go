package main

import (
	"fmt"
	"math/rand"
)

// Swaps two patients in a route. Accepts if fitness is improved.
func swapMutationRoute(originalRoute Route, instance Instance) (Route, bool) {
    mutatedRoute := deepCopyRoute(originalRoute)
    originalFitness := calculateRouteFitness(originalRoute, instance)

    for i := 0; i < len(mutatedRoute.Patients)-1; i++ {
        for j := i + 1; j < len(mutatedRoute.Patients); j++ {
            mutatedRoute.Patients[i], mutatedRoute.Patients[j] = mutatedRoute.Patients[j], mutatedRoute.Patients[i]
            
            // if both violates
            if !notViolatesTimeWindowConstraints(mutatedRoute, mutatedRoute.Patients[i], instance) && !notViolatesTimeWindowConstraints(mutatedRoute, mutatedRoute.Patients[j], instance) {

                newRoute := createRouteFromPatientsVisited(mutatedRoute.Patients, instance)
                newFitness := calculateRouteFitness(newRoute, instance)

                if newFitness < originalFitness {
                    return newRoute, true
                }
            }
            mutatedRoute.Patients[i], mutatedRoute.Patients[j] = mutatedRoute.Patients[j], mutatedRoute.Patients[i]
        }
    }
    return originalRoute, false 
}

// performs one swap mutation on every route if it imrpoves fitness
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





// Selects a random route and inverts the order of patients between two random points in the route. Accepts if fitness is improved.
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

// performs inversion mutation on every route if it improves fitness
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

// make inter-route swap mutation. Use canAddPatientEnforced


// Local seacrch. Takes in an individual,, mutates it x number of times, using several muations functions, until 
    // no improvement is made, aka reaced local minima.


// performs patient swaps between routes
func interRouteSwapMutation(individual Individual, instance Instance) (Individual) {
    mutatedIndividual := deepCopyIndividual(individual)
    //originalFitness := individual.Fitness
    anyMutationOccurred := false

    randomRoute1, randomRoute2 := getTwoRandomRoutes(mutatedIndividual)

    for p_i := 0; p_i < len(randomRoute1.Patients); p_i++ {
        for p_j := 0; p_j < len(randomRoute2.Patients); p_j++ {
            randomRoute1.Patients[p_i], randomRoute2.Patients[p_j] = randomRoute2.Patients[p_j], randomRoute1.Patients[p_i]

            if !notViolatesTimeWindowConstraints(randomRoute1, randomRoute1.Patients[p_i], instance) && !notViolatesTimeWindowConstraints(randomRoute2, randomRoute2.Patients[p_j], instance) {
                
                _, canAdd1 := randomRoute1.canAddPatientEnforced(randomRoute1.Patients[p_j].ID, instance)
                _, canAdd2 := randomRoute2.canAddPatientEnforced(randomRoute2.Patients[p_i].ID, instance)
                
                if canAdd1 && canAdd2 {
                    anyMutationOccurred = true
                    fmt.Println("MUTATION OCCURRED -- ", "Nurse ", p_i + 1, " and Nurse ", p_j + 1, "swapped", randomRoute1.Patients[p_i].ID, "and", randomRoute2.Patients[p_j].ID)
                    return mutatedIndividual
                }
            }
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




// simulated annealing

/*

func simulatedAnnealing(individual Individual, probabilty float64) (individual){
    mutatedIndividual := deepCopyIndividual(individual)
    originalFitness := individual.Fitness

    Temperature := 1000

    for temperature > 0 {
       // selct mutation function at random
    random_int = rand.Intn(3)
    if random_int == 0 {
        mutatedIndividual := interRouteSwapMutation(individual, instance)
    } else if random_int == 1 {
        mutatedIndividual := swapMutationIndividual(individual, instance)
    } else {
        mutatedIndividual := inversionMutationIndividual(individual, instance)
    }

    newFitness := mutatedIndividual.Fitness

    if exp(- temperature * (newFitness - originalFitness)) > probability {
        return mutatedIndividual
    } else {
        return individual
    }

    temperature --
    }
}

*/

// returns two random non-empty routes from individual
func getTwoRandomRoutes(individual Individual) (route1 Route, route2 Route) {
    source := rand.NewSource(0)
    random := rand.New(source)

    randomRoute1 := random.Intn(len(individual.Routes)) 
    randomRoute2 := random.Intn(len(individual.Routes))
    
    if len(individual.Routes[randomRoute1].Patients) == 0 || len(individual.Routes[randomRoute2].Patients) == 0 {
        return getTwoRandomRoutes(individual)
    }

    return individual.Routes[randomRoute1], individual.Routes[randomRoute2]
}