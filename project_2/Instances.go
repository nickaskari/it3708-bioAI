package main


// Needed to read all the JSON file data
type Instance struct {
    Instance_name  string            `json:"instance_name"`
    Nbr_nurses     int               `json:"nbr_nurses"`
    Capacity_nurse int               `json:"capacity_nurse"`
    Benchmark      float64           `json:"benchmark"`
    Depot          Depot             `json:"depot"`
    Patients       map[string]Patient `json:"patients"`
    Travel_times   [][]float64      `json:"travel_times"`
}


/* 
    Returns the traveltime between nurses and/or depot.
    Inputs are the row/column indicies of the nurses/depots in the
    travel matrix.
*/
func (i Instance) getTravelTime(source int, destination int) float64 {
	return i.Travel_times[source][destination]
}


