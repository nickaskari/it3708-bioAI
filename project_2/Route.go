package main

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




