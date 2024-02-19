package main

import (
	"fmt"
	"os"
    "math/rand"
    "strconv"
)


func readJSON (filename string) []byte {
	data, err := os.ReadFile(filename)
    if err != nil {
      fmt.Print(err)
    }
	return data
}


type Randomizer struct {
    *rand.Rand
}

func (r *Randomizer) getRandomRoute(routes []Route) (Route, int) {

    randomIndex := r.Intn(len(routes))

    return routes[randomIndex], randomIndex
}


func strToInt(str string) int {

  n , err := strconv.Atoi(str)
  if err != nil{
    fmt.Println(err)
  }
  return n
}