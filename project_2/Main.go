package main

import (
	"fmt"
)

// Declare what file you want problem instance from
var train_file string = "train/train_0.json"

// GA paramters
var numParents int = 5
var populationSize int = 100
var crossoverRate float64 = 0.2

func main() {
	instance := getProblemInstance(train_file)

	population := initPopulation(instance, populationSize)

	//population.printPopulationStats()
	//population.printBestIndividual(instance)
	population.BestIndividual.checkIndividualRoutes(instance)

	fmt.Println("\n ------------------CHILD INCOMING --------------------------")

	parents := population.tournamentSelection(numParents)
	child := mpic(parents, numParents, instance, crossoverRate)
	printSolution(child, instance)
	child.checkIndividualRoutes(instance)
}
