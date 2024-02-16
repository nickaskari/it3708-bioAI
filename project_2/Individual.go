package main

func generateIndividualSolution(instance Instance) []Route {
    // Initialize nurses and routes
    nurses := make([]Nurse, instance.Nbr_nurses)
    for i := range nurses {
        nurses[i] = Nurse{Capacity: instance.Capacity_nurse}
    }

    routes := make([]Route, len(nurses))
    for i := range routes {
        routes[i] = Route{
            Depot:    instance.Depot,
            Nurse:    nurses[i],
            Patients: make([]Patients, 0),
        }
    }

    // Distribute patients among nurses
    for _, patient := range instance.Patients {
		for i := range routes {
			nurseRoute := &routes[i] // Work with a reference to the route

			// Check if the nurse has enough capacity and can visit within the patient's time window
			if nurseRoute.Nurse.Capacity >= patient.Demand {
				
				// Need to actually calculate travel time based on matrix
				arrivalTime := nurseRoute.Nurse.CurrentTime 

				// Check if the nurse can start care within the patient's time window and finish before the end time
				if arrivalTime >= patient.Start_time && (arrivalTime + patient.Care_time) <= patient.End_time {

					nurseRoute.Patients = append(nurseRoute.Patients, patient)

					nurseRoute.Nurse.Capacity -= patient.Demand

					nurseRoute.Nurse.CurrentTime = arrivalTime + patient.Care_time

					break // Move to the next patient 
				}
			}
		}
	}

    return routes
}