package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// The genetic algorithm. Returns best individual and whether benchmark was hit
func GA(populationSize int, gMax int, numParents int, temp int,
	crossoverRate float64, mutationRate float64, elitismPercentage float64, coolingRate float64,
	annealingRate float64, benchmark float64, ctx context.Context, migrationFrequency int, numMigrants int,
	migrationEvent *MigrationEvent, islandID int, initiateBestCostRepair bool, genocideWhenStuck int, instance Instance) (Individual, bool) {

	// initialize an emtpy array
	bestFitnesses := []float64{}

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	generation := 0
	var newIndividuals []Individual
	stuck := 0
	lastFitness := math.Inf(1)
	benchmarkWasReached := false
	detectedLongRoutes := false

	fmt.Println("Initalzing population..")
	population := initPopulation(instance, populationSize)
	printBestIndividual(population.Individuals, instance)

	for generation < gMax {
		newIndividuals = []Individual{}

		if generation%migrationFrequency == 0 && generation > 0 {
			migrants := population.selectRandomMigrants(numMigrants)
			migrationEvent.DepositMigrants(generation, migrants)

			// Wait for all islands to deposit migrants and pick up new ones
			newMigrants, didCancel := migrationEvent.WaitForMigration(generation, islandID, ctx)

			if didCancel {
				return getBestIndividual(population.Individuals), false
			}

			// Incorporate new migrants into the population
			population.insertNewMigrants(newMigrants)
		}

		// Age population
		population = population.agePopulation()

		parents := population.tournamentSelection(numParents)

		threshold := int(math.Floor(float64(population.size()) * 1.25))

		for threshold > len(newIndividuals) {
			// Check if other islands found the solution
			select {
			case <-ctx.Done():
				fmt.Println("GA was canceled.")
				return getBestIndividual(population.Individuals), false
			default:
				// Continue with GA processing
			}

			source := rand.NewSource(time.Now().UnixNano())
			r := rand.New(source)

			i, j := getTwoRandomParents(parents)
			//i, j := getTwoSimilarParents(parents)
			parent1, parent2 := parents[i], parents[j] // choose based on similarity
			if r.Float64() < crossoverRate {
				child1, child2 := destroyRepairCrossover(parent1, parent2, instance)
				child1.calculateFitness(instance)
				child2.calculateFitness(instance)

				if r.Float64() < mutationRate {
					var mutated1 Individual
					var mutated2 Individual

					if annealingRate > random.Float64() {
						mutated1 = simulatedAnnealing(child1, temp, coolingRate, instance)
						if initiateBestCostRepair {
							mutated1 = destroyRepairCluster(mutated1, instance)
						}
						mutated1.calculateFitness(instance)
					} else {
						mutated1 = hillClimbing(child1, temp, instance)
						if initiateBestCostRepair {
							mutated1 = destroyRepairCluster(mutated1, instance)
						}
						mutated1.calculateFitness(instance)
					}

					if annealingRate > random.Float64() {
						mutated2 = simulatedAnnealing(child2, temp, coolingRate, instance)
						if initiateBestCostRepair {
							mutated2 = destroyRepairCluster(mutated1, instance)
						}
						mutated2.calculateFitness(instance)
					} else {
						mutated2 = hillClimbing(child2, temp, instance)
						if initiateBestCostRepair {
							mutated2 = destroyRepairCluster(mutated1, instance)
						}
						mutated2.calculateFitness(instance)
					}

					newIndividuals = addToPopulation(mutated1, threshold, newIndividuals)
					newIndividuals = addToPopulation(mutated2, threshold, newIndividuals)
				} else {
					newIndividuals = addToPopulation(child1, threshold, newIndividuals)
					newIndividuals = addToPopulation(child2, threshold, newIndividuals)
				}
			} else {
				newIndividuals = addToPopulation(parent1, threshold, newIndividuals)
				newIndividuals = addToPopulation(parent2, threshold, newIndividuals)
			}
		}

		// Survivor selection -- AGE
		newIndividuals = ageSurvivorSelection(populationSize, newIndividuals)

		// Educate the elite
		newIndividuals = educateTheElite(elitismPercentage, newIndividuals, temp, coolingRate, initiateBestCostRepair, instance)

		// Survivor selection -- ELITISM

		newPopulation := deepCopyPopulation(population.applyElitismWithPercentage(newIndividuals, elitismPercentage))
		population = deepCopyPopulation(newPopulation)

		fmt.Println("GENEREATION", generation+1)
		population.printPopulationStats()
		printBestIndividual(population.Individuals, instance)

		bestFitness := getBestIndividual(population.Individuals).Fitness
		bestFitnesses = append(bestFitnesses, bestFitness)

		if bestFitness == lastFitness {
			stuck++
		} else {
			lastFitness = bestFitness
		}

		// 5 or 15
		if stuck%genocideWhenStuck == 0 && stuck > 0 {
			fmt.Println("\nPERFORM GENOCIDE AND REBUILD POPULATION..\n")
			newPopulation = deepCopyPopulation(population.applyGenecoideWithElitism(elitismPercentage, instance))
			population = newPopulation

			bestIndex := getBestIndividualIndex(population.Individuals)
			population.Individuals[bestIndex] = millitaryCamp(population.Individuals[bestIndex], instance)
		}

		if stuck%25 == 0 && stuck > 0 {
			genocideWhenStuck = 2
		}

		if bestFitness <= benchmark {
			fmt.Println("Found an individual with lower fitness than the benchmark..")
			benchmarkWasReached = true
			break
		}

		// Check if there are abnormally long routes.

		if !detectedLongRoutes {
			numNonEmptyRoutes := getBestIndividual(population.Individuals).numberOfNonEmptyRoutes()

			if numNonEmptyRoutes < 10 {
				crossoverRate = 0.2
				initiateBestCostRepair = false
				detectedLongRoutes = true
			}
		}

		if generation > 200 {
			detectedLongRoutes = true
		}

		generation++
	}

	getBestIndividual(population.Individuals).writeIndividualToJson()
	getBestIndividual(population.Individuals).writeIndividualToVismaFormat()
	writeBestFitnessesToJSON(bestFitnesses)

	return getBestIndividual(population.Individuals), benchmarkWasReached
}

// Get Two random indexes that are not the same
func getTwoRandomParents(parents []Individual) (int, int) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for {
		// indexes
		i := r.Intn(len(parents))
		j := r.Intn(len(parents))
		if i != j {
			return i, j
		}
	}
}

// Appends individual to array only if size resitriction is good. Returns array.
func addToPopulation(toAdd Individual, size int, individuals []Individual) []Individual {
	if len(individuals) < size {
		return append(individuals, toAdd)
	}
	return individuals
}
