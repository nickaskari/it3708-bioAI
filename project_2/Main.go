package main

import (
	"fmt"
)

// Declare what file you want problem instance from
var train_file string = "train/train_0.json"

// GA paramters
var numParents int = 25
var populationSize int = 100
var crossoverRate float64 = 0.2

func main() {
	fmt.Println("")
	instance := getProblemInstance(train_file)

	population := initPopulation(instance, populationSize)
	population.printBestIndividual(instance)

	//population.printPopulationStats()
	//population.printBestIndividual(instance)
	//population.printBestIndividual(instance)
	//population.BestIndividual.checkIndividualRoutes(instance, false)

	//fmt.Println("\n -----------------------CHILD INCOMING -----------------------------")

	parents := population.tournamentSelection(numParents)
	fmt.Println("Im here 2")
	destroyRepairCrossover(parents[0], parents[1], instance)

	//child := mpic(parents, numParents, instance, crossoverRate)
	//printSolution(child, instance)
	//child.checkIndividualRoutes(instance, false)
}
