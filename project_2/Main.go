package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// https://it3708.resolve.visma.com/

// Declare what file you want problem instance from
var train_file string = "train/train_2.json"

// GA paramters
var numParents int = 50
var populationSize int = 100
var crossoverRate float64 = 0.6
var mutationRate float64 = 1
var gMax int = 1000
var temp int = 1000
var coolingRate float64 = 0.5
var elitismPercentage float64 = 0.01
var annealingRate float64 = 1

func main() {
	fmt.Println("")
	instance := getProblemInstance(train_file)

	pop := initPopulation(instance, 10)
	testInd := pop.Individuals[0]

	printSolution(testInd, instance)

	GA(populationSize, gMax, numParents, temp, crossoverRate, mutationRate,
		elitismPercentage, coolingRate, annealingRate, instance)

	fmt.Println("\nTRAVEL TIME =", instance.getTravelTime(0, 1))
	fmt.Println("Calculating fintess from json..", readFromJson(instance))
}

/*
PARAMTERS FOR TRAIN 8
var numParents int = 50
var populationSize int = 100
var crossoverRate float64 = 0.8
var mutationRate float64 = 0.1
var gMax int = 50
var temp int = 100
*/

/*
TA help

Island Mode - Feks, hver 25 gen lar man et individ fra en øy flytte til en annen, for å introdusere diversity. Øyene
vil være stuck forskjellige steder.

Niching og crowding - trenger diversity



*/

func readFromJson(instance Instance) float64 {
	file, err := os.Open("plotting/IndividualVisma.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return 0.0
	}
	defer file.Close()

	// Read JSON data from file
	var jsonData [][]int
	err = json.NewDecoder(file).Decode(&jsonData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return 0.0
	}

	var fitness float64 = 0
	for _, subArray := range jsonData {
		start := 0
		for _, element := range subArray {
			fitness += instance.getTravelTime(start, element)
			start = element
		}
		fitness += instance.getTravelTime(start, 0)
	}
	return fitness
}
