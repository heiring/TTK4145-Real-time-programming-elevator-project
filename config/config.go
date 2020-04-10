package config

import "time"

const (
	StateTransmissionInterval           = 1000 * time.Millisecond
	ElevatorTimeout                     = 10000 * time.Millisecond
	LastUpdateInterval                  = 2000 * time.Millisecond
	ActiveElevatorsTransmissionInterval = 4000 * time.Millisecond
	TransmissionPort                    = 19569
	PollRate                            = 20 * time.Millisecond
)

type ElevatorState struct {
	ID         string
	StateTable [7][9]int
}
