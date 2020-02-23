package FSM

import (
	"../elevio"
)

type Elev_state int

const (
	INIT Elev_state = iota
	IDLE
	MOVE
	WAIT
	EM_STOP
)

func FSM( state <- chan Elev_state) {
		for {
		switch <-state {
		case INIT:
			numFloors := 4
			elevio.Init("localhost:15659", numFloors)
			elevio.SetMotorDirection(elevio.MD_Down)
		case IDLE:

		case MOVE:

		case WAIT:

		case EM_STOP:
		}
	}
}
