package main


// needed to read all the JSON file data
type Instance struct {
    Instance_name  string            `json:"instance_name"`
    Nbr_nurses     int               `json:"nbr_nurses"`
    Capacity_nurse int               `json:"capacity_nurse"`
    Benchmark      float64           `json:"benchmark"`
    Depot          Depot             `json:"depot"`
    Patients       map[string]Patient `json:"patients"`
    Travel_times   [][]float64      `json:"travel_times"`
}

func (i Instance) getTravelTime(location int, destination int) float64 {
	return i.Travel_times[location][destination]
}


