package main

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

// https://it3708.resolve.visma.com/

// Declare what file you want problem instance from
var train_file string = "train/train_9.json"

// Benchmark stop criteria. 0 essentially deactivates this.
var benchmark float64 = 0
var migrationFrequency int = 25
var numMigrants int = 10
var gmax int = 1000

// GA paramters OLD
/*
var numParents int = 50
var populationSize int = 100
var crossoverRate float64 = 0.4
var mutationRate float64 = 0.8
var gMax int = 2000
var temp int = 1000
var coolingRate float64 = 0.5
var elitismPercentage float64 = 0.05
var annealingRate float64 = 1
*/

// Island parameters
var islandConfigs = []struct {
	numParents        int
	populationSize    int
	crossoverRate     float64
	mutationRate      float64
	gMax              int
	temp              int
	coolingRate       float64
	elitismPercentage float64
	annealingRate     float64
}{
	{25, 50, 0.9, 0.2, gmax, 500, 0.1, 0.05, 1},
	{25, 50, 0.9, 0.2, gmax, 1000, 0.1, 0.05, 1},
	{25, 50, 0.9, 0.2, gmax, 1000, 0.1, 0.05, 1},
	{25, 50, 0.8, 0.4, gmax, 1000, 0.1, 0.05, 1},
	{25, 50, 0.8, 0.4, gmax, 1000, 0.5, 0.05, 1},
	{25, 50, 0.8, 0.4, gmax, 1000, 0.5, 0.05, 1},
	{25, 50, 0.8, 0.4, gmax, 1000, 0.5, 0.05, 1},
	{25, 50, 0.8, 0.4, gmax, 1000, 0.5, 0.05, 1},
	{25, 50, 0.9, 0.2, gmax, 1000, 0.5, 0.05, 1},
}

func main() {
	fmt.Println("Starting GA on islands...")

	instance := getProblemInstance(train_file)

	// Use a context with cancel to signal goroutines to stop
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure all paths cancel to avoid context leak

	// Use a WaitGroup to wait for all islands to finish
	var wg sync.WaitGroup

	// Channel to collect the best individuals from each island
	bestIndividuals := make(chan Individual, len(islandConfigs))

	migrationEvent := NewMigrationEvent(len(islandConfigs))
	//migrationEvent = NewMigrationEvent(len(islandConfigs))

	for islandID, config := range islandConfigs {
		wg.Add(1)
		go func(c struct {
			numParents        int
			populationSize    int
			crossoverRate     float64
			mutationRate      float64
			gMax              int
			temp              int
			coolingRate       float64
			elitismPercentage float64
			annealingRate     float64
		}) {
			defer wg.Done()

			// Run GA on each island with its configuration and capture the best individual
			best, reachedBenchmark := GA(c.populationSize, c.gMax, c.numParents, c.temp, c.crossoverRate, c.mutationRate,
				c.elitismPercentage, c.coolingRate, c.annealingRate, benchmark, ctx, migrationFrequency, numMigrants,
				migrationEvent, islandID, instance)

			bestIndividuals <- best

			if reachedBenchmark {
				fmt.Println("BENCHMARK WAS REACHED -- EXITING ALL CURRENT GO ROUTINES..")
				fmt.Println("INDIVIDUAL WAS FOUND BY ISLAND", islandID, "AND CONFIG", config)
				cancel() // Reached benchmark, signal other goroutines to stop
			}

		}(config)
	}

	wg.Wait()              // Wait for all goroutines to finish
	close(bestIndividuals) // Close the channel after all sends are complete

	// Slice to collect all best individuals
	allBest := make([]Individual, 0, len(islandConfigs))
	for ind := range bestIndividuals {
		allBest = append(allBest, ind)
	}

	fmt.Println("\nAll islands have completed. Best individual collected:\n")

	sort.Slice(allBest, func(i, j int) bool {
		return allBest[i].Fitness < allBest[j].Fitness
	})

	best := allBest[0]

	best.fixAllRoutesAndCalculateFitness(instance)
	printSolution(best, instance)
	best.checkIndividualRoutes(instance, true)
	best.writeIndividualToJson()
	best.writeIndividualToVismaFormat()
}

/*
func main() {
	fmt.Println("")
	instance := getProblemInstance(train_file)

	GA(populationSize, gMax, numParents, temp, crossoverRate, mutationRate,
		elitismPercentage, coolingRate, annealingRate, benchmark, instance)

	fmt.Println("Calculating fintess from json..", readFromJson(instance))
}
*/
