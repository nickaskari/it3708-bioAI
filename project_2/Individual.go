package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"reflect"
	"slices"
	"strings"
	"time"
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

// Prints violations of an individual
func (individual Individual) checkIndividualRoutes(instance Instance, force bool) {
	timeWindowViolation := Violation{Count: 0, Example: ""}
	capacityViolation := Violation{Count: 0, Example: ""}
	returnTimeViolation := Violation{Count: 0, Example: ""}

	distinctVisitedPatients := make([]Patient, 0)
	totalPatients := make([]Patient, 0)
	for routeIndex, route := range individual.Routes {
		demandCovered := 0
		for _, patient := range route.Patients {
			demandCovered += patient.Demand
			if !patient.IsPatientInList(distinctVisitedPatients) {
				distinctVisitedPatients = append(distinctVisitedPatients, patient)
			}
			totalPatients = append(totalPatients, patient)
			if !notViolatesTimeWindowConstraints(route, patient, instance) {
				if timeWindowViolation.Example == "" {
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
	if timeWindowViolation.Count > 0 || capacityViolation.Count > 0 || returnTimeViolation.Count > 0 || force {
		reportViolation(timeWindowViolation, capacityViolation, returnTimeViolation, distinctVisitedPatients, totalPatients, instance)
	}
}

// Helper function for checkIndividualRoutes(). Prints out a report of violations.
func reportViolation(timeWindow Violation, capacity Violation, returnTime Violation, visitedPatients []Patient, totalPatients []Patient, instance Instance) {
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

	if len(visitedPatients) != len(instance.PatientArray) {
		fmt.Println("\nNumber of distinct patients", len(visitedPatients), "is not correct!")
	}
	fmt.Println("\nTotal number of patients", len(totalPatients))

	fmt.Println(strings.Repeat("*", 75) + " VIOLATION " + strings.Repeat("*", 75))
}

// Helper function for checkIndividualRoutes() that checks for time window constraint violations
func notViolatesTimeWindowConstraints(route Route, patient Patient, instance Instance) bool {

	if patient.VisitTime < float64(patient.StartTime) {
		return false
	}

	if patient.VisitTime > float64(patient.EndTime) {
		return false
	}

	if patient.LeavingTime > float64(patient.EndTime) {
		return false
	}

	if patient.LeavingTime < float64(patient.StartTime) {
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

// Deep copy method for Individual
func deepCopyIndividual(original Individual) Individual {
	copy := original

	// Deep copy any slices, maps, or other reference types
	copy.Routes = make([]Route, len(original.Routes))
	for i, route := range original.Routes {
		// If Route itself contains reference types, you'll need to deep copy those too
		copy.Routes[i] = deepCopyRoute(route) // Assuming Route can be shallow copied; adjust if necessary
	}

	// Deep copy other fields as necessary
	copy.Fitness = original.Fitness
	copy.Age = original.Age

	return copy
}

func (i *Individual) addRoute(route Route) {
	i.Routes = append(i.Routes, route)
}

// return num patients in individual, and if duplicates
func (i Individual) getNumPatients() (int, bool) {
	num := 0
	visited := []int{}
	for _, r := range i.Routes {
		for _, p := range r.Patients {
			num++
			visited = append(visited, p.ID)
		}
	}
	return num, hasDuplicates(visited)
}

// Finds the worst cost route within an individual. Returns the an int array of
func (i Individual) findWorstCostRoute(instance Instance) []int {
	var worstCostRouteIndex int

	worstFitness := math.Inf(-1)
	for index, r := range i.Routes {
		routeFitness := calculateRouteFitness(r, instance)
		if routeFitness > worstFitness {
			worstFitness = routeFitness
			worstCostRouteIndex = index
		}
	}

	return i.Routes[worstCostRouteIndex].extractAllVisitedPatients()
}

/*
Removes a list of patients from the routes in the individual. Does NOT recalculate individual fitness.
Does NOT update patient attributes!
*/
func (i *Individual) removePatients(patientsToRemove []int, instance Instance) {
	removed := 0
	for rIndex, r := range i.Routes {
		var newPatients []Patient
		modified := false
		for _, p := range r.Patients {
			if !slices.Contains(patientsToRemove, p.ID) {
				newPatients = append(newPatients, p)
			} else {
				removed++
				modified = true
			}
		}
		if modified {
			r.Patients = newPatients
			i.Routes[rIndex] = r
		}
		if removed == len(patientsToRemove) {
			break
		}
	}
}

// Adds patients from list to random routes. Does NOT fix patient values for routes!
func (i *Individual) distributePatientsOnRoutes(patients []int, instance Instance) {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	for _, pID := range patients {
		patientAdded := false

		for !patientAdded {
			randomRouteIndex := random.Intn(instance.NbrNurses)
			newRoute, ok := i.Routes[randomRouteIndex].canAddPatientEnforced(pID, instance)
			if ok {
				i.Routes[randomRouteIndex] = newRoute
				patientAdded = true
			}
		}
	}
}

// Fixes all routes to contain the correct values. Also calculates and updates fitness.
func (i *Individual) fixAllRoutesAndCalculateFitness(instance Instance) {
	for index, r := range i.Routes {
		i.Routes[index] = createRouteFromPatientsVisited(r.Patients, instance)
	}
	i.calculateFitness(instance)
}

// Inserts a list of patients in the best routes. Recalculates fitness in the end.
func (i *Individual) findBestRoutesForPatients(patients []int, instance Instance) {
	for _, pID := range patients {
		leastChange := math.Inf(1)
		var index int
		var newRoute Route
		for rIndex, r := range i.Routes {
			possibleRoute, change := r.findBestInsertion(pID, instance)
			if change < leastChange {
				leastChange = change
				newRoute = possibleRoute
				index = rIndex
			}
		}
		i.Routes[index] = newRoute
	}
	i.calculateFitness(instance)
}
