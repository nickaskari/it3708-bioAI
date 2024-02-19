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


func (p Patient) IsPatientInList(patients []Patient) bool {
	for _, patient := range patients {
		if patient.X_coord == p.X_coord && patient.Y_coord == p.Y_coord {
			return true
		}
	}
	return false
}

// For some reason, every time i print the patients, the order is different. But the individual values are correct.
func getPatients(filename string, instance Instance, data []byte) []Patient {
	var patientsSlice []Patient
	for id, patient := range instance.Patients {

		int_id := strToInt(id)

		patient.ID = int_id

		patientsSlice = append(patientsSlice, patient)
	}

	return patientsSlice
}


