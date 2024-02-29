package main

/*
	Performs education (local search) by doing mutations if there is improved objective value. Returns
	educated individual. Purpose of this is to only find LOCAL optima.
*/
// understand why this kills people when zero violations
func hillClimbing(individual Individual, temp int, instance Instance) Individual {
	currentState := deepCopyIndividual(individual)

	for temp != 0 {
		invertedIndividual := randomInversionMutation(currentState, instance)
		invertedIndividual.fixAllRoutesAndCalculateFitness(instance)

		if invertedIndividual.Fitness <= currentState.Fitness {
			currentState = invertedIndividual
			swappedIndividual := randomSwapMutation(currentState, instance)
			swappedIndividual.fixAllRoutesAndCalculateFitness(instance)

			if swappedIndividual.Fitness <= currentState.Fitness {
				currentState = swappedIndividual

				interSwappedIndividual := randomInterRouteSwapMutation(currentState, instance)
				interSwappedIndividual.fixAllRoutesAndCalculateFitness(instance)

				if interSwappedIndividual.Fitness <= currentState.Fitness {
					currentState = interSwappedIndividual
				}
			}
		}
		currentState.fixAllRoutesAndCalculateFitness(instance)

		temp--
	}
	return currentState
}
