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

func InitFSM(transmitStateTableCh chan<- [7][3]int) {
	moveInDir(elevio.MD_Down)
	go pollHardwareActions(transmitStateTableCh)
}

func pollHardwareActions(stateTableTransmitCh chan<- [7][3]int) {
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
	localID := statetable.GetLocalID()
	fmt.Println("PRE FOR LOOP")

	for {
		select {
		case butn := <-drvButtons:
			var row int = 3 + butn.Floor
			var col int = int(butn.Button)
			elevio.SetButtonLamp(butn.Button, butn.Floor, true)
			statetable.UpdateStateTableIndex(row, col, localID, 1, true)
			stateTableTransmitCh <- statetable.Get()
		case floor := <-drvFloors:
			lastFloor := statetable.GetCurrentFloor()
			curDir := statetable.GetElevDirection(localID)
			stateTableTransmitCh <- statetable.Get()

			if lastFloor != statetable.UnknownFloor {
				elevio.SetFloorIndicator(lastFloor)
			}
			elevio.SetFloorIndicator(floor)
			statetable.UpdateElevLastFLoor(floor)

			if currentOrder == floor {
				moveInDir(elevio.MD_Stop)
				completeCurOrder()
				stateTableTransmitCh <- statetable.Get()
				time.Sleep(3 * time.Second)
			} else if currentOrder == -1 {
				moveInDir(elevio.MD_Stop)
				stateTableTransmitCh <- statetable.Get()
			} else if currentOrder > floor {
				moveInDir(elevio.MD_Up)
				stateTableTransmitCh <- statetable.Get()
			} else if currentOrder < floor {
				moveInDir(elevio.MD_Down)
				stateTableTransmitCh <- statetable.Get()
			} else if (floor == 0 && curDir == int(elevio.MD_Down)) || (floor == 3 && curDir == int(elevio.MD_Up)) {
				moveInDir(elevio.MD_Stop)
				stateTableTransmitCh <- statetable.Get()
			}

		case a := <-drvObstr:
			fmt.Printf("%+v\n", a)

		case a := <-drvStop:
			fmt.Printf("%+v\n", a)
			moveInDir(elevio.MD_Stop)
			stateTableTransmitCh <- statetable.Get()

		case order := <-newOrder:
			currentOrder = order
			if currentOrder == -1 {
				moveInDir(elevio.MD_Down)
				stateTableTransmitCh <- statetable.Get()
			} else {
				currentFloor := statetable.GetCurrentFloor()
				currentDirection := statetable.GetElevDirection(localID)
				if currentOrder == currentFloor {
					moveInDir(elevio.MD_Stop)
					completeCurOrder()
					stateTableTransmitCh <- statetable.Get()
					time.Sleep(3 * time.Second)
				} else if currentDirection == elevio.MD_Stop {
					newDirection, _ := tools.DivCheck((currentOrder - currentFloor), int(math.Abs(float64(currentOrder-currentFloor))))
					moveInDir(elevio.MotorDirection(newDirection))
					stateTableTransmitCh <- statetable.Get()
				}
			}
		}

	}
}

func moveInDir(dir elevio.MotorDirection) {
	elevio.SetMotorDirection(dir)
	statetable.UpdateElevDirection(int(dir))

}

func completeCurOrder() {
	curFloor := statetable.GetCurrentFloor()
	row := curFloor + 3
	statetable.ResetRow(row)
	orderdistributor.CompleteCurrentOrder()
	for butn := elevio.BT_HallUp; butn <= elevio.BT_Cab; butn++ {
		elevio.SetButtonLamp(butn, curFloor, false)
	}
}
