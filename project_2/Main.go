package main

import (
	"fmt"
)

// Declare what file you want problem instance from
var train_file string = "train/train_0.json"

// GA paramters
var numParents int = 5
var populationSize int = 10

func main() {
	instance := getProblemInstance(train_file)

	fmt.Println("")

	population := initPopulation(instance, populationSize)

	population.printPopulationStats()
	population.printBestIndividual(instance)
	population.BestIndividual.checkIndividualRoutes(instance)

	printDivider(105, "<-><->")

	bro := swapMutation(population.BestIndividual, instance)
	printSolution(bro, instance)
	bro.checkIndividualRoutes(instance)

}
