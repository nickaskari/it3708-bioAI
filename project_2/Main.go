package main

import (
	"fmt"
)

// Declare what file you want problem instance from
var train_file string = "train/train_7.json"

// GA paramters
var numParents int = 40
var populationSize int = 100
var crossoverRate float64 = 0.8
var mutationRate float64 = 0.2
var gMax int = 100
var temp int = 1000

func main() {
	fmt.Println("")
	instance := getProblemInstance(train_file)

	GA(populationSize, gMax, numParents, temp, crossoverRate, mutationRate, instance)

}
