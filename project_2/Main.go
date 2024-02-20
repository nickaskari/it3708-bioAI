package main

import (
	"fmt"
)

// Declare what file you want problem instance from
var train_file string = "train/train_0.json"

func main() {
    instance := getProblemInstance(train_file)

	patients := getPatients(instance)

	for _, patient := range patients {
		fmt.Printf("Patient: %+v\n", patient)
	}
	fmt.Println(len(patients))
}
