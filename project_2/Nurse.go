package main

type Nurse struct {
	Capacity int `json:"capacity_nurse"`
	CurrentTime int // to keep track of the time while visiting patients
}