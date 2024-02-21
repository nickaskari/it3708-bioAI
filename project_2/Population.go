package main


type Population struct {
	Individuals []Individual
}

// Initializes a new population
func createPopulation(instance Instance, populationSize int) Population {
	population := Population{make([]Individual, populationSize)}

	for i := 0; i < populationSize; i++ {
		population.Individuals[i] = createGreedyIndividual(instance)
	}

	return population
}

// Prints average fitnees of the population, best fitness and worst fitness
func printPopulationStats(population Population) {
	averageFitness := 0.0
	bestFitness := 0.0
	worstFitness := 0.0

	for _, individual := range population.Individuals {
		averageFitness += individual.Fitness

		if individual.Fitness > bestFitness {
			bestFitness = individual.Fitness
		}

		if individual.Fitness < worstFitness {
			worstFitness = individual.Fitness
		}
	}

	averageFitness = averageFitness / float64(len(population.Individuals))

	println("Average fitness: ", averageFitness)
	println("Best fitness: ", bestFitness)
	println("Worst fitness: ", worstFitness)
}
