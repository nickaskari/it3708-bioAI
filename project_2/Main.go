package main

import (
	"fmt"
)

// Declare what file you want problem instance from
var train_file string = "train/train_6.json"

// GA paramters
var numParents int = 10
var populationSize int = 100
var crossoverRate float64 = 0.8
var mutationRate float64 = 0.2
var gMax int = 200
var temp int = 10000
var coolingRate float64 = 0.99
var elitismPercentage float64 = 0.01
var annealingRate float64 = 1

func main() {
	fmt.Println("")
	instance := getProblemInstance(train_file)

	GA(populationSize, gMax, numParents, temp, crossoverRate, mutationRate,
		elitismPercentage, coolingRate, annealingRate, instance)

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
