package main

import (
	"math/rand"
	"time"
)

// Swaps two patients in a random route of an individual. Returns individual.
func randomSwapMutation(individual Individual, instance Instance) Individual {
	mutated := deepCopyIndividual(individual)

	// Get route of length > 1
	randomRouteIndex := mutated.getRandomRoute()
	route := mutated.Routes[randomRouteIndex]

	for i := 0; i < len(route.Patients)-1; i++ {
		for j := i + 1; j < len(route.Patients); j++ {
			route.Patients[i], route.Patients[j] = route.Patients[j], route.Patients[i]

			// If both violates time window constraint
			if !notViolatesTimeWindowConstraints(route, route.Patients[i], instance) && !notViolatesTimeWindowConstraints(route, route.Patients[j], instance) {

				newRoute := createRouteFromPatientsVisited(route.Patients, instance)
				if newRoute.CurrentTime <= float64(instance.Depot.ReturnTime) {
					mutated.Routes[randomRouteIndex] = newRoute
					return mutated
				}
			}
			route.Patients[i], route.Patients[j] = route.Patients[j], route.Patients[i]
		}
	}

	// No swap occured? Just swap at random. Only do swap if return time is not violated

	i, j := getTwoSwapIndexes(route)
	route.Patients[i], route.Patients[j] = route.Patients[j], route.Patients[i]
	newRoute := createRouteFromPatientsVisited(route.Patients, instance)
	if newRoute.CurrentTime <= float64(instance.Depot.ReturnTime) {
		mutated.Routes[randomRouteIndex] = newRoute
		return mutated
	}

	return mutated
}

func getTwoSwapIndexes(route Route) (int, int) {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	for {
		i := random.Intn(len(route.Patients))
		j := random.Intn(len(route.Patients))
		if i != j {
			return i, j
		}
	}
}

// Performs inversion mutations on a random route in individual. Returns individual
func randomInversionMutation(individual Individual, instance Instance) Individual {
	mutated := deepCopyIndividual(individual)

	// Get route of length > 1
	randomRouteIndex := mutated.getRandomRoute()
	r := mutated.Routes[randomRouteIndex]

	if len(r.Patients) < 2 {
		return mutated
	}

	start, end := getTwoSwapIndexes(r)
	if start > end {
		start, end = end, start
	}

	for i, j := start, end; i < j; i, j = i+1, j-1 {
		r.Patients[i], r.Patients[j] = r.Patients[j], r.Patients[i]
	}

	
	newRoute := createRouteFromPatientsVisited(r.Patients, instance)
	if newRoute.CurrentTime <= float64(instance.Depot.ReturnTime) {
		mutated.Routes[randomRouteIndex] = newRoute
		return mutated
	}
	mutated.Routes[randomRouteIndex] = newRoute
	return mutated
}

// Performs random inter route mutation. Returns mutated individual (if mutated)..

func randomInterRouteSwapMutation(individual Individual, instance Instance) Individual {
	mutated := deepCopyIndividual(individual)

	i, j := getTwoRouteIndexes(mutated)

	route1, route2 := deepCopyRoute(mutated.Routes[i]), deepCopyRoute(mutated.Routes[j])

	patient1, patient2 := route1.getRandomPatient(), route2.getRandomPatient()

	route1, ok1 := route1.performPatientSwap(patient1, patient2, instance)
	route2, ok2 := route2.performPatientSwap(patient2, patient1, instance)

	if ok1 && ok2 {
		mutated.Routes[i], mutated.Routes[j] = route1, route2
		return mutated
	}

	return mutated
}

/*
Returns two random route indexes based on an individual. Returns two route indexes that are not the same.
The routes are not empty.
*/
func getTwoRouteIndexes(individual Individual) (int, int) {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	for {
		i := random.Intn(len(individual.Routes))
		j := random.Intn(len(individual.Routes))
		if i != j {
			if len(individual.Routes[i].Patients) > 0 && len(individual.Routes[j].Patients) > 0 {
				return i, j
			}
		}
	}
}

// Random inter move. Returns mutated individual (if mutated).. NOT FINISHED

func randomInterRouteMoveMutation(individual Individual, instance Instance) Individual {
	mutated := deepCopyIndividual(individual)

	i, j := getTwoRouteIndexes(mutated)

	route1, route2 := mutated.Routes[i], mutated.Routes[j]

	patient1, patient2 := route1.getRandomPatient(), route2.getRandomPatient()

	route1, ok1 := route1.performPatientSwap(patient1, patient2, instance)
	route2, ok2 := route2.performPatientSwap(patient2, patient1, instance)

	if ok1 && ok2 {
		mutated.Routes[i], mutated.Routes[j] = route1, route2
		return mutated
	}

	return mutated
}
