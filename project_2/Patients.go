package main

import (
	"fmt"
	"math/rand"
	"time"
)

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

// Deletes patient from a list of patients. Returns the new list of patients.
func (p Patient) deletePatientFrom(patients []Patient) ([]Patient) {
	for i, patient := range patients {
		if patient.ID == p.ID {
			return append(patients[:i], patients[i+1:]...)
		}
	}
	fmt.Println("Unable to delete patient", p.ID, "from list of patients...")
	return patients 
}

func createDummyPatient() Patient {
	return Patient{-1, -1, -1, -1, -1, -1, -1, -1, -1}
}

// Generate a list of random patient ID's
func generateRandomPatientIDs(Instance Instance) []int {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	numbers := make([]int, Instance.NbrNurses)
    for i := range numbers {
        numbers[i] = i + 1
    }

    // Shuffle the array.
    random.Shuffle(len(numbers), func(i, j int) {
        numbers[i], numbers[j] = numbers[j], numbers[i]
    })

    // Decide the size of the subset. This example uses a random size between 1 and 25.
    subsetSize := rand.Intn(25) + 1

    // Slice the shuffled array to get a random subset of the desired size.
    return numbers[:subsetSize]
}