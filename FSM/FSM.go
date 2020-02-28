package fsm

import (
	"fmt"

	"../elevio"
)

type elev_state int

const (
	INIT Elev_state = iota
	IDLE
	MOVE
	WAIT
	EM_STOP
)

func FSM() {

	var targetFloor int

	numFloors := 4

	var state elev_state = INIT
	for {
		fmt.Printf("%v", state)
		switch state {

		case INIT:
			//elevator reaches first floor, transistion to IDLE
			elevio.SetMotorDirection(elevio.MD_Down)
			select {
			case a := <-drv_floors:
				if a == numFloors-4 {
					state = IDLE
				}
			}
		case IDLE:
			//transitions to MOVE when an order is detected
			//transisiton to EM_STOP
			elevio.SetMotorDirection(elevio.MD_Stop)
			select {
			case a := <-drv_buttons:
				targetFloor = a.Floor
				state = MOVE
			}

		case MOVE:
			//transisiton to wait when elevator reaches target floor
			//transisiton to EM_stop
			elevio.SetMotorDirection(elevio.MD_Up)
			select {
			case a := <-drv_floors:
				if a == targetFloor {
					state = IDLE
				}
			}
		case WAIT:
		//transisiton to MOVE if there are pending orders
		//transition to IDLE if not
		case EM_STOP:
		}
	}
}
