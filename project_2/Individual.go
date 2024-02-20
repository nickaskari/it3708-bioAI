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
	total_travel_time := 0.0

	visited_patients := make([]Patient, 0)
	for _, patient := range instance.Patients {

		// Assigning patients to routes(nurses), not the other way around
		if !patient.IsPatientInList(visited_patients) {

			available_routes := routes
			search_for_route := true

			for search_for_route {
				route, route_index := random.getRandomRoute(available_routes)

				if satisfiesConstraints(route, patient, instance) {
					total_travel_time += visitPatient(routes, route_index, patient, instance)
					visited_patients = append(visited_patients, patient)
					search_for_route = false
					//fmt.Println("IM HERE 1")
				} else {
					// If the route does not satisfy constraints, remove it from the set of possible routes.
					available_routes = removeRouteFromArray(available_routes, route_index)
					//fmt.Println("IM HERE 2", len(available_routes))
				}

				if len(available_routes) == 0 {
					// If there no routes that satisfies the constraints, start from scratch
					fmt.Println("IM HERE 3")
					return createIndividual(instance)
				}
			}
		}
	}
	routes = returnToDepot(routes, instance)

	return Individual{total_travel_time, routes}
}

// Creates an array of Nbr_nurses number of routes. Each initialized with t=0 and zero patients.
func createInitialRoutes(instance Instance) []Route {
	routes := make([]Route, (instance.Nbr_nurses))
	for i := range routes {
		routes[i] = Route{
			Depot:          instance.Depot,
			Nurse_capacity: instance.Capacity_nurse,
			Current_time:   0,
			Patients:       make([]Patient, 0),
		}
	}
	return routes
}

// Checks whether a given nurse can visit a potential patient.
func satisfiesConstraints(nurseRoute Route, potentialPatient Patient, instance Instance) bool {

	current_patient := 0
	if len(nurseRoute.Patients) > 0 {
		current_patient = nurseRoute.Patients[len(nurseRoute.Patients)-1].ID
	}

	potentialPatient_id := potentialPatient.ID

	currentTime := nurseRoute.Current_time

	potential_patient_to_depot := instance.getTravelTime(potentialPatient.ID, 0)
	from_curent_to_potential_patient := instance.getTravelTime(current_patient, potentialPatient_id)

	// should we check whether a patient has a start and end time such that the care time is sufficient?
	if (nurseRoute.Nurse_capacity >= potentialPatient.Demand) &&
		(potentialPatient.End_time-potentialPatient.Start_time >= potentialPatient.Care_time) && //Not scam
		(currentTime+float64(potentialPatient.Care_time) <= float64(potentialPatient.End_time)) &&
		(nurseRoute.Current_time+from_curent_to_potential_patient+float64(potentialPatient.Care_time)+potential_patient_to_depot <=
			float64(instance.Depot.ReturnTime)) {
		return true
	} else {
		return false
	}
}

// Visit a patient and wait and/or care for them. Returns travel time and route.
func visitPatient(routes []Route, route_index int, patient Patient, instance Instance) float64 {
	route := routes[route_index]

	route.Patients = append(route.Patients, patient)
	route.Nurse_capacity -= patient.Demand

	current_time := route.Current_time

	if current_time < float64(patient.Start_time) {
		waitingTime := float64(patient.Start_time) - current_time
		current_time += waitingTime
	}

	last_visited_patient_ID := 0
	if len(route.Patients) > 0 {
		last_visited_patient_ID = route.Patients[len(route.Patients)-1].ID
	}

	travel_time := instance.getTravelTime(last_visited_patient_ID, patient.ID)
	route.Current_time += float64(patient.Care_time) + travel_time

	routes[route_index] = route

	return travel_time
}

// Takes in all routes, checks if they are not empty, then return those to the depot.
func returnToDepot(routes []Route, instance Instance) []Route {
	for _, route := range routes {
		patients := route.Patients
		if len(patients) != 0 {
			last_patient_id := patients[len(patients)-1].ID
			travel_time_to_depot := instance.getTravelTime(last_patient_id, 0)
			route.Current_time += travel_time_to_depot
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
