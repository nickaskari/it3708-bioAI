package main

import (
	"fmt"
	"math/rand"
	"time"
)


func GA(populationSize int, gMax int, numParents int, temp int, crossoverRate float64, mutationRate float64, instance Instance) {
	
	generation := 0
	
	fmt.Println("Initalzing population..")
	population := initPopulation(instance, populationSize)
	population.printBestIndividual(instance)

	var newIndividuals []Individual

	for generation < gMax {

		newIndividuals = make([]Individual, 0)

		parents := population.tournamentSelection(numParents)
			
			for population.size() > len(newIndividuals) {
				source := rand.NewSource(time.Now().UnixNano())
				r := rand.New(source)
				
				if r.Float64() < crossoverRate {
					i, j := getTwoRandomParents(parents)
					parent1, parent2 := parents[i], parents[j]
					child1, child2 := destroyRepairCrossover(parent1, parent2, instance)

					if r.Float64() < mutationRate {
						mutated1, mutated2 := hillClimbing(child1, temp, instance), hillClimbing(child2, temp, instance)
						newIndividuals = append(newIndividuals, mutated1, mutated2)
					} else {
						newIndividuals = append(newIndividuals, child1, child2)
					}
					
				}
			}
		
			// GHETTO solution
		fmt.Println("\n\n")
		fmt.Println("GENEREATION", generation)
		population = Population {
			Individuals: newIndividuals,
			BestIndividual: getBestIndividual(newIndividuals),
		}
		population.printPopulationStats()
		population.printBestIndividual(instance)
		fmt.Println("\n\n")

		generation++
	}
	// Perform survivor selection SKIPPING FOR NOW
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