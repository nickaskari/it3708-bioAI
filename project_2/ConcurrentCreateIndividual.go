package main

import (
	"fmt"
	"sync"
)

// Acts as a constructor for creating an individual.
func createIndividualConcurrent(instance Instance) Individual {
    routes := createInitialRoutes(instance)
    totalTravelTime := 0.0
    var mu sync.Mutex // For synchronizing access to shared variables

    visitedPatients := make([]Patient, 0)
    patients := instance.getPatients()
    patientChan := make(chan Patient, len(patients)) // Channel for patients to process
    doneChan := make(chan bool)                      // Channel to signal completion

    // Start a goroutine for each patient
    for _, patient := range patients {
        go func(p Patient) {
            availableRoutes := sliceToMap(routes) // Assume this is a thread-safe operation or make it safe
            searchForRoute := true

            for searchForRoute {
                route, routeIndex := getRandomRoute(availableRoutes) // Ensure this is thread-safe
				
                if satisfiesConstraints(route, p, instance) {
                    mu.Lock()
                    travelTime := visitPatient(routes, routeIndex, p, instance) // Ensure visitPatient is thread-safe
                    totalTravelTime += travelTime
                    visitedPatients = append(visitedPatients, p)
                    mu.Unlock()

                    searchForRoute = false
                } else {
                    mu.Lock()
                    availableRoutes = removeRouteFromMap(availableRoutes, routeIndex) // Ensure thread-safety
                    mu.Unlock()
                }

                mu.Lock()
                if len(availableRoutes) == 0 {
                    mu.Unlock()
                    fmt.Println("Restarting due to no available routes")
                    patientChan <- p // Send patient back to channel if no route available
                    return
                }
                mu.Unlock()
            }

            doneChan <- true // Signal completion for this patient
        }(patient)
    }

    // Wait for all goroutines to complete
    for i := 0; i < len(patients); i++ {
        <-doneChan
    }

    routes = returnToDepot(routes, instance) // Ensure thread-safety if necessary
    return Individual{totalTravelTime, routes}
}