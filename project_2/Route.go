package main

import (
	"math"
	"math/rand"
	"time"
)

type Route struct {
	Depot         Depot     `json:"depot"`
	NurseCapacity int       `json:"nurse_capacity"`
	CurrentTime   float64   `json:"current_time"`
	Patients      []Patient `json:"patients"`
}

// Outputs the current location for nurse. 0 means depot. 1, 2, 3, .. denotes the patient.
func (r Route) getCurrentLocation() int {
	if len(r.Patients) == 0 {
		return 0
	} else {
		lastPatientID := r.Patients[len(r.Patients)-1].ID
		return lastPatientID
	}
}

// Returns a random patient from route.
func (r Route) getRandomPatient() Patient {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	randomIndex := random.Intn(len(r.Patients))

	return r.Patients[randomIndex]
}

// Creates a route based on a slice of patients visited. Updates patients visit/leave time.
func createRouteFromPatientsVisited(patients []Patient, instance Instance) Route {
	var currentTime float64 = 0
	lastLocation := 0
	capacity := instance.CapacityNurse
	newPatients := make([]Patient, 0)

	for _, patient := range patients {
		patientCopy := instance.getPatientAtID(patient.ID)

		currentTime += instance.getTravelTime(lastLocation, patientCopy.ID)
		if currentTime < float64(patientCopy.StartTime) {
			currentTime += float64(patientCopy.StartTime) - currentTime
		}
		patientCopy.VisitTime = currentTime

		currentTime += float64(patientCopy.CareTime)
		patientCopy.LeavingTime = currentTime

		lastLocation = patientCopy.ID
		newPatients = append(newPatients, patientCopy)

		capacity -= patientCopy.Demand
	}
	// Go back to depot
	currentTime += instance.getTravelTime(lastLocation, 0)

	return Route{
		Depot:         instance.Depot,
		NurseCapacity: capacity,
		CurrentTime:   currentTime,
		Patients:      newPatients,
	}
}

// Returns a route with currentTime = 0 and zero patients.
func initalizeOneRoute(instance Instance) Route {
	return Route{
		Depot:         instance.Depot,
		NurseCapacity: instance.CapacityNurse,
		CurrentTime:   0,
		Patients:      make([]Patient, 0),
	}
}

// Visits a patient to a Route. Updates currentTime
func (r *Route) visitPatient(patient Patient, instance Instance) {
	currentLocation := r.getCurrentLocation()

	// Travel
	r.CurrentTime += instance.getTravelTime(currentLocation, patient.ID)
	// Wait
	if r.CurrentTime < float64(patient.StartTime) {
		r.CurrentTime += float64(patient.StartTime) - r.CurrentTime
	}
	// Visit
	patient.VisitTime = r.CurrentTime
	// Care
	r.CurrentTime += float64(patient.CareTime)
	// Leave
	patient.LeavingTime = r.CurrentTime

	r.NurseCapacity -= patient.Demand

	r.Patients = append(r.Patients, patient)
}

// Deep copy function for Route
func deepCopyRoute(originalRoute Route) Route {
	var r Route
	r.Depot = originalRoute.Depot
	r.NurseCapacity = originalRoute.NurseCapacity
	r.CurrentTime = originalRoute.CurrentTime

	// Manually copying the slice
	r.Patients = make([]Patient, len(originalRoute.Patients))
	copy(r.Patients, originalRoute.Patients) // This is correct usage of copy for slice

	return r
}

// Calculates fitness of route. Returns fitness
func calculateRouteFitness(route Route, instance Instance) float64 {
	var fitness float64 = 0
	if len(route.Patients) > 0 {
		lastLocation := 0
		for _, patient := range route.Patients {
			fitness += instance.getTravelTime(lastLocation, patient.ID)
			lastLocation = patient.ID
			fitness += calculatePenalty(patient)
		}

		fitness += instance.getTravelTime(lastLocation, 0)
	}
	return fitness
}

// Outputs all patient ID's visited
func (r Route) extractAllVisitedPatients() []int {
	visited := make([]int, 0)
	for _, p := range r.Patients {
		visited = append(visited, p.ID)
	}
	return visited
}

/*
Checks whether a patient can be added to a route. Checks capacity.
Returns Route and bool on whether this can indeed happen.
*/
func (r Route) canAddPatient(patientID int, instance Instance) (Route, bool) {
	patientToAdd := instance.getPatientAtID(patientID)

	demandCovered := 0
	//canReturnInTime := false
	finalRoute := initalizeOneRoute(instance)
	for _, patient := range r.Patients {
		demandCovered += patient.Demand

		/*
			if !canReturnInTime {
				newPatientOrder := r.Patients
				newPatientOrder = append(newPatientOrder[:index+1], newPatientOrder[index:]...)
				newPatientOrder[index] = patientToAdd

				newRoute := createRouteFromPatientsVisited(newPatientOrder, instance)

				if newRoute.CurrentTime <= float64(instance.Depot.ReturnTime) {
					finalRoute = newRoute
					canReturnInTime = true
				}
			}*/
	}

	if instance.CapacityNurse <= (demandCovered + patientToAdd.Demand) {
		return finalRoute, false
	} else {
		oldPatients := r.Patients
		newPatients := append(oldPatients, patientToAdd)
		return createRouteFromPatientsVisited(newPatients, instance), true
	}
}

