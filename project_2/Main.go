package main

import (
	"fmt"
)

// Declare what file you want problem instance from
var train_file string = "train/train_6.json"

// GA paramters
var numParents int = 150
var populationSize int = 300
var crossoverRate float64 = 0.8
var mutationRate float64 = 0.9
var gMax int = 200
var temp int = 1000
var coolingRate float64 = 0.9
var elitismPercentage float64 = 0.01

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
