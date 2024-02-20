package main

import (
	"fmt"
)

// Acts as a constructor for creating an individual.
func createIndividual(instance Instance) Individual {
	routes := createInitialRoutes(instance)
	totalTravelTime := 0.0

	visitedPatients := make([]Patient, 0)
	for _, patient := range instance.getPatients() {

		// Assigning patients to routes(nurses), not the other way around
		if !patient.IsPatientInList(visitedPatients) {

			availableRoutes := sliceToMap(routes)
			searchForRoute := true

			for searchForRoute {
				route, routeIndex := getRandomRoute(availableRoutes)

				if satisfiesConstraints(route, patient, instance) {

					travelTime := visitPatient(routes, routeIndex, patient, instance)
					totalTravelTime += travelTime
					visitedPatients = append(visitedPatients, patient)
					searchForRoute = false
					//fmt.Println("IM HERE 1")
				} else {
					// If the route does not satisfy constraints, remove it from the set of possible routes.
					availableRoutes = removeRouteFromMap(availableRoutes, routeIndex)
					//fmt.Println("IM HERE 2", len(availableRoutes))
				}

				if len(availableRoutes) == 0 {
					// If there no routes that satisfies the constraints, start from scratch
					fmt.Println("visited", len(visitedPatients))
					return createIndividual(instance)

					// UNCOMMMENT TO GET "HALFWAY" SOLUTION WHEN RUNNING
					//routes = returnToDepot(routes, instance)
					//return Individual{totalTravelTime, routes}
				}
			}
		}
	}
	fmt.Println("visited", len(visitedPatients))
	routes = returnToDepot(routes, instance)

	return Individual{totalTravelTime, routes}
}

// Creates an array of Nbr_nurses number of routes. Each initialized with t=0 and zero patients.
func createInitialRoutes(instance Instance) []Route {
	routes := make([]Route, (instance.NbrNurses))
	for i := range routes {
		routes[i] = Route{
			Depot:          instance.Depot,
			NurseCapacity: instance.CapacityNurse,
			CurrentTime:   0,
			Patients:       make([]Patient, 0),
		}
	}
	return routes
}

// Checks whether a given nurse can visit a potential patient.
func satisfiesConstraints(nurseRoute Route, potentialPatient Patient, instance Instance) bool {

	currentPatient := 0
	if len(nurseRoute.Patients) > 0 {
		currentPatient = nurseRoute.Patients[len(nurseRoute.Patients)-1].ID
	}

	potentialPatientID := potentialPatient.ID

	//currentTime := nurseRoute.CurrentTime

	potentialPatientToDepot := instance.getTravelTime(potentialPatient.ID, 0)
	curentToPotentialPatient := instance.getTravelTime(currentPatient, potentialPatientID)

	if (nurseRoute.NurseCapacity < potentialPatient.Demand) {
		return false
	}

	timeAtArival := nurseRoute.CurrentTime + curentToPotentialPatient
	if timeAtArival < float64(potentialPatient.StartTime) {
		if (float64(potentialPatient.StartTime) + float64(potentialPatient.CareTime) + potentialPatientToDepot) <= float64(instance.Depot.ReturnTime) {
			// Assuming end - start >= caretime. Hence if care is at start time the nurse will be done in time for the end time. 
			return true
		}
	} else {
		if (timeAtArival + float64(potentialPatient.CareTime) + potentialPatientToDepot) <= float64(instance.Depot.ReturnTime) {
			// THIS ONE IS TOO STRICT??
			if (timeAtArival + float64(potentialPatient.CareTime)) <= float64(potentialPatient.EndTime) {
				// If the nurse arrives late, check if the nurse will treat in time before the end time.
				return true
			}
		}
	}
	return false

	// CHECK THESE CONDITIONS --> ARE THEY CORRECT?
	/*
	if (nurseRoute.NurseCapacity >= potentialPatient.Demand) && 
		(potentialPatient.EndTime - potentialPatient.StartTime >= potentialPatient.CareTime) && //Not scam
		(currentTime + float64(potentialPatient.CareTime) <= float64(potentialPatient.EndTime)) &&
		(nurseRoute.CurrentTime + curentToPotentialPatient + float64(potentialPatient.CareTime) + potentialPatientToDepot <=
			float64(instance.Depot.ReturnTime)) {
		return true
	} else {
		return false
	}*/
}

// Visit a patient and wait and/or care for them. Returns travel time and route.
func visitPatient(routes []Route, index int, patient Patient, instance Instance) float64 {
	routes[index].NurseCapacity -= patient.Demand

	lastVisitedPatientID := 0
	if len(routes[index].Patients) > 0 {
		lastVisitedPatientID = routes[index].Patients[len(routes[index].Patients)-1].ID
	}

	travelTime := instance.getTravelTime(lastVisitedPatientID, patient.ID)

	// Travel
	routes[index].CurrentTime += travelTime

	// Wait if neccesary
	if routes[index].CurrentTime < float64(patient.StartTime) {
		waitingTime := float64(patient.StartTime) - routes[index].CurrentTime
		routes[index].CurrentTime += waitingTime
		fmt.Println("ROUTE NUM", index, "WAITED")
	}

	// Now you can visit the patient
	patient.VisitTime = routes[index].CurrentTime 

	// Care for patient
	routes[index].CurrentTime += float64(patient.CareTime)

	// Time to leave
	patient.LeavingTime = routes[index].CurrentTime

	routes[index].Patients = append(routes[index].Patients, patient)

	return travelTime
}

// Takes in all routes, checks if they are not empty, then return those to the depot.
func returnToDepot(routes []Route, instance Instance) []Route {
	for i, route := range routes {
		patients := route.Patients
		if len(patients) != 0 {
			lastPatientID := patients[len(patients)-1].ID
			travelTimeToDepot := instance.getTravelTime(lastPatientID, 0)
			routes[i].CurrentTime += travelTimeToDepot
		}
	}
	return routes
}

// Removes the route at index spesified from an array of routes.
func removeRouteFromArray(routes []Route, index int) []Route {
	if index < len(routes) - 1 {
		routes = append(routes[:index], routes[index+1:]...)
	} else {
		// It's the last element or out of bounds; just truncate the slice if it's the last element
		routes = routes[:index]
	}
	return routes
}

