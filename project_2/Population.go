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

	return Population{Individuals: individuals, BestIndividual: bestIndividual}
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
	p.BestIndividual.checkIndividualRoutes(instance, true)
}

// Performs tournamentselection for parent selection. Input is number of desired parents. Returns all chosen parents. (deterministic)
func (p Population) tournamentSelection(numParents int) []Individual {
	contestants := p.Individuals
	winners := make([]Individual, 0)
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for len(winners) != numParents {
		if len(contestants) > 1 {
			size := 2 + r.Intn(p.size() -1)
			match := chooseRandomUnique[Individual](contestants, size)
			winner := getBestIndividual(match)
			winner.removeIndividualFrom(match)
			winners = append(winners, winner)
		} else {
			winner := contestants[0]
			winners = append(winners, winner)
		}
	}
	return winners
}

// A helper function for tournamentSelection(). Returns best individual from list of individuals.
func getBestIndividual(individuals []Individual) Individual {
	bestIndividual := Individual{Fitness: math.Inf(1), Age: 0, Routes: make([]Route, 0)}

	for _, individual := range individuals {
		if individual.Fitness < bestIndividual.Fitness {
			bestIndividual = individual
		}
	}
	return bestIndividual
}

// Performs elitism for surivior selection. Returns the new population
func (p *Population) applyElitismWithPercentage(newPopulation []Individual, elitismPercentage float64) ([]Individual, Individual) {
	numToPreserve := int(float64(len(p.Individuals))*elitismPercentage/100.0 + 0.5) // Percentage to absolute number

	// Sort the old population by fitness to find the fittest individuals, by making a copy
	sortedIndividuals := make([]Individual, len(p.Individuals))
	copy(sortedIndividuals, p.Individuals)
	sort.Slice(sortedIndividuals, func(i, j int) bool {
		return sortedIndividuals[i].Fitness < sortedIndividuals[j].Fitness // For minimization
	})

	// Select the n fittest individuals based on the elitism percentage
	fittestIndividuals := sortedIndividuals[:numToPreserve]

	for _, fittest := range fittestIndividuals {
		// Check if this fittest individual is already in the new generation
		found := false
		for _, individual := range newPopulation {
			if individual.Fitness == fittest.Fitness {
				found = true
				// If found, break the loop
				break
			}
		}

		// If not found, replace the least fit individual in the new generation with this fittest individual from the old generation
		if !found {
			// Find the least fit individual in the new generation
			worstFitnessIndex := -1
			worstFitness := -1.0
			for i, individual := range newPopulation {
				if worstFitnessIndex == -1 || individual.Fitness > worstFitness {
					worstFitness = individual.Fitness
					worstFitnessIndex = i
				}
			}

			if worstFitnessIndex != -1 {
				newPopulation[worstFitnessIndex] = fittest
			}
		}
	}
	
	bestIndividual := getBestIndividual(newPopulation)
	return newPopulation, bestIndividual
}

// Returns the size of the population
func (p Population) size() int {
	return len(p.Individuals)
}

// gets two random parents that are not the same
