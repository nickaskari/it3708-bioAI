package main

import "fmt"

/*
NOTE:
	An individual is simply a feasible solution for this problem.
	In other words an array of routes for each nurse
	such that it is a valid solution.
*/

type Individual struct {
	fitness float64
	routes  []Route
}

// A randomizer struct for arbitrary use
var random *Randomizer = NewRandomizer()

// Acts as a constructor for creating an individual.
func createIndividual(instance Instance) Individual {
	routes := createInitialRoutes(instance)
	totalTravelTime := 0.0

	visitedPatients := make([]Patient, 0)
	for _, patient := range instance.getPatients() {

		// Assigning patients to routes(nurses), not the other way around
		if !patient.IsPatientInList(visitedPatients) {

			availableRoutes := routes
			searchForRoute := true

			for searchForRoute {
				route, routeIndex := random.getRandomRoute(availableRoutes)

				if satisfiesConstraints(route, patient, instance) {
					totalTravelTime += visitPatient(routes, routeIndex, patient, instance)
					visitedPatients = append(visitedPatients, patient)
					searchForRoute = false
					//fmt.Println("IM HERE 1")
				} else {
					// If the route does not satisfy constraints, remove it from the set of possible routes.
					availableRoutes = removeRouteFromArray(availableRoutes, routeIndex)
					//fmt.Println("IM HERE 2", len(available_routes))
				}

				if len(availableRoutes) == 0 {
					// If there no routes that satisfies the constraints, start from scratch
					fmt.Println("IM HERE 3")
					return createIndividual(instance)
				}
			}
		}
	}
	if len(visitedPatients) < len(instance.getPatients()) {
		fmt.Println("TOO FEW PATIENTS", len(visitedPatients), "INSTEAD OF", instance.getPatients())
	} else if len(visitedPatients) == len(instance.getPatients()) {
		fmt.Println("WE GOOD", len(visitedPatients), "=", len(instance.getPatients()))
	}
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

	currentTime := nurseRoute.CurrentTime

	potentialPatientToDepot := instance.getTravelTime(potentialPatient.ID, 0)
	curentToPotentialPatient := instance.getTravelTime(currentPatient, potentialPatientID)

	// CHECK THESE CONDITIONS --> ARE THEY CORRECT?
	if (nurseRoute.NurseCapacity >= potentialPatient.Demand) &&
		(potentialPatient.EndTime - potentialPatient.StartTime >= potentialPatient.CareTime) && //Not scam
		(currentTime + float64(potentialPatient.CareTime) <= float64(potentialPatient.EndTime)) &&
		(nurseRoute.CurrentTime + curentToPotentialPatient + float64(potentialPatient.CareTime) + potentialPatientToDepot <=
			float64(instance.Depot.ReturnTime)) {
		return true
	} else {
		return false
	}
}

// Visit a patient and wait and/or care for them. Returns travel time and route.
func visitPatient(routes []Route, routeIndex int, patient Patient, instance Instance) float64 {
	route := routes[routeIndex]
	route.NurseCapacity -= patient.Demand


	if route.CurrentTime < float64(patient.StartTime) {
		waitingTime := float64(patient.StartTime) - route.CurrentTime
		route.CurrentTime += waitingTime
	}

	lastVisitedPatientID := 0
	if len(route.Patients) > 0 {
		lastVisitedPatientID = route.Patients[len(route.Patients)-1].ID
	}

	patient.VisitTime = route.CurrentTime

	travelTime := instance.getTravelTime(lastVisitedPatientID, patient.ID)
	route.CurrentTime += travelTime + float64(patient.CareTime)

	patient.LeavingTime = route.CurrentTime

	route.Patients = append(route.Patients, patient)
	routes[routeIndex] = route

	return travelTime
}

// Takes in all routes, checks if they are not empty, then return those to the depot.
func returnToDepot(routes []Route, instance Instance) []Route {
	for _, route := range routes {
		patients := route.Patients
		if len(patients) != 0 {
			lastPatientID := patients[len(patients)-1].ID
			travelTimeToDepot := instance.getTravelTime(lastPatientID, 0)
			route.CurrentTime += travelTimeToDepot
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

