package main

import (
    "fmt"
    "encoding/json"
    "os"
)

// Needed to read all the JSON file data
type Instance struct {
    InstanceName  string            `json:"instance_name"`
    NbrNurses     int               `json:"nbr_nurses"`
    CapacityNurse int               `json:"capacity_nurse"`
    Benchmark      float64           `json:"benchmark"`
    Depot          Depot             `json:"depot"`
    Patients       map[string]Patient `json:"patients"`
    TravelTimes   [][]float64      `json:"travel_times"`
    PatientArray    []Patient       `json:""`    
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
func (i *Instance) getPatients() []Patient {
	var patientsSlice []Patient
	for id, patient := range i.Patients {
		int_id := strToInt(id)
		patient.ID = int_id
		patientsSlice = append(patientsSlice, patient)
	}
    i.PatientArray = patientsSlice
	return patientsSlice
}

// Returns patient with the specified id.
func (i Instance) getPatientAtID(patientID int) Patient {
    allPatients := i.PatientArray
    var resultPatient Patient
    for _, p := range allPatients {
        if p.ID == patientID {
            resultPatient = p
        }
    }
    return resultPatient
}

// Reads problem instance from JSON and returns an Instance object.
func getProblemInstance(filename string) Instance {
	data, err := readJSON(filename)
	if err != nil {
		fmt.Println("failed to read JSON data: %w", err)
		os.Exit(1)
	}

	if data == nil {
		// Check for empty data
		fmt.Println("empty JSON data")
		os.Exit(1)
	}

	var instance Instance
	err = json.Unmarshal(data, &instance)
	if err != nil {
		fmt.Println("error unmarshaling JSON: %w", err)
		os.Exit(1)
	}
    instance.getPatients()
	return instance
}


