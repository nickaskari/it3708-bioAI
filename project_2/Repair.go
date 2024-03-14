package main

// Takes in an individual, finds a cluster in some route, repairs it with best insertion. Returns repair individual

func destroyRepairCluster(individual Individual, instance Instance) Individual {
	repairedIndividual := deepCopyIndividual(individual)

	for i := 0; i < instance.NbrNurses; i ++ {
		routeIndex := repairedIndividual.getRandomRoute()

		patientsInCluster := repairedIndividual.Routes[routeIndex].findPatientsInCluster(instance)
	
		repairedIndividual.removePatients(patientsInCluster, instance)
	
		// Updates Route values and assignes fitness
		repairedIndividual.fixAllRoutesAndCalculateFitness(instance)
	
		// create function that distributes patients on routes based on best insertion
		repairedIndividual.findBestRoutesForPatients(patientsInCluster, instance)
	}

	return repairedIndividual
}

// DESTROYS ONE RANDOM ROUTE, AND PERFORMS BEST COST INSERTION ON IT
func destroyRepair(individual Individual, instance Instance) Individual {
	repairedIndividual := deepCopyIndividual(individual)

	for i := 0; i < instance.NbrNurses; i ++ {
		routeIndex := repairedIndividual.getRandomRoute()

		patients := repairedIndividual.Routes[routeIndex].extractAllVisitedPatients()
	
		repairedIndividual.removePatients(patients, instance)
	
		// Updates Route values and assignes fitness
		repairedIndividual.fixAllRoutesAndCalculateFitness(instance)
	
		// create function that distributes patients on routes based on best insertion
		repairedIndividual.findBestRoutesForPatients(patients, instance)
	}

	return repairedIndividual
}

// DESTROYS ONE RANDOM ROUTE, AND PERFORMS RANDOM INSERTION ON IT
func destroyRepairRandomly(individual Individual, instance Instance) Individual {
	repairedIndividual := deepCopyIndividual(individual)

	for i := 0; i < instance.NbrNurses; i ++ {
		routeIndex := repairedIndividual.getRandomRoute()

		patients := repairedIndividual.Routes[routeIndex].extractAllVisitedPatients()
	
		repairedIndividual.removePatients(patients, instance)
	
		// Updates Route values and assignes fitness
		repairedIndividual.fixAllRoutesAndCalculateFitness(instance)
	
		// create function that distributes patients on routes based on best insertion
		repairedIndividual.distributePatientsOnRoutes(patients, instance)
	}

	return repairedIndividual
}

// Removes a few patients, AND PERFORMS BEST COST INSERTION ON IT
func destroyRepaiLite(individual Individual, instance Instance) Individual {
	repairedIndividual := deepCopyIndividual(individual)

	for i:=0;i<2;i++ {
		patients := generateRandomPatientIDs(instance)

		repairedIndividual.removePatients(patients, instance)
	
		// Updates Route values and assignes fitness
		repairedIndividual.fixAllRoutesAndCalculateFitness(instance)
	
		// create function that distributes patients on routes based on best insertion
		repairedIndividual.findBestRoutesForPatients(patients, instance)
	}

	return repairedIndividual
}