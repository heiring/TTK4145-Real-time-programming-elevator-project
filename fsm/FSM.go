package fsm

import (
	"fmt"
	"math"
	"time"

	"../elevio"
	"../orderdistributor"
	"../statetable"
	"../tools"
)

type elev_state int

const (
	INIT elev_state = iota
	IDLE
	MOVE
	WAIT
	EM_STOP
)

func InitFSM(elev_nr int, transmitStateTableCh chan<- [7][9]int) {
	moveInDir(elevio.MD_Down)
	go pollHardwareActions(elev_nr, transmitStateTableCh)
}

func pollHardwareActions(elev_nr int, stateTableTransmitCh chan<- [7][9]int) {
	drvButtons := make(chan elevio.ButtonEvent)
	drvFloors := make(chan int)
	drvObstr := make(chan bool)
	drvStop := make(chan bool)
	newOrder := make(chan int)

	go elevio.PollButtons(drvButtons)
	go elevio.PollFloorSensor(drvFloors)
	go elevio.PollObstructionSwitch(drvObstr)
	go elevio.PollStopButton(drvStop)
	go orderdistributor.PollOrders(newOrder)

	var currentOrder int

	for {
		select {
		case butn := <-drvButtons:
			fmt.Println("SW: button - ", butn)
			var row int = 3 + butn.Floor
			var col int = int(butn.Button)
			elevio.SetButtonLamp(butn.Button, butn.Floor, true)
			statetable.UpdateStateTableIndex(row, col, 1, true)
			stateTableTransmitCh <- statetable.Get()
		case floor := <-drvFloors:
			fmt.Println("SW: floor - ", floor)
			lastFloor := statetable.GetCurrentFloor(elev_nr)
			curDir := statetable.GetElevDirection(elev_nr)
			stateTableTransmitCh <- statetable.Get()

			if lastFloor != statetable.UnknownFloor {
				elevio.SetFloorIndicator(lastFloor)
			}
			elevio.SetFloorIndicator(floor)
			statetable.UpdateElevLastFLoor(floor)

			if currentOrder == floor {
				fmt.Println("REACHED ORDER FLOOR: ", currentOrder)
				moveInDir(elevio.MD_Stop)
				completeCurOrder(elev_nr)
				stateTableTransmitCh <- statetable.Get()
				// curOrder := orderdistributor.GetCurrentOrder()
				// fmt.Println("NEW ORDER: ", curOrder)
			} else if currentOrder == -1 {
				fmt.Println("NO NEW ORDERS, NOT MOVING")
				moveInDir(elevio.MD_Stop)
				stateTableTransmitCh <- statetable.Get()
			} else if currentOrder > floor {
				fmt.Println("MOVING UP")
				moveInDir(elevio.MD_Up)
				stateTableTransmitCh <- statetable.Get()
			} else if currentOrder < floor {
				fmt.Println("MOVING DOWN")
				moveInDir(elevio.MD_Down)
				stateTableTransmitCh <- statetable.Get()
			} else if (floor == 0 && curDir == int(elevio.MD_Down)) || (floor == 3 && curDir == int(elevio.MD_Up)) {
				fmt.Println("MOVING STOPPED (FLOOR)")
				moveInDir(elevio.MD_Stop)
				stateTableTransmitCh <- statetable.Get()
			}

		case a := <-drvObstr:
			fmt.Printf("%+v\n", a)

		case a := <-drvStop:
			// Stop the elevator
			fmt.Printf("%+v\n", a)
			moveInDir(elevio.MD_Stop)
			stateTableTransmitCh <- statetable.Get()

		case order := <-newOrder:
			fmt.Println("SW - order: ", order)
			currentOrder = order
			if currentOrder == -1 {
				moveInDir(elevio.MD_Down)
				stateTableTransmitCh <- statetable.Get()
			} else {
				currentFloor := statetable.GetCurrentFloor(elev_nr)
				currentDirection := statetable.GetElevDirection(elev_nr)
				if currentOrder == currentFloor {
					moveInDir(elevio.MD_Stop)
					completeCurOrder(elev_nr)
					stateTableTransmitCh <- statetable.Get()
				} else if currentDirection == elevio.MD_Stop {
					newDirection, _ := tools.DivCheck((currentOrder - currentFloor), int(math.Abs(float64(currentOrder-currentFloor))))
					moveInDir(elevio.MotorDirection(newDirection))
					stateTableTransmitCh <- statetable.Get()
				}
			}
		}

	}
}

// func updateBtnLampAndStateTable(butn elevio.ButtonEvent) {
// 	var row int = 3 + butn.Floor
// 	var col int = int(butn.Button)
// 	elevio.SetButtonLamp(butn.Button, butn.Floor, true)
// 	statetable.UpdateStateTableIndex(row, col, 1, true)
// }

func moveInDir(dir elevio.MotorDirection) {
	elevio.SetMotorDirection(dir)
	statetable.UpdateElevDirection(int(dir))

}

func completeCurOrder(elev_nr int) {
	time.Sleep(3 * time.Second)
	curFloor := statetable.GetCurrentFloor(elev_nr)
	row := curFloor + 3
	statetable.ResetRow(row)
	orderdistributor.CompleteCurrentOrder()
	for butn := elevio.BT_HallUp; butn <= elevio.BT_Cab; butn++ {
		elevio.SetButtonLamp(butn, curFloor, false)
	}
}
