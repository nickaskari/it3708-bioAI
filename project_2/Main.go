package main

import (
	"fmt"
)

// Declare what file you want problem instance from
var train_file string = "train/train_8.json"

// GA paramters
var numParents int = 60
var populationSize int = 100
var crossoverRate float64 = 1
var mutationRate float64 = 0.8
var gMax int = 20
var temp int = 200

func main() {
	fmt.Println("")
	instance := getProblemInstance(train_file)

	GA(populationSize, gMax, numParents, temp, crossoverRate, mutationRate, instance)

}
