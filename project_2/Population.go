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
			size := 2 + r.Intn(p.size()-1)
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

// Performs elitism for surivior selection. Returns the new population.
func (p Population) applyElitismWithPercentage(newPopulation []Individual, elitismPercentage float64) Population {
	numToPreserve := int(math.Floor(float64(p.size()) * elitismPercentage))

	// Sort the old population by fitness to find the fittest individuals, by making a copy. BEST TO WORST
	sortedOldIndividuals := make([]Individual, p.size())
	copy(sortedOldIndividuals, p.Individuals)
	sort.Slice(sortedOldIndividuals, func(i, j int) bool {
		return sortedOldIndividuals[i].Fitness < sortedOldIndividuals[j].Fitness // For minimization
	})

	// Sort the new population by fitness for WORST TO BEST (opposite of last one)
	sortedNewIndividuals := make([]Individual, len(newPopulation))
	copy(sortedNewIndividuals, newPopulation)
	sort.Slice(sortedNewIndividuals, func(i, j int) bool {
		return sortedNewIndividuals[i].Fitness > sortedNewIndividuals[j].Fitness
	})

	finalIndividuals := []Individual{}

	fmt.Println("NUM TO PRESERVE:", numToPreserve)
	for index, individual := range sortedNewIndividuals {
		if index < numToPreserve {
			oldFitIndividual := deepCopyIndividual(sortedOldIndividuals[index])
			finalIndividuals = append(finalIndividuals, oldFitIndividual)
			//fmt.Println("\n\n\n\n\n\nIIMMMM HERE \n\n\n\n")
		} else {
			newIndividual := deepCopyIndividual(individual)
			finalIndividuals = append(finalIndividuals, newIndividual)
			//fmt.Println("\n\n\n\n\n\nNOOOOO IM GAAAA \n\n\n\n")
		}
	}

	return Population{
		Individuals:    finalIndividuals,
		BestIndividual: getBestIndividual(finalIndividuals),
	}

	/*

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
				// Find the least fit individual in the new generation.
				// Does not make sense
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
	*/
}

// Returns the size of the population
func (p Population) size() int {
	return len(p.Individuals)
}

// gets two random parents that are not the same
