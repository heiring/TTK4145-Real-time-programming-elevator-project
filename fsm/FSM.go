package fsm

import (
	"fmt"
	"time"

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
	fmt.Println("FSM")
	initFSM(elev_nr)
	// var targetFloor int

	chState := make(chan elev_state)
	fmt.Println("state")

	fmt.Println("prePoll")
	go pollHardwareActions(chState, elev_nr)
	fmt.Println("preFor")

	// for {
	// 	state := <-chState
	// 	fmt.Printf("%v", state)
	// 	curFloor := statetable.GetCurrentFloor(elev_nr)
	// 	fmt.Println("preSwitch")

	// 	switch state {

	// 	case INIT:
	// 		fmt.Println("INIT")
	// 		//elevator reaches first floor, transistion to IDLE
	// 		elevio.SetMotorDirection(elevio.MD_Down)
	// 		statetable.UpdateElevDirection(elev_nr, elevio.MD_Down)
	// 	case IDLE:
	// 		fmt.Println("IDLE")
	// 		elevio.SetMotorDirection(elevio.MD_Stop)
	// 		statetable.UpdateElevDirection(elev_nr, elevio.MD_Stop)

	// 	case MOVE:
	// 		fmt.Println("MOVE")
	// 		curOrder := orderdistributor.GetCurrentOrder()
	// 		if curOrder > curFloor {
	// 			fmt.Println("MOVING UP")
	// 			elevio.SetMotorDirection(elevio.MD_Up)
	// 			statetable.UpdateElevDirection(elev_nr, elevio.MD_Down)
	// 		} else if curOrder < curFloor {
	// 			fmt.Println("MOVING DOWN")
	// 			elevio.SetMotorDirection(elevio.MD_Down)
	// 			statetable.UpdateElevDirection(elev_nr, elevio.MD_Down)
	// 		} else if curOrder == curFloor {
	// 			fmt.Println("MOVING STOPPED")
	// 			// state <- WAIT
	// 			row := curFloor + 3
	// 			statetable.ResetRow(row)
	// 			orderdistributor.CompleteCurrentOrder()
	// 		}
	// 	case WAIT:
	// 		fmt.Println("WAIT")
	// 		elevio.SetMotorDirection(elevio.MD_Stop)
	// 		statetable.UpdateElevDirection(elev_nr, elevio.MD_Stop)
	// 	//transisiton to MOVE if there are pending orders
	// 	//transition to IDLE if not
	// 	case EM_STOP:
	// 		fmt.Println("EM_STOP")
	// 	}
	// }
}

func initFSM(elev_nr int) {
	numFloors := 4

	elevio.Init("localhost:15657", numFloors)
	fmt.Println("elevioinit")
	statetable.InitStateTable()

	moveInDir(elev_nr, elevio.MD_Down)
	for true {
		if statetable.GetCurrentFloor(elev_nr) == 0 {
			moveInDir(elev_nr, elevio.MD_Stop)
			break
		}
	}
}

func pollHardwareActions(state chan elev_state, elev_nr int) {
	// Todo:
	// *

	state <- INIT
	drvButtons := make(chan elevio.ButtonEvent)
	drvFloors := make(chan int)
	drvObstr := make(chan bool)
	drvStop := make(chan bool)

	go elevio.PollButtons(drvButtons)
	go elevio.PollFloorSensor(drvFloors)
	go elevio.PollObstructionSwitch(drvObstr)
	go elevio.PollStopButton(drvStop)
	// Detect hw changes and update StateTable
	fmt.Println("HWPreLoop")
	for {
		fmt.Println("HWLoop")

		select {
		case a := <-drvButtons:
			// Update statetable
			// Update lights
			// If not moving, begin to move
			updateBtnLampAndStateTable(elev_nr, a)
			curDir := statetable.GetElevDirection(elev_nr)
			curFloor := statetable.GetCurrentFloor(elev_nr)
			if curDir == elevio.MD_Stop {
				curOrder := orderdistributor.GetCurrentOrder()
				if curOrder > curFloor {
					fmt.Println("MOVING UP")
					moveInDir(elev_nr, elevio.MD_Up)
				} else if curOrder < curFloor {
					fmt.Println("MOVING DOWN")
					moveInDir(elev_nr, elevio.MD_Down)
				} else if curOrder == curFloor {
					moveInDir(elev_nr, elevio.MD_Stop)
					completeCurOrder(elev_nr, curFloor)
				}
			}
			fmt.Printf("%+b\n", a)

		case floor := <-drvFloors:
			// Update statetable
			// Update lights
			// Check if orderfloor is reached
			lastFloor := statetable.GetCurrentFloor(elev_nr)
			curOrder := orderdistributor.GetCurrentOrder()
			elevio.SetFloorIndicator(lastFloor)
			elevio.SetFloorIndicator(floor)
			statetable.UpdateElevLastFLoor(elev_nr, floor)
			fmt.Println("Current floor: ", floor)

			if curOrder == floor {
				completeCurOrder(elev_nr, floor)
				curOrder := orderdistributor.GetCurrentOrder()
				time.Sleep(3 * time.Second)

				if curOrder == -1 {
					fmt.Println("NO NEW ORDERS, NOT MOVING")
				} else if curOrder > floor {
					fmt.Println("MOVING UP")
					moveInDir(elev_nr, elevio.MD_Up)
				} else if curOrder < floor {
					fmt.Println("MOVING DOWN")
					moveInDir(elev_nr, elevio.MD_Down)
				}
			}

		case a := <-drvObstr:
			fmt.Printf("%+v\n", a)

		case a := <-drvStop:
			// Stop the elevator
			fmt.Printf("%+v\n", a)
			moveInDir(elev_nr, elevio.MD_Stop)
		}
	}
}

func updateBtnLampAndStateTable(elev_nr int, butn elevio.ButtonEvent) {
	var row int = 3 + butn.Floor
	var col int = int(butn.Button)
	elevio.SetButtonLamp(butn.Button, butn.Floor, true)
	statetable.UpdateStateTableIndex(row, col, elev_nr, 1)
}

func moveInDir(elev_nr int, dir elevio.MotorDirection) {
	fmt.Println("UP")
	elevio.SetMotorDirection(dir)
	statetable.UpdateElevDirection(elev_nr, int(dir))
}

func completeCurOrder(elev_nr, curFloor int) {
	moveInDir(elev_nr, elevio.MD_Stop)
	row := curFloor + 3
	statetable.ResetRow(row)
	orderdistributor.CompleteCurrentOrder()
	for butn := elevio.BT_HallUp; butn < elevio.BT_Cab; butn++ {
		elevio.SetButtonLamp(butn, curFloor, false)
	}
}
