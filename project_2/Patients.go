package main


type Patient struct {
	ID        int `json:"-"`
	X_coord    int `json:"x_coord"`
	Y_coord    int `json:"y_coord"`
	Demand     int `json:"demand"`
	Start_time int `json:"start_time"`
	End_time   int `json:"end_time"`
	Care_time  int `json:"care_time"`
}

// Check if patient is in a list of patients.
func (p Patient) IsPatientInList(patients []Patient) bool {
	for _, patient := range patients {
		if patient.X_coord == p.X_coord && patient.Y_coord == p.Y_coord {
			return true
		}
	}
	return false
}

// Converts JSON data to Patient objects. Returns an array of patients.
func getPatients(instance Instance) []Patient {
	var patientsSlice []Patient
	for id, patient := range instance.Patients {
		int_id := strToInt(id)
		patient.ID = int_id
		patientsSlice = append(patientsSlice, patient)
	}

	return patientsSlice
}


