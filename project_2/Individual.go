package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
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
	timeWindowViolation := Violation{Count: 0, Example: ""}
	capacityViolation := Violation{Count: 0, Example: ""}
	returnTimeViolation := Violation{Count: 0, Example: ""}

	distinctVisitedPatients := make([]Patient, 0)
	for routeIndex, route := range individual.Routes {
		demandCovered := 0
		for _, patient := range route.Patients {
			demandCovered += patient.Demand
			if !patient.IsPatientInList(distinctVisitedPatients) {
				distinctVisitedPatients = append(distinctVisitedPatients, patient)
			}
			if !violatesTimeWindowConstraints(route, patient, instance) {
				if (timeWindowViolation.Example == "") {
					timeWindowViolation.registerExample(fmt.Sprintf("Route %d violates time window of patient %d.", routeIndex+1, patient.ID))
				}
				timeWindowViolation.countViolation()
			}
		}
		if demandCovered > instance.CapacityNurse {
			capacityViolation.registerExample(fmt.Sprintf("Route %d violates capacity.", routeIndex+1))
			capacityViolation.countViolation()
		}
		if route.CurrentTime > float64(instance.Depot.ReturnTime) {
			returnTimeViolation.registerExample(fmt.Sprintf("Route %d violates return time.", routeIndex+1))
			returnTimeViolation.countViolation()
		}
	}
	if timeWindowViolation.Count > 0 || capacityViolation.Count > 0 || returnTimeViolation.Count > 0 {
		reportViolation(timeWindowViolation, capacityViolation, returnTimeViolation, distinctVisitedPatients, instance)
	}

}

// Helper function for checkIndividualRoutes(). Prints out a report of violations.
func reportViolation(timeWindow Violation, capacity Violation, returnTime Violation, visitedPatients []Patient, instance Instance) {
	const consoleWidth = 150
	const countLabel = "Count = "
	const padding = 2

	countWidth := len(countLabel) + 20

	fmt.Println(strings.Repeat("*", 75) + " VIOLATION " + strings.Repeat("*", 75))

	printViolation := func(v Violation) {
		exampleWidth := consoleWidth - countWidth - padding
		exampleFmt := fmt.Sprintf("%%-%ds%%%ds%%d\n", exampleWidth, 30)
		fmt.Printf(exampleFmt, v.Example, countLabel, v.Count)
	}

	printViolation(timeWindow)
	printViolation(capacity)
	printViolation(returnTime)

	sumViolations := timeWindow.Count + capacity.Count + returnTime.Count
	fmt.Printf("\nSum violations = %d\n", sumViolations)

	if len(visitedPatients) != len(instance.getPatients()) {
		fmt.Println("\nNumber of distinct patinets", len(visitedPatients), "is not correct!")
	}

	fmt.Println(strings.Repeat("*", 75) + " VIOLATION " + strings.Repeat("*", 75))
}

// Helper function for checkIndividualRoutes() that checks for time window constraint violations
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

// isEqual compares two Individual instances for equality.
func isEqual(a, b Individual) bool {
	// Compare all other fields for equality.
	// Assuming Routes is the only field for simplicity; add comparisons for any other fields as necessary.
	return reflect.DeepEqual(a.Routes, b.Routes)
}

// RemoveIndividualFrom removes the first occurrence of `individual` from `individuals`.
func (individual Individual) removeIndividualFrom(individuals []Individual) []Individual {
	for i, j := range individuals {
		if isEqual(j, individual) {
			return append(individuals[:i], individuals[i+1:]...)
		}
	}
	return individuals
}
