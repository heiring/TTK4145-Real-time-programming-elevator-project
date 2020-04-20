package mapsync

import "sync"

// This module enables safe concurrent use of the maps StateTables and ActiveLights.

type StateTablesSync struct {
	sync.RWMutex
	StateTables map[string][7][3]int
}

type ActiveLightsSync struct {
	sync.RWMutex
	ActiveLights map[[2]int]bool
}

func (mapSync *StateTablesSync) Read(key string) ([7][3]int, bool) {
	mapSync.RLock()
	stateTable, ok := mapSync.StateTables[key]
	mapSync.RUnlock()
	return stateTable, ok
}

func (mapSync *StateTablesSync) Write(key string, value [7][3]int) {
	mapSync.Lock()
	mapSync.StateTables[key] = value
	mapSync.Unlock()
}

func (mapSync *StateTablesSync) ReadWholeMap() map[string][7][3]int {
	mapSync.RLock()
	stateTables := mapSync.StateTables
	mapSync.RUnlock()
	return stateTables
}

func (mapSync *ActiveLightsSync) Write(butn int, floor int, value bool) {
	mapSync.Lock()
	mapSync.ActiveLights[[2]int{butn, floor}] = value
	mapSync.Unlock()
}

func (mapSync *ActiveLightsSync) ReadWholeMap() map[[2]int]bool {
	mapSync.RLock()
	activeLights := mapSync.ActiveLights
	mapSync.RUnlock()
	return activeLights
}

func (mapSync *ActiveLightsSync) WriteWholeMap(activeLights map[[2]int]bool) {
	mapSync.Lock()
	mapSync.ActiveLights = activeLights
	mapSync.Unlock()
}
