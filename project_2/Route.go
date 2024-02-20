package main

import (
	"sort"
)

type Route struct {
	Depot Depot `json:"depot"`
	NurseCapacity int `json:"nurse_capacity"`
	CurrentTime float64 `json:"current_time"`
	Patients []Patient `json:"patients"`
}

func sortRoutesByPatientCount(routes []Route) {
    sort.Slice(routes, func(i, j int) bool {
        return len(routes[i].Patients) < len(routes[j].Patients)
    })
}


