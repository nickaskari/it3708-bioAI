package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

/*
NOTE:
	An individual is simply a feasible solution for this problem.
	In other words an array of routes for each nurse
	such that it is a valid solution.
*/

type Individual struct {
	Fitness float64 `json:"fitness"`
	Age     int     `json:"age"`
	Routes  []Route `json:"routes"`
}

// Takes in an individual and writes the struct in a JSON-file
func (i Individual) writeIndividualToJson() {
	jsonData, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		fmt.Printf("Error marshaling to JSON: %v", err)
	}

	err = os.WriteFile("plotting/Individual.json", jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing JSON to file: %v", err)
	}
}

// Calculates the fitness of an individual, and assignes a penalty if neccesary. Updated the individuals fitness.
func (i *Individual) calculateFitness(instance Instance) {
	var fitness float64 = 0
	for _, route := range i.Routes {
		if len(route.Patients) > 0 {
			lastLocation := 0
			for pNum, patient := range route.Patients {
				fitness += instance.getTravelTime(lastLocation, pNum)
				lastLocation = pNum
				fitness += calculatePenalty(patient)
			}
		}
	}
	i.Fitness = fitness
}

// Calculates a penalty if patient is visited after endtime, or nurse leaves after endtime.
func calculatePenalty(patient Patient) float64 {
	var penaltyFactor float64 = 10

	if patient.VisitTime > float64(patient.EndTime) {
		return (patient.VisitTime - float64(patient.EndTime)) * penaltyFactor
	} else if patient.LeavingTime > float64(patient.EndTime) {
		return (patient.LeavingTime - float64(patient.EndTime)) * penaltyFactor
	}
	return 0
}

// Function that makes individual one generation older.
func (i *Individual) growOlder() {
	i.Age++
}

func (individual Individual) checkIndividualRoutes(instance Instance) {
routesLoop:
	for routeIndex, route := range individual.Routes {
		for _, patient := range route.Patients {
			if !violatesTimeWindowConstraints(route, patient, instance) {
				// If any constraint is violated, it's printed inside checkRouteConstraints
				fmt.Println(strings.Repeat("*", 75) + " VIOLATION " + strings.Repeat("*", 75))
				fmt.Printf("Route %d has constraints violations with patient %d\n", routeIndex+1, patient.ID)
				fmt.Println(strings.Repeat("*", 75) + " VIOLATION " + strings.Repeat("*", 75))

				break routesLoop
			}
		}
	}
}

func violatesTimeWindowConstraints(route Route, patient Patient, instance Instance) bool {

	if patient.VisitTime < float64(patient.StartTime) {
		//fmt.Printf("Route has a patient it starts treating before its treating window \n")
		return false
	}

	if patient.VisitTime > float64(patient.EndTime) {
		//fmt.Printf("Route has a patient it starts treating after its treating window \n")
		return false
	}

	if patient.LeavingTime > float64(patient.EndTime) {
		//fmt.Printf("Route has a patient it ends treating after its treating window \n")
		return false
	}

	if patient.LeavingTime < float64(patient.StartTime) {
		//fmt.Printf("Route has a patient it ends treating before its treating window \n")
		return false
	}

	return true
}
