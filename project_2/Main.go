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

	population := initPopulation(instance, populationSize)

	population.printPopulationStats()
	population.printBestIndividual(instance)
	population.BestIndividual.checkIndividualRoutes(instance)

	parents := population.tournamentSelection(numParents)
	child := crossover(parents[0], parents[1], instance, 0)
	printSolution(child, instance)
	child.checkIndividualRoutes(instance)
	fmt.Println("")
}
