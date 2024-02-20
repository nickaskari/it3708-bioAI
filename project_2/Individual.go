package main

import (
	"fmt"
)

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
var random Randomizer

// Acts as a constructor for creating an individual.
func createIndividual(instance Instance) Individual {
	routes := createInitialRoutes(instance)
	total_travel_time := 0.0

	visited_patients := make([]Patient, 0)
	for _, patient := range instance.Patients {

		if !patient.IsPatientInList(visited_patients) {

			available_routes := routes
			search_for_route := true

			for search_for_route {
				route, route_index := random.getRandomRoute(routes)

				if satisfiesConstraints(route, patient, instance) {
					total_travel_time += visitPatient(route, patient, instance)
					visited_patients = append(visited_patients, patient)
					search_for_route = false
				} else {
					available_routes = append(available_routes[:route_index], available_routes[route_index+1:]...)
				}

				if len(available_routes) == 0 {
					return createIndividual(instance)
				}
			}

		}
	}

	return Individual{total_travel_time, routes}
}

// Creates an array of Nbr_nurses number of routes. Each initialized with t=0 and zero patients.
func createInitialRoutes(instance Instance) []Route {
	routes := make([]Route, (instance.Nbr_nurses))
	for i := range routes {
		routes[i] = Route{
			Depot:         instance.Depot,
			Nurse_capacity: instance.Capacity_nurse,
			Current_time:   0,
			Patients:      make([]Patient, 0),
		}
	}
	return routes
}

// Checks whether a given nurse can visit a potential patient.
func satisfiesConstraints(nurseRoute Route, potentialPatient Patient, instance Instance) bool {

	last_patient_id := nurseRoute.Patients[len(nurseRoute.Patients)-1].ID
	potentialPatient_id := potentialPatient.ID

	currentTime := nurseRoute.Current_time

	last_patient_to_depot := instance.getTravelTime(last_patient_id, 0)
	from_curent_to_potential_patient := instance.getTravelTime(last_patient_id, potentialPatient_id)

	// should we check whether a patient has a start and end time such that the care time is sufficient?
	if (nurseRoute.Nurse_capacity >= potentialPatient.Demand) &&
		(currentTime+float64(potentialPatient.Care_time) <= float64(potentialPatient.End_time)) &&
		(nurseRoute.Current_time + from_curent_to_potential_patient+last_patient_to_depot <= float64(instance.Depot.ReturnTime)) {
			return true
	} else {
		return false
	}
}

// Visit a patient and wait and/or care for them. Returns travel time.
func visitPatient(route Route, patient Patient, instance Instance) float64{
	route.Patients = append(route.Patients, patient)
	route.Nurse_capacity -= patient.Demand

	current_time := route.Current_time

	if current_time < float64(patient.Start_time) {
		waitingTime := float64(patient.Start_time) - current_time
		current_time += waitingTime
	}

	last_visited_patient_ID := route.Patients[len(route.Patients) - 1].ID
	travel_time := instance.getTravelTime(last_visited_patient_ID, patient.ID)
	route.Current_time += float64(patient.Care_time) + travel_time

	return travel_time
}
