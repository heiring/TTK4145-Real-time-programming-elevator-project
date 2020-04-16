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

type StateTablesSync struct {
	sync.RWMutex
	Internal map[string][7][3]int
}

func (m *StateTablesSync) Read(key string) ([7][3]int, bool) {
	m.RLock()
	result, ok := m.Internal[key]
	m.RUnlock()
	return result, ok
}

func (m *StateTablesSync) Write(key string, value [7][3]int) {
	m.Lock()
	m.Internal[key] = value
	m.Unlock()
}

func (m *StateTablesSync) ReadWholeMap() map[string][7][3]int {
	m.RLock()
	result := m.Internal
	m.RUnlock()
	return result
}
