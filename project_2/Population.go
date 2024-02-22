package main

import (
	"math"
	"fmt"
)

type Population struct {
	Individuals []Individual
	BestIndividual Individual
}

// Prints average fitnees of the population, best fitness and worst fitness
func (p Population) printPopulationStats() {
	averageFitness := 0.0
	bestFitness := 0.0
	worstFitness := 0.0

	for _, individual := range p.Individuals {
		averageFitness += individual.Fitness

		if individual.Fitness > bestFitness {
			bestFitness = individual.Fitness
		}

		if individual.Fitness < worstFitness {
			worstFitness = individual.Fitness
		}
	}

	averageFitness = averageFitness / float64(len(p.Individuals))

	fmt.Println("Average fitness: ", averageFitness)
	fmt.Println("Best fitness: ", bestFitness)
	fmt.Println("Worst fitness: ", worstFitness)
}

func (p Population) printBestIndividual(instance Instance) {
	printSolution(p.BestIndividual, instance)
}


// Initializes a new population
func createPopulation(instance Instance, populationSize int) Population {
	individuals := make([]Individual, populationSize)

	bestIndividual := Individual{math.Inf(1), make([]Route, 0)}

	for i := 0; i < populationSize; i++ {
		randomIndividual := createIndividual(instance)
		individuals[i] = randomIndividual

		if randomIndividual.Fitness < bestIndividual.Fitness {
			bestIndividual = randomIndividual
		}
	}

	return Population{Individuals: individuals, BestIndividual: bestIndividual}
}