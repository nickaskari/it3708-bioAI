package main

import (
	"fmt"
)

func main() {
	//data := readJSON("train/train_0.json")

	/*
	patients, err := getPatients("train/train_0.json")
    if err != nil {
        fmt.Println("Error getting patients:", err)
        return
    }

    for _, patient := range patients {
        fmt.Printf("Patient: %+v\n", patient)
    }
	fmt.Println(len(patients))
	*/

	
	matrix, err := getTravelTimeMatrix("train/train_0.json")
    if err != nil {
        fmt.Println("Error getting travel times matrix:", err)
        return
    }
	fmt.Println(matrix)	
}




/*
train/train_0.json

1. Funksjon som tar inn data, extracte patient data. Lag objekter og sett i en array.

2. Lag travel matrix. Array of Array. 

3. Store "nbr_nurses" ,"capacity_nurse" og "benchmark" som global variabler.


*/


