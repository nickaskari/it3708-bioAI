package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

type Population struct {
	Individuals []Individual
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

	return Population{Individuals: individuals}
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
func printBestIndividual(i []Individual, instance Instance) {
	best := getBestIndividual(i)
	printSolution(best, instance)
	best.checkIndividualRoutes(instance, true)
}

// Performs tournamentselection for parent selection. Input is number of desired parents. Returns all chosen parents. (deterministic)
func (p Population) tournamentSelection(numParents int) []Individual {
	contestants := deepCopyIndividuals(p.Individuals)
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
	return deepCopyIndividual(bestIndividual)
}

// get best individual index
func getBestIndividualIndex(individuals []Individual) int {
	bestIndividual := math.Inf(1)
	index := 0
	for i, individual := range individuals {
		if individual.Fitness < bestIndividual {
			index = i
			bestIndividual = individual.Fitness
		}
	}
	return index
}
// Performs elitism for surivior selection. Returns the new population.
func (p Population) applyElitismWithPercentage(newPopulation []Individual, elitismPercentage float64) Population {
	numToPreserve := int(math.Floor(float64(p.size()) * elitismPercentage))

	// Sort the old population by fitness to find the fittest individuals, by making a copy. BEST TO WORST
	sortedOldIndividuals := deepCopyIndividuals(p.Individuals)

	sort.Slice(sortedOldIndividuals, func(i, j int) bool {
		return sortedOldIndividuals[i].Fitness < sortedOldIndividuals[j].Fitness // For minimization
	})

	// Sort the new population by fitness for WORST TO BEST (opposite of last one)
	sortedNewIndividuals := deepCopyIndividuals(newPopulation)
	sort.Slice(sortedNewIndividuals, func(i, j int) bool {
		return sortedNewIndividuals[i].Fitness > sortedNewIndividuals[j].Fitness
	})

	finalIndividuals := []Individual{}

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
		Individuals: finalIndividuals,
	}
}

// Returns the size of the population
func (p Population) size() int {
	return len(p.Individuals)
}

func deepCopyPopulation(original Population) Population {
	copy := original

	copy.Individuals = make([]Individual, len(original.Individuals))
	for i, individual := range original.Individuals {
		copy.Individuals[i] = deepCopyIndividual(individual)
	}

	return copy
}

// Does elitism with destruction of all the other individuals. Performs create individual again. Returns new population
func (p Population) applyGenecoideWithElitism(elitismPercentage float64, instance Instance) Population {
	numToPreserve := int(math.Floor(float64(p.size()) * elitismPercentage))

	// Sort the old population by fitness to find the fittest individuals, by making a copy. BEST TO WORST
	sortedOldIndividuals := deepCopyIndividuals(p.Individuals)

	sort.Slice(sortedOldIndividuals, func(i, j int) bool {
		return sortedOldIndividuals[i].Fitness < sortedOldIndividuals[j].Fitness // For minimization
	})

	finalIndividuals := []Individual{}

	for index := range p.size() {
		if index < numToPreserve {
			oldFitIndividual := deepCopyIndividual(sortedOldIndividuals[index])
			finalIndividuals = append(finalIndividuals, oldFitIndividual)
		} else {
			individual := createIndividual(instance)
			newIndividual := deepCopyIndividual(individual)
			finalIndividuals = append(finalIndividuals, newIndividual)
		}
	}

	return Population{
		Individuals: finalIndividuals,
	}
}

// Spreads a disease on all individual except elites. (Patients are scrambled) returns new population
func (p Population) spreadDisease(elitismPercentage float64, instance Instance) Population {
	//numToPreserve := int(math.Floor(float64(p.size()) * elitismPercentage))

	// Sort the old population by fitness to find the fittest individuals, by making a copy. BEST TO WORST
	sortedOldIndividuals := deepCopyIndividuals(p.Individuals)

	sort.Slice(sortedOldIndividuals, func(i, j int) bool {
		return sortedOldIndividuals[i].Fitness < sortedOldIndividuals[j].Fitness // For minimization
	})

	finalIndividuals := []Individual{}

	for index, elite := range sortedOldIndividuals {
		if index < len(sortedOldIndividuals) {
			oldFitIndividual := deepCopyIndividual(elite)
			finalIndividuals = append(finalIndividuals, oldFitIndividual)
		} else {
			newIndividual := deepCopyIndividual(elite)

			randomPatients := generateRandomPatientIDs(instance)
			newIndividual.removePatients(randomPatients, instance)
			//newIndividual.findBestRoutesForPatients(randomPatients, instance)
			newIndividual.distributePatientsOnRoutes(randomPatients, instance)
			newIndividual.fixAllRoutesAndCalculateFitness(instance)
			finalIndividuals = append(finalIndividuals, newIndividual)
		}
	}

	return Population{
		Individuals: finalIndividuals,
	}
}

// Ages population, returns new population
func (p Population) agePopulation() Population {
	aged := deepCopyPopulation(p)

	for i, ind := range aged.Individuals {
		ind.growOlder()
		aged.Individuals[i] = ind
	}

	return aged
}

// Perform age selection, returns new population
func ageSurvivorSelection(popSize int, newIndividuals []Individual) []Individual {
	newIndividuals = deepCopyIndividuals(newIndividuals)

	sort.Slice(newIndividuals, func(i, j int) bool {
		return newIndividuals[i].Age < newIndividuals[j].Age
	})

	newIndividuals = newIndividuals[:popSize]

	return newIndividuals
}


// Returns random migrants. Removes those from the population.
func (p *Population) selectRandomMigrants(numMigrants int) []Individual {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	migrants := make([]Individual, 0)

	newIndividuals := deepCopyIndividuals(p.Individuals)
	for len(migrants) < numMigrants {
		index := random.Intn(len(newIndividuals))

		migrants = append(migrants, newIndividuals[index])
		
		// Delete migrated individual from new individuals
		newIndividuals[index] = newIndividuals[len(newIndividuals) - 1]
		newIndividuals = newIndividuals[:len(newIndividuals) - 1]
	}

	p.Individuals = newIndividuals

	return migrants
}

// Insert new migrants into population. 
func (p *Population) insertNewMigrants(migrants []Individual) {
	newIndividuals := deepCopyIndividuals(p.Individuals)

	newIndividuals = append(newIndividuals, migrants...)

	p.Individuals = newIndividuals
}
