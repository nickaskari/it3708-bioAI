package main

import (
	"fmt"
	"encoding/json"
	"os"
)

// needed to read all the JSON file data
type Instance struct {
    Instance_name  string            `json:"instance_name"`
    Nbr_nurses     int               `json:"nbr_nurses"`
    Capacity_nurse int               `json:"capacity_nurse"`
    Benchmark      float64           `json:"benchmark"`
    Depot          Depot             `json:"depot"`
    Patients       map[string]Patients `json:"patients"`
    Travel_times   [][]float64       `json:"travel_times"`
}

func readJSON (filename string) []byte {
	data, err := os.ReadFile(filename)
    if err != nil {
      fmt.Print(err)
    }
	return data
}

// Two-dimensional array. Each inner array is a row in the matrix
// Matrix dimension is (n+1) X (n+1) where n is the number of patients
func getTravelTimeMatrix(filename string) ([][]float64, error) {
    data := readJSON(filename)
    if data == nil {
        // Since readJSON doesn't return an error, we need to create our own
        return nil, fmt.Errorf("failed to read JSON data from file: %s", filename)
    }

    var instance Instance
    err := json.Unmarshal(data, &instance)
    if err != nil {
        return nil, err
    }

    travelTimes := instance.Travel_times
    return travelTimes, nil
}


