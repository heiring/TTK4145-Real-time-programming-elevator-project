package config

import "time"

const (
	TRANSMIT_INTERVAL    = 1000 * time.Millisecond
	ELEVATOR_TIMEOUT     = 10000 * time.Millisecond
	LAST_UPDATE_INTERVAL = 2000 * time.Millisecond
	TRANSMIT_PORT        = 19569
	PollRate             = 20 * time.Millisecond
)

type orderDistributionChannels struct {
	stateUpdateCh    chan ElevatorState
	transmitPacketCh chan ElevatorState
	activeElevators  chan map[string]bool
}

type ElevatorState struct {
	ID         string
	StateTable [7][9]int
}
