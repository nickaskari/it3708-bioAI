package main

import (
	"fmt"
)

// Declare what file you want problem instance from
var train_file string = "train/train_0.json"

// GA paramters
var numParents int = 5
var populationSize int = 100

func main() {
	instance := getProblemInstance(train_file)

	fmt.Println("")

	population := initPopulation(instance, populationSize)

	population.printPopulationStats()
	population.printBestIndividual(instance)
	population.BestIndividual.checkIndividualRoutes(instance)
	fmt.Println("\n\n1\n\n")
	var bro Individual = inversionMutationIndividual(population.BestIndividual, instance)
	for i := 0; i < 100000; i ++ {

		bro = inversionMutationIndividual(bro, instance)
		bro = swapMutationIndividual(bro, instance)
	
	}

	printSolution(bro, instance)
	bro.checkIndividualRoutes(instance)

}
