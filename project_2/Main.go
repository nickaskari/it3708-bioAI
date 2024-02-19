package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	data := readJSON("train/train_0.json")
	// Since readJSON does not return an error, we can only check if data is not nil
	if data == nil {
		fmt.Println("Failed to read or empty JSON data")
		return
	}

	var instance Instance
	err := json.Unmarshal(data, &instance)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %s\n", err)
		return
	}

	patients := getPatients("train/train_0.json", instance, data)

	for _, patient := range patients {
		fmt.Printf("Patient: %+v\n", patient)
	}
	fmt.Println(len(patients))

	/*
			matrix, err := getTravelTimeMatrix("train/train_0.json")
		    if err != nil {
		        fmt.Println("Error getting travel times matrix:", err)
		        return
		    }
			fmt.Println(matrix)
	*/

	/*
			NbrNurses = instance.Nbr_nurses
			CapacityNurse = instance.Capacity_nurse
			Benchmark = instance.Benchmark
			fmt.Printf("NbrNurses: %d, CapacityNurse: %d, Benchmark: %f\n", NbrNurses, CapacityNurse, Benchmark)

			depot := instance.Depot
		    fmt.Printf("Depot: %+v\n", depot)

			patients := instance.Patients
			fmt.Printf("Patients: %+v\n", patients)
	*/

	//individualSolution := generateIndividualSolution(instance)
	//fmt.Println("Generated Individual Solution:", individualSolution)

	//CalculateDistance(instance.Travel_times, '0', '1')
}
