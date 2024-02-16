package main



type Route struct {
	Depot Depot
	Nurse Nurse
	// strore the patients in the route, depot start and finish is implicit
	Patients []Patients
}

