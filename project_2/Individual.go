package main

import (
	"encoding/json"
	"fmt"
	"os"
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
    for routeIndex, route := range individual.Routes {
        for _, patient := range route.Patients {
            if !checkRouteConstraints(route, patient, instance) {
                // If any constraint is violated, it's printed inside checkRouteConstraints
                fmt.Printf("Route %d has constraints violations with patient ID %d\n", routeIndex+1, patient.ID)
            }
        }
    }
}

func checkRouteConstraints(nurseRoute Route, potentialPatient Patient, instance Instance) bool {
    currentPatientID := 0
    if len(nurseRoute.Patients) > 0 {
        currentPatientID = nurseRoute.Patients[len(nurseRoute.Patients)-1].ID
    }

    potentialPatientID := potentialPatient.ID
    potentialPatientToDepot := instance.getTravelTime(potentialPatientID, 0) 
    currentToPotentialPatient := instance.getTravelTime(currentPatientID, potentialPatientID)

    // Check nurse capacity constraint
    if nurseRoute.NurseCapacity < potentialPatient.Demand {
        fmt.Println("Constraint violated: Nurse capacity is less than the patient's demand.")
        return false
    }

    // Calculate time of arrival and evaluate timing constraints
    timeAtArrival := nurseRoute.CurrentTime + currentToPotentialPatient
    if timeAtArrival < float64(potentialPatient.StartTime) {
        // Nurse arrives before start time
        if (float64(potentialPatient.StartTime) + float64(potentialPatient.CareTime) + potentialPatientToDepot) > float64(instance.Depot.ReturnTime) {
            fmt.Println("Constraint violated: Nurse cannot return to depot in time after providing care.")
            return false
        }
    } else {
        // Nurse arrives at or after start time
        if (timeAtArrival + float64(potentialPatient.CareTime) + potentialPatientToDepot) > float64(instance.Depot.ReturnTime) {
            fmt.Println("Constraint violated: Nurse cannot return to depot in time after providing care.")
            return false
        }
		if (timeAtArrival + float64(potentialPatient.CareTime)) > float64(potentialPatient.EndTime) {
			// If the nurse arrives late, check if the nurse will treat in time before the end time.
			return false
		}
    }

    // If we reach here, no constraints are violated for this patient
    return true
}
