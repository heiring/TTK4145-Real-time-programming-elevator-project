package fsm

import (
	"fmt"

	"../orderdistributor"

	"../elevio"
	"../statetable"
)

type elev_state int

const (
	INIT elev_state = iota
	IDLE
	MOVE
	WAIT
	EM_STOP
)

func FSM(elev_nr int) {

	// var targetFloor int
	// numFloors := 4
	state := make(chan elev_state)

	go pollHardwareActions(state, elev_nr)

	for {
		fmt.Printf("%v", state)
		curFloor := statetable.GetCurrentFloor(elev_nr)
		switch <-state {

		case INIT:
			//elevator reaches first floor, transistion to IDLE
			elevio.SetMotorDirection(elevio.MD_Down)

			if curFloor == 0 {
				state <- IDLE
			}
		case IDLE:
			//transitions to MOVE when an order is detected
			//transisiton to EM_STOP
			elevio.SetMotorDirection(elevio.MD_Stop)
			// select {
			// case a := <-drv_buttons:
			// 	targetFloor = a.Floor
			// 	state = MOVE
			// }

		case MOVE:
			//transisiton to wait when elevator reaches target floor
			//transisiton to EM_stop
			elevio.SetMotorDirection(elevio.MD_Up)
			// select {
			// case a := <-drv_floors:
			// 	if a == targetFloor {
			// 		state = IDLE
			// 	}
			// }
		case WAIT:
		//transisiton to MOVE if there are pending orders
		//transition to IDLE if not
		case EM_STOP:
		}
	}
}

func pollHardwareActions(state chan elev_state, elev_nr int) {
	drvButtons := make(chan elevio.ButtonEvent)
	drvFloors := make(chan int)
	drvObstr := make(chan bool)
	drvStop := make(chan bool)

	go elevio.PollButtons(drvButtons)
	go elevio.PollFloorSensor(drvFloors)
	go elevio.PollObstructionSwitch(drvObstr)
	go elevio.PollStopButton(drvStop)

	// Detect state changes and update StateTable
	for {
		select {
		case a := <-drvButtons:
			fmt.Printf("%+b\n", a)
			var row int = 3 + a.Floor
			var col int = int(a.Button)
			elevio.SetButtonLamp(a.Button, a.Floor, true)
			statetable.UpdateStateTableIndex(row, col, elev_nr, 1)

		case a := <-drvFloors:
			fmt.Printf("%+v\n", a)
			elevio.SetFloorIndicator(a)

			curOrder := orderdistributor.GetCurrentOrder()

			if curOrder == a {
				state <- WAIT
				row := a + 3
				statetable.ResetRow(row)
				orderdistributor.CompleteCurrentOrder()
			}

			statetable.UpdateElevLastFLoor(elev_nr, a)

		case a := <-drvObstr:
			fmt.Printf("%+v\n", a)

		case a := <-drvStop:
			fmt.Printf("%+v\n", a)
		}
	}
}
