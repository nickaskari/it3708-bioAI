package main

type Route struct {
	Depot Depot `json:"depot"`
	NurseCapacity int `json:"nurse_capacity"`
	CurrentTime float64 `json:"current_time"`
	Patients []Patient `json:"patients"`
}



