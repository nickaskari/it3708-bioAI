package main

import (
	"math"
	"time"
	"math/rand"
)

// Performs edge recombination of two parents. Returns offspring Route. 
func edgeRecombination(parentRoute1 Route, parentRoute2 Route, instance Instance) Route {
	
	if len(parentRoute1.Patients) == 0 && len(parentRoute2.Patients) == 0 {
        return initalizeOneRoute(instance) 
    }

	// If either parent route is empty, return the other parent route
	if len(parentRoute1.Patients) == 0 {
		return parentRoute2
	}
	if len(parentRoute2.Patients) == 0 {
		return parentRoute1
	}

	// This way of dealing with empty routes makes sure that satisfiesConstraints never gets called with an empty route

	// If both parent routes are nonempty, perform edge recombination
	matrix1 := createEdgeConnectivityMatrix(parentRoute1)
	matrix2 := createEdgeConnectivityMatrix(parentRoute2)

	unionMatrix := matrixUnionEdges(matrix1, matrix2)

	offspringRoute := initalizeOneRoute(instance)

	// Select a random starting patient from one of the parent routes
	// Need to handle empty routes
	currentPatient := parentRoute1.getRandomPatient()

	for len(offspringRoute.Patients) < getTotalNumberOfUniquePatients(unionMatrix) {
		if satisfiesConstraints(offspringRoute, currentPatient, instance) { 
			offspringRoute.visitPatient(currentPatient, instance)
		
			removePatientFromAdjacencyMatrix(currentPatient, unionMatrix)
	
			var nextPatient Patient
	
			if len(unionMatrix[currentPatient.ID]) > 0 {
				nextPatient = getLeastConnectedNeighbor(currentPatient, unionMatrix)
	
			} else {
				nextPatient = selectRandomRemainingPatient(unionMatrix)
			}
	
			currentPatient = nextPatient
		}
	}

	return offspringRoute
}


// Creates an adjcancy matrix of patients visited from a certain patient. Start and end patient are "tied" together.
func createEdgeConnectivityMatrix(route Route) map[int][]Patient {
	matrix := make(map[int][]Patient, 0)
	patients := route.Patients

	for i, patient := range patients {
		connection := make([]Patient, 0)
		if i == 0 {
			connection = append(connection, patients[len(patients)-1], patients[i+1])
		} else if i == (len(patients) - 1) {
			connection = append(connection, patients[0], patients[i-1])
		} else {
			connection = append(connection, patients[i-1], patients[i+1])
		}
		matrix[patient.ID] = connection
	}
	return matrix
}

// matrixUnionEdges takes the union of two adjacency matrices (m1 and m2) and returns the union matrix.
func matrixUnionEdges(m1, m2 map[int][]Patient) map[int][]Patient {
	union := make(map[int][]Patient)

	// Helper function to add patients to the union map safely (avoiding duplicates).
	addPatients := func(key int, patients []Patient) {
		for _, patient := range patients {
			if !patient.IsPatientInList(union[key]) {
				union[key] = append(union[key], patient)
			}
		}
	}

	// Add all patients from m1 to the union.
	for key, patients := range m1 {
		addPatients(key, patients)
	}

	// Add all patients from m2 to the union. This will also add new keys from m2 that were not in m1.
	for key, patients := range m2 {
		addPatients(key, patients)
	}

	return union
}

// Returns total number of unique patients from an adjacency matrix
func getTotalNumberOfUniquePatients(matrix map[int][]Patient) int {
	uniquePatients := make([]Patient, 0)
	for _, patients := range matrix {
		for _, patient := range patients {
			if !patient.IsPatientInList(uniquePatients) {
				uniquePatients = append(uniquePatients, patient)
			}
		}
	}

	return len(uniquePatients)
}

// Remove a patient from a all lists within adjacency matrix
func removePatientFromAdjacencyMatrix(patient Patient, matrix map[int][]Patient) {
	for key, patients := range matrix {
		patient.deletePatientFrom(patients)
		matrix[key] = patient.deletePatientFrom(patients)
	}
}

// Choose the next patient from the current patient's connections Prefer patients with the fewest connections in the union matrix for the next step.
func getLeastConnectedNeighbor(patient Patient, matrix map[int][]Patient) Patient {
	neighbors := matrix[patient.ID]

	leastConnections := math.Inf(1)
	candidates := make([]Patient, 0)
	for _, neighbor := range neighbors {
		numNeighbors := len(matrix[neighbor.ID])
		if numNeighbors < int(leastConnections) {
			candidates = []Patient{neighbor}
		} else if numNeighbors == int(leastConnections) {
			candidates = append(candidates, neighbor)
		}
	}

	if len(candidates) > 0 {
		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)
	
		randomIndex := random.Intn(len(candidates))
	
		return candidates[randomIndex]
	} else {
		return candidates[0]
	}
}

// Select and return random remaining patient
func selectRandomRemainingPatient(matrix map[int][]Patient) Patient {
	for _, patients := range matrix {
		if len(patients) > 0 {
			return patients[0]
		}
	}
	return createDummyPatient()
}


