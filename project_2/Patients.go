package main

import (
	"fmt"
	"encoding/json"
)


// Really only need 'json:XXXX' if the JSON key is different from the variable name
type Patients struct {
	X_coord    int `json:"x_coord"`
    Y_coord    int `json:"y_coord"`
    Demand     int `json:"demand"`
    Start_time int `json:"start_time"`
    End_time   int `json:"end_time"`
    Care_time  int `json:"care_time"`
}


// for some reason, every time i print the patients, the order is different. But the individual values are correct.
func getPatients(filename string) ([]Patients, error) {
    data := readJSON(filename)
    if data == nil { 
        return nil, fmt.Errorf("failed to read JSON data from file: %s", filename)
    }

    var instance Instance
    err := json.Unmarshal(data, &instance)
    if err != nil {
        return nil, err
    }

    patientsSlice := make([]Patients, 0, len(instance.Patients))
    for _, patient := range instance.Patients {
        patientsSlice = append(patientsSlice, patient)
    }

    return patientsSlice, nil
}