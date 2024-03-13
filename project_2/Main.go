package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
)

// When running you can press enter --> GA will stop and output best individual so far

// https://it3708.resolve.visma.com/

// Declare what file you want problem instance from
var train_file string = "train/train_7.json"

// Benchmark stop criteria. 0 essentially deactivates this.
var benchmark float64 = 1102

var migrationFrequency int = 25
var numMigrants int = 8
var initiateBestCostRepair bool = true
var genocideWhenStuck int = 15

var gmax int = 1500
var numParents int = 25
var populationSize int = 50

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
	{numParents, populationSize, 0.8, 0.2, gmax, 500, 0.1, 0.05, 1},
	{numParents, populationSize, 0.8, 0.2, gmax, 1000, 0.1, 0.05, 1},
	{numParents, populationSize, 0.8, 0.2, gmax, 1000, 0.1, 0.05, 1},
	{numParents, populationSize, 0.7, 0.4, gmax, 1000, 0.1, 0.05, 1},
	{numParents, populationSize, 0.8, 0.4, gmax, 1000, 0.5, 0.05, 1},
	{numParents, populationSize, 0.8, 0.4, gmax, 1000, 0.5, 0.05, 1},
	{numParents, populationSize, 0.7, 0.4, gmax, 1000, 0.5, 0.05, 1},
	{numParents, populationSize, 0.8, 0.4, gmax, 1000, 0.5, 0.05, 1},
	{numParents, populationSize, 0.7, 0.2, gmax, 1000, 0.5, 0.05, 1},

	/*
		ALL TRAIN EXCEPT 5 AND 6
		var genocideWhenStuck int = 15
		var migrationFrequency int = 25
		{numParents, populationSize, 0.8, 0.2, gmax, 500, 0.1, 0.05, 1},
		{numParents, populationSize, 0.8, 0.2, gmax, 1000, 0.1, 0.05, 1},
		{numParents, populationSize, 0.8, 0.2, gmax, 1000, 0.1, 0.05, 1},
		{numParents, populationSize, 0.7, 0.4, gmax, 1000, 0.1, 0.05, 1},
		{numParents, populationSize, 0.8, 0.4, gmax, 1000, 0.5, 0.05, 1},
		{numParents, populationSize, 0.8, 0.4, gmax, 1000, 0.5, 0.05, 1},
		{numParents, populationSize, 0.7, 0.4, gmax, 1000, 0.5, 0.05, 1},
		{numParents, populationSize, 0.8, 0.4, gmax, 1000, 0.5, 0.05, 1},
		{numParents, populationSize, 0.7, 0.2, gmax, 1000, 0.5, 0.05, 1},

		"train/train_6.json", and train 5 (WHEN INDIVIDUAL CONVERGES TO LONG ROUTES THESE SHOULD BE THE PARAMTERS)
		migrationNum 8 migrationFrquency = 10
		initiateBestCostRepair bool = false
		var genocideWhenStuck int = 5
		{numParents, populationSize, 0.1, 0.5, gmax, 50, 0.1, 0.05, 1},
		{numParents, populationSize, 0.1, 0.2, gmax, 100, 0.1, 0.05, 1},
		{numParents, populationSize, 0.0, 0.2, gmax, 100, 0.1, 0.05, 1},
		{numParents, populationSize, 0.2, 0.8, gmax, 50, 0.1, 0.05, 1},
		{numParents, populationSize, 0.2, 0.8, gmax, 100, 0.9, 0.05, 1},
		{numParents, populationSize, 0.1, 0, gmax, 100, 0.5, 0.05, 1},
		{numParents, populationSize, 0.25, 0.8, gmax, 100, 0.1, 0.05, 1},

		LONGEST TRAINING TIMES ARE UP TO 5-6 MIN.
		EASY PROBLEMS GET SOLVED IN ABOUT 2 MIN.
	*/

}

func main() {
	fmt.Println("Starting GA on islands...")

	startTime := time.Now()

	instance := getProblemInstance(train_file)

	// Use a context with cancel to signal goroutines to stop
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Use a WaitGroup to wait for all islands to finish
	var wg sync.WaitGroup

	// Channel to collect the best individuals from each island
	bestIndividuals := make(chan Individual, len(islandConfigs))

	migrationEvent := NewMigrationEvent(len(islandConfigs))

	go func() {
		fmt.Println("Press Enter at any time to stop all operations.")
		bufio.NewReader(os.Stdin).ReadBytes('\n') // Block until Enter is pressed
		migrationEvent.Lock()
		defer migrationEvent.Unlock()
		migrationEvent.Ready.Broadcast()
		migrationEvent.signalCancelEvent()
		cancel() // Trigger cancellation
	}()

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
				migrationEvent, islandID, initiateBestCostRepair, genocideWhenStuck, instance)

			bestIndividuals <- best

			if reachedBenchmark {
				fmt.Println("BENCHMARK WAS REACHED -- EXITING ALL CURRENT GO ROUTINES..")
				fmt.Println("INDIVIDUAL WAS FOUND BY ISLAND", islandID, "AND CONFIG", config)
				migrationEvent.Lock()
				defer migrationEvent.Unlock()
				migrationEvent.Ready.Broadcast()
				migrationEvent.signalCancelEvent()
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

	endTime := time.Now() // Capture the end time
	duration := endTime.Sub(startTime)

	fmt.Println("\nThe Genetic Algorithm took", duration.Seconds(), "to run.")
}
