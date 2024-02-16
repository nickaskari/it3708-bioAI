package main



type Route struct {
	Depot Depot
	Nurse Nurse
	// List of all the patients that the nurse will visit, and the order in which they will be visited
	Patients []Patients

}