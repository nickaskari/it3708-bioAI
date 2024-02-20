package main


type Patient struct {
	ID        int `json:"-"`
	XCoord    int `json:"x_coord"`
	YCoord    int `json:"y_coord"`
	Demand     int `json:"demand"`
	StartTime int `json:"start_time"`
	EndTime   int `json:"end_time"`
	CareTime  int `json:"care_time"`
	VisitTime float64 `json:"-"`
	LeavingTime float64 `json:"-"`
}

// Check if patient is in a list of patients.
func (p Patient) IsPatientInList(patients []Patient) bool {
	for _, patient := range patients {
		if patient.XCoord == p.XCoord && patient.YCoord == p.YCoord {
			return true
		}
	}
	return false
}


