package config

import "time"

const (
	TRANSMIT_INTERVAL    = 5000 * time.Millisecond
	ELEVATOR_TIMEOUT     = 10000 * time.Millisecond
	LAST_UPDATE_INTERVAL = 4000 * time.Millisecond
)
