package main

import (
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

	newPatients := make([]Patient, 0)
	for _, patient := range patients {
		currentTime += instance.getTravelTime(lastLocation, patient.ID)
		if currentTime < float64(patient.StartTime) {
			currentTime += float64(patient.StartTime) - currentTime
		}
		patient.VisitTime = currentTime

		currentTime += float64(patient.CareTime)
		patient.LeavingTime = currentTime

		lastLocation = patient.ID
		newPatients = append(newPatients, patient)
	}
	// Go back to depot
	currentTime += instance.getTravelTime(lastLocation, 0)

	return Route{
		Depot:         instance.Depot,
		NurseCapacity: instance.CapacityNurse,
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
