package main


import (
	"fmt"
)


/*
NOTE:
	An individual is simply a feasible solution for this problem. In other words an array of routes for each nurse
	such that it is a valid solution.
*/

type Individual struct {
	fitness float64
	routes  []Route
}


// the total time it takes for all the nurses to visit all the patients
func (individual Individual) calculateFitness() int {
	totalTime := 0	
	for _, route := range individual.routes {
		totalTime += int(route.CurrentTime)
	}
	return totalTime
}

// A randomizer struct for arbitrary use
var random Randomizer

// Acts as a constructor for creating an individual.
func createIndividual(instance Instance) Individual {
	routes := createInitialRoutes(instance)

	visited_patients := make([]Patient, 0)
		for num, patient := range instance.Patients {
			fmt.Println("This is patient number", num)

			if !patient.IsPatientInList(visited_patients) {

				available_routes := routes
				search_for_route := true

				for search_for_route {
					route, route_index := random.getRandomRoute(routes)

					if satisfiesConstraints(route, patient, instance) {
						route.Patients = append(route.Patients, patient)
						route.NurseCapacity -= patient.Demand
						route.CurrentTime += float64(patient.Care_time) // + Travel time
					
						visited_patients = append(visited_patients, patient)
						search_for_route = false
					} else {
						available_routes = append(available_routes[:route_index], available_routes[route_index+1:]...)
					}

					if len(available_routes) == 0 {
						createIndividual(instance)
					}
				}
				
			}
		}
		newIndivual := Individual{0, routes}
		newIndivual.calculateFitness()

		return newIndivual
	}

func createInitialRoutes(instance Instance) []Route {
	routes := make([]Route, (instance.Nbr_nurses))
	for i := range routes {
		routes[i] = Route{
			Depot:         instance.Depot,
			NurseCapacity: instance.Capacity_nurse,
			CurrentTime:   0,
			Patients:      make([]Patient, 0),
		}
	}
	return routes
}


func satisfiesConstraints(nurseRoute Route, potentialPatient Patient, instance Instance) bool {
	
	last_patient_id := nurseRoute.Patients[len(nurseRoute.Patients)-1].ID	
	potentialPatient_id := potentialPatient.ID

	currentTime := nurseRoute.CurrentTime

	last_patient_to_depot := instance.getTravelTime(last_patient_id, 0)
	from_curent_to_potential_patient := instance.getTravelTime(last_patient_id, potentialPatient_id)

	if (nurseRoute.NurseCapacity >= potentialPatient.Demand) &&
		(int(currentTime) + potentialPatient.Care_time <= potentialPatient.End_time) &&
		(nurseRoute.CurrentTime + from_curent_to_potential_patient + last_patient_to_depot <= float64(instance.Depot.ReturnTime)) {
			if (currentTime >= float64(potentialPatient.Start_time)) {
				// DONT WAIT
			} else {
				// WAIT 
			}

		}
	return true
}

