package main


// Performs crossover between two parents with the destroy/repair method as detailed in the visma lecture
func destroyRepairCrossover(parent1 Individual, parent2 Individual, instance Instance) (Individual, Individual) {
	// Finds per now the route with worst objectve value. Maybe want to change to most inefficient route in the future. 
	worstRoutePatients1 := parent1.findWorstCostRoute(instance) 
	worstRoutePatients2 := parent2.findWorstCostRoute(instance)

	offspring1 := deepCopyIndividual(parent1)
	offspring2 := deepCopyIndividual(parent2)

	offspring1.Age, offspring2.Age = 0, 0

	offspring1.removePatients(worstRoutePatients2, instance)
	offspring2.removePatients(worstRoutePatients1, instance)

	// Updates Route values and assignes fitness
	offspring1.fixAllRoutesAndCalculateFitness(instance)
	offspring2.fixAllRoutesAndCalculateFitness(instance)

	// create function that distributes patients on routes based on best insertion
	offspring1.findBestRoutesForPatients(worstRoutePatients2, instance)
	offspring2.findBestRoutesForPatients(worstRoutePatients1, instance)

	return offspring1, offspring2
}



