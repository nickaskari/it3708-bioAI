package main

import (
	"fmt"
)

// Performs crossover between two parents with the destroy/repair method as detailed in the visma lecture
func destroyRepairCrossover(parent1 Individual, parent2 Individual, instance Instance) {
	worstRoutePatients1 := parent1.findWorstCostRoute(instance)
	worstRoutePatients2 := parent2.findWorstCostRoute(instance)

	offspring1 := deepCopyIndividual(parent1)
	offspring2 := deepCopyIndividual(parent2)

	offspring1.Age, offspring2.Age = 0, 0

	// When the routes of parent 1 and 2 get manipulated, the patient values, return times, and individual fitness are not fixed
	// These should be taken care of in the local search algorithm?

	offspring1.removePatients(worstRoutePatients2, instance)
	offspring2.removePatients(worstRoutePatients1, instance)

	// Assign these removed patients to new routes randomly.
	offspring1.distributePatientsOnRoutes(worstRoutePatients2, instance)
	offspring2.distributePatientsOnRoutes(worstRoutePatients1, instance)


	// Updates Route values and assignes fitness
	offspring1.fixAllRoutesAndCalculateFitness(instance)
	offspring2.fixAllRoutesAndCalculateFitness(instance)

	// Perform local search with mutations

	fmt.Println("yo")
}