/*
Checks whether a patient can be added to a route. Checks capacity AND returntime constraints.
Returns Route and bool on whether this can indeed happen.
*/
func (r Route) canAddPatientEnforced(patientID int, instance Instance) (Route, bool) {
	patientToAdd := instance.getPatientAtID(patientID)

	demandCovered := 0
	canReturnInTime := false
	finalRoute := initalizeOneRoute(instance)
	for index, patient := range r.Patients {
		demandCovered += patient.Demand

		if canReturnInTime == false {
			routeCopy := deepCopyRoute(r)
			newPatientOrder := routeCopy.Patients
			newPatientOrder = append(newPatientOrder[:index+1], newPatientOrder[index:]...)
			newPatientOrder[index] = patientToAdd
			newRoute := createRouteFromPatientsVisited(newPatientOrder, instance)

			if newRoute.CurrentTime <= float64(instance.Depot.ReturnTime) {
				finalRoute = newRoute
				canReturnInTime = true
			}
		}
	}

	if instance.CapacityNurse > (demandCovered+patientToAdd.Demand) && canReturnInTime {
		return finalRoute, true
	} else {
		return finalRoute, false
	}
}

// Checks if route contains duplicate patients. Returns array of patient id visited, and bool on whether there is duplicate
func (r Route) checkIfRouteContainsDuplicates() ([]int, bool) {
	visited := []int{}
	for _, p := range r.Patients {
		visited = append(visited, p.ID)
	}
	return visited, hasDuplicates(visited)
}

// Finds best insertion of patient in a route. Outputs new route and route objective value changed (after - before)
func (r Route) findBestInsertion(patientID int, instance Instance) (Route, float64) {
	patient := instance.getPatientAtID(patientID)

	//if r.NurseCapacity < patient.Demand {
	//	return r, math.Inf(1)
	//}

	oldRouteFitness := calculateRouteFitness(r, instance)
	var bestRoute Route
	changedObjectiveValue := math.Inf(1)

	// Handle Empty Route
	if len(r.Patients) == 0 {
		newPatientOrder := []Patient{patient}
		newRoute := createRouteFromPatientsVisited(newPatientOrder, instance)
		// check for return time violation
		//if newRoute.CurrentTime <= float64(instance.Depot.ReturnTime) {
		newRouteFitness := calculateRouteFitness(newRoute, instance)
		return newRoute, newRouteFitness
		//}
	}

	for index := 0; index < len(r.Patients); index++ {
		routeCopy := deepCopyRoute(r)
		newPatientOrder := routeCopy.Patients
		newPatientOrder = append(newPatientOrder[:index+1], newPatientOrder[index:]...)

		if index == len(r.Patients)-1 {
			newPatientOrder[index+1] = patient
		} else {
			newPatientOrder[index] = patient
		}

		newRoute := createRouteFromPatientsVisited(newPatientOrder, instance)

		// check for return time violation
		//if newRoute.CurrentTime <= float64(instance.Depot.ReturnTime) {
			newRouteFitness := calculateRouteFitness(newRoute, instance)

			change := newRouteFitness - oldRouteFitness
			if change < changedObjectiveValue {
				changedObjectiveValue = change
				bestRoute = newRoute
			}
		//}
	}

	return bestRoute, changedObjectiveValue
}

/*
Performs a patient swap for inter route mutations.
Takes in patient from current route, changes it with the spesified patient.
Returns the changed route. Performs swap only if capacity and returntime are not violated.
*/
func (r Route) performPatientSwap(exsistingPatient Patient, outsidePatient Patient, instance Instance) (Route, bool) {
	newPatients := []Patient{}

	for _, p := range r.Patients {
		if p.ID == exsistingPatient.ID {
			newPatients = append(newPatients, instance.getPatientAtID(outsidePatient.ID))
		} else {
			newPatients = append(newPatients, instance.getPatientAtID(p.ID))
		}
	}

	newRoute := createRouteFromPatientsVisited(newPatients, instance)
	if newRoute.NurseCapacity <= instance.CapacityNurse {
		if newRoute.CurrentTime <= float64(instance.Depot.ReturnTime) {
			return newRoute, true
		}
	}
	return r, false
}
