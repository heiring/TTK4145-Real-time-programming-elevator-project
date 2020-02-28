package havai_main

import (
	"fmt"

	"./elevio"
)

type elev_state int

const (
	INIT elev_state = iota
	IDLE
	MOVE
	WAIT
	EM_STOP
)

func havai_main() {

	var targetFloor int

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

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
	/*
							numFloors := 4

							elevio.Init("localhost:15657", numFloors)

							var d elevio.MotorDirection = elevio.MD_Up
							//elevio.SetMotorDirection(d)

							drv_buttons := make(chan elevio.ButtonEvent)
							drv_floors := make(chan int)
							drv_obstr := make(chan bool)
							drv_stop := make(chan bool)

							go elevio.PollButtons(drv_buttons)
							go elevio.PollFloorSensor(drv_floors)
							go elevio.PollObstructionSwitch(drv_obstr)
							go elevio.PollStopButton(drv_stop)

							for {
								select {
								case a := <-drv_buttons:
									fmt.Printf("%+v\n", a)
									elevio.SetButtonLamp(a.Button, a.Floor, true)

								case a := <-drv_floors:
									fmt.Printf("%+v\n", a)
									if a == numFloors-1 {
										d = elevio.MD_Down
									} else if a == 0 {
										d = elevio.MD_Up
									}
									//elevio.SetMotorDirection(d)

								case a := <-drv_obstr:
									fmt.Printf("%+v\n", a)
									if a {
										elevio.SetMotorDirection(elevio.MD_Stop)
									} else {
										elevio.SetMotorDirection(d)
									}

								case a := <-drv_stop:
									fmt.Printf("%+v\n", a)
									for f := 0; f < numFloors; f++ {
										for b := elevio.ButtonType(0); b < 3; b++ {
									var targetFloor int

					numFloors := 4

					elevio.Init("localhost:15657", numFloors)

					drv_buttons := make(chan elevio.ButtonEvent)
					drv_floors := make(chan int)
					drv_obstr := make(chan bool)
					drv_stop := make(chan bool)

					go elevio.PollButtons(drv_buttons)
					go elevio.PollFloorSensor(drv_floors)
					go elevio.PollObstructionSwitch(drv_obstr)
					go elevio.PollStopButton(drv_stop)

					var state elev_state = INIT
					for {
						switch state {

						case INIT:
							//elevator reaches first floor, transistion to IDLE
							elevio.SetMotorDirection(elevio.MD_Down)
							select {
							case a := <-drv_floors:
								if a == numFloors-3 {
									elevio.SetMotorDirection(elevio.MD_Stop)
								}
							}
							state = IDLE
						case IDLE:
							//transitions to MOVE when an order is detected
							//transisiton to EM_STOP
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
									elevio.SetMotorDirection(elevio.MD_Stop)
								}
							}
						case WAIT:
						//transisiton to MOVE if there are pending orders
						//transition to IDLE if not
						case EM_STOP:
						}
		            }
	*/

}
