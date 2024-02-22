package main

import (
	"fmt"
)

// Declare what file you want problem instance from
var train_file string = "train/train_0.json"

func main() {
	instance := getProblemInstance(train_file)

	fmt.Println("")

	population := initPopulation(instance, 100)
	population.printPopulationStats()

	population.printBestIndividual(instance)

	/*
	individual := createIndividual(instance)

	printSolution(individual, instance)
	individual.writeIndividualToJson()
	*/
}
