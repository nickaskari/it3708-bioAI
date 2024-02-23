package main

import (
	"time"
	"math/rand"
)

type Route struct {
	Depot Depot `json:"depot"`
	NurseCapacity int `json:"nurse_capacity"`
	CurrentTime float64 `json:"current_time"`
	Patients []Patient `json:"patients"`
}

// Outputs the current location for nurse. 0 means depot. 1, 2, 3, .. denotes the patient.
func (r Route) getCurrentLocation() int {
	if len(r.Patients) == 0 {
		return 0
	} else {
		lastPatientID := r.Patients[len(r.Patients) - 1].ID
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

func createRouteFromPatientsVisited(patients []Patient, instance Instance) Route {
	var currentTime float64 = 0
	lastLocation := 0
	for _, patient := range patients {
		currentTime += instance.getTravelTime(lastLocation, patient.ID) 
		if currentTime < float64(patient.StartTime) {
			currentTime += float64(patient.StartTime) - currentTime
		} 
		patient.VisitTime = currentTime
		currentTime += float64(patient.CareTime)
		patient.LeavingTime = currentTime

		lastLocation = patient.ID
	}
	// Go back to depot
	currentTime += instance.getTravelTime(lastLocation, 0)

	return Route{
		Depot:          instance.Depot,
		NurseCapacity:  instance.CapacityNurse,
		CurrentTime:    currentTime,
		Patients:       patients,
	}
}





