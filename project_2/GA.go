package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func GA(populationSize int, gMax int, numParents int, temp int,
	crossoverRate float64, mutationRate float64, elitismPercentage float64, coolingRate float64,
	annealingRate float64, instance Instance) {
	
	// initialize an emtpy array
	bestFitnesses := []float64{}

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	generation := 0
	var newIndividuals []Individual
	stuck := 0
	lastFitness := math.Inf(1)

	fmt.Println("Initalzing population..")
	population := initPopulation(instance, populationSize)
	printBestIndividual(population.Individuals, instance)

	for generation < gMax {
		newIndividuals = []Individual{}

		parents := population.tournamentSelection(numParents)

		for population.size() > len(newIndividuals) {
			source := rand.NewSource(time.Now().UnixNano())
			r := rand.New(source)

			i, j := getTwoRandomParents(parents)
			//i, j := getTwoSimilarParents(parents)
			parent1, parent2 := parents[i], parents[j] // choose based on similarity
			if r.Float64() < crossoverRate {
				child1, child2 := destroyRepairCrossover(parent1, parent2, instance)
				child1.calculateFitness(instance)
				child2.calculateFitness(instance)

				if r.Float64() < mutationRate {
					var mutated1 Individual
					var mutated2 Individual

					if annealingRate > random.Float64() {
						mutated1 = simulatedAnnealing(child1, temp, coolingRate, instance)
						mutated1.calculateFitness(instance)
					} else {
						mutated1 = hillClimbing(child1, temp, instance)
						mutated1.calculateFitness(instance)
					}

					if annealingRate > random.Float64() {
						mutated2 = simulatedAnnealing(child2, temp, coolingRate, instance)
						mutated2.calculateFitness(instance)
					} else {
						mutated2 = hillClimbing(child2, temp, instance)
						mutated2.calculateFitness(instance)
					}

					newIndividuals = addToPopulation(mutated1, population.size(), newIndividuals)
					newIndividuals = addToPopulation(mutated2, population.size(), newIndividuals)
				} else {
					newIndividuals = addToPopulation(child1, population.size(), newIndividuals)
					newIndividuals = addToPopulation(child2, population.size(), newIndividuals)
				}
			}
			newIndividuals = addToPopulation(parent1, population.size(), newIndividuals)
			newIndividuals = addToPopulation(parent2, population.size(), newIndividuals)
		}

		// Survivor selection -- ELITISM

		// grows populationx
		newPopulation := deepCopyPopulation(population.applyElitismWithPercentage(newIndividuals, elitismPercentage))
		population = deepCopyPopulation(newPopulation)

		fmt.Println("GENEREATION", generation+1)
		population.printPopulationStats()
		printBestIndividual(population.Individuals, instance)

		bestFitness := getBestIndividual(population.Individuals).Fitness
		bestFitnesses = append(bestFitnesses, bestFitness)

		if bestFitness == lastFitness {
			stuck++
		} else {
			lastFitness = bestFitness
		}

		if stuck > 5 {
			var newPopulation Population
			if 0.8 > random.Float64() {
				fmt.Println("\nPERFORM GENOCIDE AND REBUILD POPULATION..\n")
				newPopulation = deepCopyPopulation(population.applyGenecoideWithElitism(elitismPercentage, instance))
			} else {
				fmt.Println("\nSPREAD DISEASE..\n")
				newPopulation = deepCopyPopulation(population.spreadDisease(elitismPercentage, instance))
			}
			population = newPopulation
			stuck = 0
		}
		generation++
	}

	getBestIndividual(population.Individuals).writeIndividualToJson()
	getBestIndividual(population.Individuals).writeIndividualToVismaFormat()
	writeBestFitnessesToJSON(bestFitnesses)

}

// Get Two random indexes that are not the same
func getTwoRandomParents(parents []Individual) (int, int) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for {
		// indexes
		i := r.Intn(len(parents))
		j := r.Intn(len(parents))
		if i != j {
			return i, j
		}
	}
}

// Appends individual to array only if size resitriction is good. Returns array.
func addToPopulation(toAdd Individual, size int, individuals []Individual) []Individual {
	if len(individuals) < size {
		return append(individuals, toAdd)
	}
	return individuals
}
