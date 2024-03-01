package main

import (
	"fmt"
)

// Declare what file you want problem instance from
var train_file string = "train/train_5.json"

// GA paramters
var numParents int = 80
var populationSize int = 100
var crossoverRate float64 = 0.8
var mutationRate float64 = 0.2
var gMax int = 300
var temp int = 100
var coolingRate float64 = 0.9
var elitismPercentage float64 = 0.02

func main() {
	fmt.Println("")
	instance := getProblemInstance(train_file)

	GA(populationSize, gMax, numParents, temp, crossoverRate, mutationRate,
		elitismPercentage, coolingRate, instance)

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
