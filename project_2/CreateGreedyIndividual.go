package main

import (
	"fmt"
	"math"
)

// Acts as a constructor for creating an individual.
func createGreedyIndividual(instance Instance) Individual {
	routes := createInitialRoutes(instance)
	totalTravelTime := 0.0

	allPatients := instance.getPatients()
	visitedPatients := make([]Patient, 0)
	for _, route := range routes {
		if len(visitedPatients) != len(allPatients) {
			searchForPatient := true

			for searchForPatient {
				patient := findClosestPatient(route, availableRoutes, instance)

				if satisfiesConstraints(route, patient, instance) {

					travelTime := visitPatient(routes, routeIndex, patient, instance)
					totalTravelTime += travelTime
					visitedPatients = append(visitedPatients, patient)
					searchForPatient = false
					//fmt.Println("IM HERE 1")
				} else {
					// If the route does not satisfy constraints, remove it from the set of possible routes.
					availableRoutes = removeRouteFromMap(availableRoutes, routeIndex)
					//fmt.Println("IM HERE 2", len(availableRoutes))
				}

				if len(availableRoutes) == 0 {
					// If there no routes that satisfies the constraints, start from scratch
					fmt.Println("visited", len(visitedPatients))
					return createIndividual(instance)

					// UNCOMMMENT TO GET "HALFWAY" SOLUTION WHEN RUNNING
					//routes = returnToDepot(routes, instance)
					//return Individual{totalTravelTime, routes}
				}
			}
		}

	}
	fmt.Println("visited", len(visitedPatients))
	routes = returnToDepot(routes, instance)

	return Individual{totalTravelTime, routes}
}

// Finds the closest nurse for a given patient.
func findClosestNurse(patient Patient, routes map[int]Route, instance Instance) (Route, int) {
	closestDistance := math.Inf(1)
	closestRoute := routes[0]
	closestRouteIndex := 0

	for index, route := range routes {
		currentLocation := route.getCurrentLocation()
		distance := instance.getTravelTime(currentLocation, patient.ID)

		if distance < closestDistance {
			closestRoute = route
			closestRouteIndex = index
		}
	}

	return closestRoute, closestRouteIndex

}


