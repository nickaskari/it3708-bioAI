package main

import (
	"fmt"
	"encoding/json"
	"os"
)

/*
NOTE:
	An individual is simply a feasible solution for this problem.
	In other words an array of routes for each nurse
	such that it is a valid solution.
*/

type Individual struct {
	Fitness float64		`json:"fitness"`
	Routes  []Route		`json:"routes"`
}

// Takes in an individual and writes the struct in a JSON-file
func (i Individual) writeIndividualToJson() {
	jsonData, err := json.MarshalIndent(i, "", "    ")
    if err != nil {
        fmt.Printf("Error marshaling to JSON: %v", err)
    }

    err = os.WriteFile("Individual.json", jsonData, 0644)
    if err != nil {
        fmt.Printf("Error writing JSON to file: %v", err)
    }
}