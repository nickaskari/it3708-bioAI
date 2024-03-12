package main

import (
	"sync"
)

type MigrationEvent struct {
    Migrants map[int][]Individual // key: generation, value: slice of individuals
    Counters map[int]int          // key: generation, value: count of islands deposited migrants
    TotalIslands int              // Total number of islands participating
	PickedUp map[int]map[int]bool
    sync.Mutex
    Ready *sync.Cond
}

func NewMigrationEvent(totalIslands int) *MigrationEvent {
    me := &MigrationEvent{
        Migrants: make(map[int][]Individual),
        Counters: make(map[int]int),
		PickedUp: make(map[int]map[int]bool),
        TotalIslands: totalIslands,
    }
    me.Ready = sync.NewCond(&me.Mutex)
    return me
}

func (me *MigrationEvent) DepositMigrants(generation int, migrants []Individual) {
    me.Lock()
    defer me.Unlock()
    me.Migrants[generation] = append(me.Migrants[generation], migrants...)
    me.Counters[generation]++
    if me.Counters[generation] == me.TotalIslands {
        me.Ready.Broadcast() // Signal that all islands have deposited for this generation
    }
}

func (me *MigrationEvent) WaitForMigration(generation, islandID int) []Individual {
    me.Lock()
    defer me.Unlock()
    for me.Counters[generation] < me.TotalIslands {
        me.Ready.Wait() // Wait until all islands have deposited their migrants
    }
    if me.PickedUp[generation][islandID] {
        // This island already picked up migrants for this generation, return an empty slice
        return []Individual{}
    }
    // Determine the subset of migrants this island should receive
    subsetSize := len(me.Migrants[generation]) / me.TotalIslands
    startIndex := subsetSize * (islandID % me.TotalIslands)
    endIndex := startIndex + subsetSize
    if endIndex > len(me.Migrants[generation]) { endIndex = len(me.Migrants[generation]) }
    
    // Extract the subset of migrants for this island
    migrants := make([]Individual, subsetSize)
    copy(migrants, me.Migrants[generation][startIndex:endIndex])
    
    // Mark as picked up for this island
    me.PickedUp[generation][islandID] = true
    
    // Optionally, remove migrants from the pool if you don't want them to be available for the next migration event
    return migrants
}