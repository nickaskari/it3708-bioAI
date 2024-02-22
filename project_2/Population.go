package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

type Population struct {
	Individuals    []Individual
	BestIndividual Individual
	Size           int
}

// Initializes a new population
func initPopulation(instance Instance, populationSize int) Population {
	individuals := make([]Individual, populationSize)

	bestIndividual := Individual{Fitness: math.Inf(1), Age: 0, Routes: make([]Route, 0)}

	for i := 0; i < populationSize; i++ {
		randomIndividual := createIndividual(instance)
		individuals[i] = randomIndividual

		if randomIndividual.Fitness < bestIndividual.Fitness {
			bestIndividual = randomIndividual
		}
	}

	bestIndividual.writeIndividualToJson()

	return Population{Individuals: individuals, BestIndividual: bestIndividual, Size: populationSize}
}

// Prints average fitnees of the population, best fitness and worst fitness
func (p Population) printPopulationStats() {
	averageFitness := 0.0
	bestFitness := math.Inf(1)
	worstFitness := math.Inf(-1)
	ages := make([]int, 0)

	for _, individual := range p.Individuals {
		averageFitness += individual.Fitness

		if individual.Fitness < bestFitness {
			bestFitness = individual.Fitness
		}

		if individual.Fitness > worstFitness {
			worstFitness = individual.Fitness
		}

		ages = append(ages, individual.Age)
	}

	maxAge := 0
	for _, age := range ages {
		if age > maxAge {
			maxAge = age
		}
	}

	minAge := 0
	for _, age := range ages {
		if age < minAge {
			minAge = age
		}
	}

	averageFitness = averageFitness / float64(len(p.Individuals))

	fmt.Println("Average fitness:", averageFitness)
	fmt.Println("Best individual:", bestFitness)
	fmt.Println("Worst individual:", worstFitness)
	fmt.Println("Oldest individual:", maxAge)
	fmt.Println("Youngest individual:", minAge)
	fmt.Println("Population count:", len(p.Individuals))
	printDivider(150, "-")
}

// Prints the best individual in a pretty format
func (p Population) printBestIndividual(instance Instance) {
	printSolution(p.BestIndividual, instance)
}

// Performs tournamentselection for parent selection. Returns all chosen parents. (deterministic)
func (p Population) tournamentSelection() []Individual {
	contestants := p.Individuals
	winners := make([]Individual, 0)
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for len(winners) != p.Size {
		if len(contestants) > 1 {
			size := 2 + r.Intn(p.Size-1)
			match := chooseRandomUnique[Individual](contestants, size)
			winner := getBestIndividual(match)
		}
	}
}

func getBestIndividual(individuals []Individual) Individual {
	bestIndividual := Individual{Fitness: math.Inf(1), Age: 0, Routes: make([]Route, 0)}

	for _, individual := range individuals {
		if individual.Fitness < bestIndividual.Fitness {
			bestIndividual = individual
		}
	}
	return bestIndividual
}

// Performs elitism for surivior selection. Returns all surviving individuals
func (p *Population) applyElitismWithPercentage(newwPopulation []Individual, elitismPercentage float64) {

	numToPreserve := int(float64(len(p.Individuals))*elitismPercentage/100.0 + 0.5) // Rounded

	// Sort the current population by fitness to find the fittest individuals
	// Make a copy of the slice to avoid modifying the original population order
	sortedIndividuals := make([]Individual, len(p.Individuals))
	copy(sortedIndividuals, p.Individuals)
	sort.Slice(sortedIndividuals, func(i, j int) bool {
		return sortedIndividuals[i].Fitness < sortedIndividuals[j].Fitness // For minimization
	})

	// Select the fittest individuals based on the calculated number to preserve
	fittestIndividuals := sortedIndividuals[:numToPreserve]

	// Ensure the fittest individuals are included in the new generation
	for _, fittest := range fittestIndividuals {
		// Check if this fittest individual is already in the new generation
		found := false
		for _, individual := range newIndividuals {
			if individual.Fitness == fittest.Fitness {
				found = true
				break
			}
		}

		// If not found, replace the least fit individual in the new generation with this fittest individual
		if !found {
			// Find the least fit individual in the new generation
			worstFitnessIndex := 0
			worstFitness := newIndividuals[0].Fitness
			for i, individual := range newIndividuals {
				if individual.Fitness > worstFitness {
					worstFitness = individual.Fitness
					worstFitnessIndex = i
				}
			}

			newIndividuals[worstFitnessIndex] = fittest
		}
	}

	// Update the population with the new generation, now including the preserved fittest individuals
	p.Individuals = newIndividuals

	// Optionally, update the BestIndividual if needed
	// This assumes BestIndividual should still reflect the overall best found so far
	if len(fittestIndividuals) > 0 {
		p.BestIndividual = fittestIndividuals[0] // The first one is the best due to sorting
	}

}
