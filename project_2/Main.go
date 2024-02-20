package main

import (
	"fmt"
)

// Declare what file you want problem instance from
var train_file string = "train/train_3.json"

func main() {
    instance := getProblemInstance(train_file)
	
	/*
	patients := instance.getPatients()

	for _, patient := range patients {
		fmt.Printf("Patient: %+v\n", patient)
	}
	fmt.Println(len(patients))*/
	
	fmt.Println("")
	individual := createIndividual(instance)
	printSolution(individual, instance)
	//fmt.Println(instance.getTravelTime(49, 0))

	individual.writeIndividualToJson()
}
