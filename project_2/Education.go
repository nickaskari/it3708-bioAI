package main

import (
	"math"
	"math/rand"
	"time"
)

/*
Performs education (local search) by doing mutations if there is improved objective value. Returns
educated individual. Purpose of this is to only find LOCAL optima.
*/
func hillClimbing(individual Individual, temp int, instance Instance) Individual {
	currentState := deepCopyIndividual(individual)

	for temp != 0 {
		invertedIndividual := randomInversionMutation(currentState, instance)
		invertedIndividual.fixAllRoutesAndCalculateFitness(instance)

		if invertedIndividual.Fitness <= currentState.Fitness {
			currentState = invertedIndividual
			swappedIndividual := randomSwapMutation(currentState, instance)
			swappedIndividual.fixAllRoutesAndCalculateFitness(instance)

			if swappedIndividual.Fitness <= currentState.Fitness {
				currentState = swappedIndividual

				interSwappedIndividual := randomInterRouteSwapMutation(currentState, instance)
				interSwappedIndividual.fixAllRoutesAndCalculateFitness(instance)

				if interSwappedIndividual.Fitness <= currentState.Fitness {
					currentState = interSwappedIndividual
				}
			}
		}
		currentState.fixAllRoutesAndCalculateFitness(instance)

		temp--
	}
	return currentState
}

// Simulated anhealing of individual. Returns optimized individual
func simulatedAnnealing(initialIndividual Individual, initialTemp int, coolingRate float64, instance Instance) Individual {

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	currentState := deepCopyIndividual(initialIndividual)
	currentTemp := float64(initialTemp)

	mutations := []func(Individual, Instance) Individual{
		randomInversionMutation,
		randomSwapMutation,
		randomInterRouteSwapMutation,
	}

	for currentTemp > 1 {
		mutation := mutations[random.Intn(len(mutations))]
		mutatedIndividual := mutation(currentState, instance)
		mutatedIndividual.fixAllRoutesAndCalculateFitness(instance)

		acceptMutation := shouldAcceptMutation(currentState.Fitness, mutatedIndividual.Fitness, currentTemp)
		if acceptMutation {
			currentState = mutatedIndividual
		}

		currentTemp *= coolingRate
	}
	currentState.fixAllRoutesAndCalculateFitness(instance)
	return currentState
}

func shouldAcceptMutation(currentFitness, newFitness, temperature float64) bool {
	if newFitness < currentFitness {
		return true
	}
	changeFitness := newFitness - currentFitness
	probability := math.Exp(-changeFitness / temperature)

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	return random.Float64() < probability
}

// Educates the Elite. Returns educated population
func educateTheElite(elitismPercentage float64, individuals []Individual, initialTemp int, coolingRate float64, instance Instance) []Individual {

	educatedIndividuals := deepCopyIndividuals(individuals)
	numToEducate := int(math.Floor(float64(len(individuals)) * elitismPercentage))
	for i := range len(individuals) {
		if (i + 1) > numToEducate {
			break
		}
		educatedIndividual := deepCopyIndividual(individuals[i])
		//individuals[i] = simulatedAnnealing(educatedIndividual, initialTemp,
		//	coolingRate, instance)

		individuals[i] = destroyRepairCluster(individuals[i], instance)
		individuals[i] = hillClimbing(educatedIndividual, 80, instance)
	}
	return educatedIndividuals
}
