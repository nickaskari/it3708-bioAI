package main

type Route struct {
	Depot Depot `json:"depot"`
	Nurse_capacity int `json:"nurse_capacity"`
	Current_time float64 `json:"current_time"`
	Patients []Patient `json:"patients"`
}



