package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func GA(populationSize int, gMax int, numParents int, temp int,
	crossoverRate float64, mutationRate float64, elitismPercentage float64, coolingRate float64,
	instance Instance) {

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
		newIndividuals = make([]Individual, 0)

		parents := population.tournamentSelection(numParents)

		// This condition is wrong? population size gets altered..
		for population.size() > len(newIndividuals) {
			source := rand.NewSource(time.Now().UnixNano())
			r := rand.New(source)

			i, j := getTwoRandomParents(parents)
			parent1, parent2 := parents[i], parents[j]
			if r.Float64() < crossoverRate {
				child1, child2 := destroyRepairCrossover(parent1, parent2, instance)

				if r.Float64() < mutationRate {
					var mutated1 Individual
					var mutated2 Individual

					if 0.8 > random.Float64() {
						mutated1 = simulatedAnnealing(child1, temp, coolingRate, instance)
					} else {
						mutated1 = hillClimbing(child1, temp, instance)
					}

					if 0.8 > random.Float64() {
						mutated2 = simulatedAnnealing(child2, temp, coolingRate, instance)
					} else {
						mutated2 = hillClimbing(child2, temp, instance)
					}

					newIndividuals = append(newIndividuals, mutated1, mutated2)
				} else {
					newIndividuals = append(newIndividuals, child1, child2)
				}
			}
			newIndividuals = append(newIndividuals, parent1, parent2)
		}

		// Survivor selection -- ELITISM

		// grows populationx
		newPopulation := deepCopyPopulation(population.applyElitismWithPercentage(newIndividuals, elitismPercentage))
		population = newPopulation

		population.printPopulationStats()
		printBestIndividual(population.Individuals, instance)
		fmt.Println("GENEREATION", generation+1)

		bestFitness := getBestIndividual(population.Individuals).Fitness
		if bestFitness == lastFitness {
			stuck--
		} else {
			lastFitness = bestFitness
		}

		if stuck > 5 {
			fmt.Println("\nPERFORM DESTRUCTION AND REBUILD POPULATION..\n")
			newPopulation := deepCopyPopulation(population.applyGenecoideWithElitism(elitismPercentage, instance))
			population = newPopulation
			stuck = 0
		}

		generation++
	}

	getBestIndividual(population.Individuals).writeIndividualToJson()
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
