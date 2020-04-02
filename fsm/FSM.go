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
	go pollHardwareActions(elev_nr)
	initFSM(elev_nr)
}

func initFSM(elev_nr int) {
	numFloors := 4

	elevio.Init("localhost:15000", numFloors)
	statetable.InitStateTable(elev_nr)
	moveInDir(elev_nr, elevio.MD_Down)
}

func pollHardwareActions(elev_nr int) {
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
			fmt.Println("SW: button - ", a)
			// Update statetable
			// Update lights
			// If not moving, begin to move
			updateBtnLampAndStateTable(elev_nr, a)
			curDir := statetable.GetElevDirection(elev_nr)
			curFloor := statetable.GetCurrentFloor(elev_nr)
			fmt.Println("CURDIR - ", curDir)
			fmt.Println("CURFLOOR - ", curFloor)
			if curDir == elevio.MD_Stop {
				curOrder := orderdistributor.GetCurrentOrder()
				fmt.Println("CURORDER - ", curOrder)
				if curOrder > curFloor {
					fmt.Println("MOVING UP")
					moveInDir(elev_nr, elevio.MD_Up)
				} else if curOrder < curFloor {
					fmt.Println("MOVING DOWN")
					moveInDir(elev_nr, elevio.MD_Down)
				} else if curOrder == curFloor {
					fmt.Println("MOVING STOPPED")
					completeCurOrder(elev_nr, curFloor)

					curOrder := orderdistributor.GetCurrentOrder()
					fmt.Println("NEW ORDER: ", curOrder)

					if curOrder == -1 {
						fmt.Println("NO NEW ORDERS, NOT MOVING")
					} else if curOrder > curFloor {
						fmt.Println("MOVING UP")
						moveInDir(elev_nr, elevio.MD_Up)
					} else if curOrder < curFloor {
						fmt.Println("MOVING DOWN")
						moveInDir(elev_nr, elevio.MD_Down)
					}
				}
			}

		case floor := <-drvFloors:
			fmt.Println("SW: floor - ", floor)
			// Update statetable
			// Update lights
			// Check if orderfloor is reached
			lastFloor := statetable.GetCurrentFloor(elev_nr)
			curOrder := orderdistributor.GetCurrentOrder()
			curDir := statetable.GetElevDirection(elev_nr)
			if lastFloor != statetable.UnknownFloor {
				elevio.SetFloorIndicator(lastFloor)
			}
			elevio.SetFloorIndicator(floor)
			statetable.UpdateElevLastFLoor(elev_nr, floor)

			if curOrder == floor {
				fmt.Println("REACHED ORDER FLOOR: ", curOrder)
				completeCurOrder(elev_nr, floor)
				curOrder := orderdistributor.GetCurrentOrder()
				fmt.Println("NEW ORDER: ", curOrder)

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

			if (floor == 0 && curDir == int(elevio.MD_Down)) || (floor == 3 && curDir == int(elevio.MD_Up)) {
				fmt.Println("MOVING STOPPED (FLOOR)")
				moveInDir(elev_nr, elevio.MD_Stop)

				curOrder := orderdistributor.GetCurrentOrder()
				fmt.Println("NEW ORDER: ", curOrder)

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
	elevio.SetMotorDirection(dir)
	statetable.UpdateElevDirection(elev_nr, int(dir))
}

func completeCurOrder(elev_nr, curFloor int) {
	moveInDir(elev_nr, elevio.MD_Stop)
	time.Sleep(3 * time.Second)
	row := curFloor + 3
	statetable.ResetRow(row)
	orderdistributor.CompleteCurrentOrder()
	for butn := elevio.BT_HallUp; butn <= elevio.BT_Cab; butn++ {
		elevio.SetButtonLamp(butn, curFloor, false)
	}
}
