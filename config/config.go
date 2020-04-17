package config

import (
	"sync"
	"time"
)

const (
	StateTransmissionInterval           = 500 * time.Millisecond
	ElevatorTimeout                     = 10000 * time.Millisecond
	LastUpdateInterval                  = 2000 * time.Millisecond
	ActiveElevatorsTransmissionInterval = 1000 * time.Millisecond
	TransmissionPort                    = 19569
	PollRate                            = 20 * time.Millisecond
)

type ElevatorState struct {
	ID         string
	StateTable [7][3]int
}

type StateTablesMutex struct {
	sync.RWMutex
	Internal map[string][7][3]int
}

type ActiveLightsMutex struct {
	sync.RWMutex
	Internal map[[2]int]bool
}

func (m *StateTablesMutex) Read(key string) ([7][3]int, bool) {
	m.RLock()
	result, ok := m.Internal[key]
	m.RUnlock()
	return result, ok
}

func (m *StateTablesMutex) Write(key string, value [7][3]int) {
	m.Lock()
	m.Internal[key] = value
	m.Unlock()
}

func (m *StateTablesMutex) ReadWholeMap() map[string][7][3]int {
	m.RLock()
	result := m.Internal
	m.RUnlock()
	return result
}

func (m *StateTablesMutex) WriteWholeMap(stateTables map[string][7][3]int) {
	m.Lock()
	m.Internal = stateTables
	m.Unlock()
}

func (m *ActiveLightsMutex) Read(butn int, floor int) (bool, bool) {
	m.RLock()
	result, ok := m.Internal[[2]int{butn, floor}]
	m.RUnlock()
	return result, ok
}

func (m *ActiveLightsMutex) Write(butn int, floor int, value bool) {
	m.Lock()
	m.Internal[[2]int{butn, floor}] = value
	m.Unlock()
}

func (m *ActiveLightsMutex) ReadWholeMap() map[[2]int]bool {
	m.RLock()
	result := m.Internal
	m.RUnlock()
	return result
}

func (m *ActiveLightsMutex) WriteWholeMap(activeLights map[[2]int]bool) {
	m.Lock()
	m.Internal = activeLights
	m.Unlock()
}
