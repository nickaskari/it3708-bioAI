package main


// Needed to read all the JSON file data
type Instance struct {
    InstanceName  string            `json:"instance_name"`
    NbrNurses     int               `json:"nbr_nurses"`
    CapacityNurse int               `json:"capacity_nurse"`
    Benchmark      float64           `json:"benchmark"`
    Depot          Depot             `json:"depot"`
    Patients       map[string]Patient `json:"patients"`
    TravelTimes   [][]float64      `json:"travel_times"`
}


/* 
    Returns the traveltime between nurses and/or depot.
    Inputs are the row/column indicies of the nurses/depots in the
    travel matrix.
    Travel Times is the Euclidean distance between two items. 
*/
func (i Instance) getTravelTime(source int, destination int) float64 {
	return i.TravelTimes[source][destination]
}

// Converts JSON data to Patient objects. Returns an array of patients.
func (i Instance) getPatients() []Patient {
	var patientsSlice []Patient
	for id, patient := range i.Patients {
		int_id := strToInt(id)
		patient.ID = int_id
		patientsSlice = append(patientsSlice, patient)
	}
	return patientsSlice
}

// Returns patient with the specified id.
func (i Instance) getPatientAtID(patientID int) Patient {
    allPatients := i.getPatients()
    var resultPatient Patient
    for _, p := range allPatients {
        if p.ID == patientID {
            resultPatient = p
        }
    }
    return resultPatient
}


