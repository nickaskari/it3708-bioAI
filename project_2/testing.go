package main

import "fmt"


// Person struct - similar to a class in other languages
type Person struct {
    Name string
    Age  int
}

// Greet method on the Person struct
func (p Person) Greet() string {
    return "Hello, my name is " + p.Name
}

func main() {
    // Creating an instance of Person
    person := Person{Name: "John", Age: 30}

    // Calling a method on the Person instance
    greeting := person.Greet()
    fmt.Println(greeting)
}


func stuff() {
	i := 0
	// THIS IS A COMMENT

	// BASICALLY A WHILE LOOP
	for i < 3 {
		fmt.Println(i, "THATS CRAZY")
		i++
	}
}

func getCategoryCounts() map[string]int {
    // Initialize the map with some data
    counts := map[string]int{
        "books": 12,
        "videos": 4,
        "articles": 9,
    }
    return counts
}
