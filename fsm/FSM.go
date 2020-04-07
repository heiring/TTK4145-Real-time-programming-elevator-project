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

func FSM(elev_nr int, port string) {
	initFSM(elev_nr, port)
	go pollHardwareActions(elev_nr)
}

func initFSM(elev_nr int, port string) {
	numFloors := 4
	ip := "localhost:" + port
	elevio.Init(ip, numFloors)
	statetable.InitStateTable(elev_nr)
	moveInDir(elev_nr, elevio.MD_Down)
	fmt.Println("IP: ", ip)
}

func pollHardwareActions(elev_nr int) {
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
			statetable.UpdateStateTableIndex(row, col, elev_nr, 1)

		case floor := <-drvFloors:
			fmt.Println("SW: floor - ", floor)
			lastFloor := statetable.GetCurrentFloor(elev_nr)
			curDir := statetable.GetElevDirection(elev_nr)

			if lastFloor != statetable.UnknownFloor {
				elevio.SetFloorIndicator(lastFloor)
			}
			elevio.SetFloorIndicator(floor)
			statetable.UpdateElevLastFLoor(elev_nr, floor)

			if currentOrder == floor {
				fmt.Println("REACHED ORDER FLOOR: ", currentOrder)
				moveInDir(elev_nr, elevio.MD_Stop)
				completeCurOrder(elev_nr)
				// curOrder := orderdistributor.GetCurrentOrder()
				// fmt.Println("NEW ORDER: ", curOrder)
			} else if currentOrder == -1 {
				fmt.Println("NO NEW ORDERS, NOT MOVING")
				moveInDir(elev_nr, elevio.MD_Stop)
			} else if currentOrder > floor {
				fmt.Println("MOVING UP")
				moveInDir(elev_nr, elevio.MD_Up)
			} else if currentOrder < floor {
				fmt.Println("MOVING DOWN")
				moveInDir(elev_nr, elevio.MD_Down)
			} else if (floor == 0 && curDir == int(elevio.MD_Down)) || (floor == 3 && curDir == int(elevio.MD_Up)) {
				fmt.Println("MOVING STOPPED (FLOOR)")
				moveInDir(elev_nr, elevio.MD_Stop)
			}

		case a := <-drvObstr:
			fmt.Printf("%+v\n", a)

		case a := <-drvStop:
			// Stop the elevator
			fmt.Printf("%+v\n", a)
			moveInDir(elev_nr, elevio.MD_Stop)

		case order := <-newOrder:
			fmt.Println("SW - order: ", order)
			currentOrder = order
			if currentOrder == -1 {
				moveInDir(elev_nr, elevio.MD_Down)
			} else {
				currentFloor := statetable.GetCurrentFloor(elev_nr)
				currentDirection := statetable.GetElevDirection(elev_nr)
				if currentOrder == currentFloor {
					moveInDir(elev_nr, elevio.MD_Stop)
					completeCurOrder(elev_nr)
				} else if currentDirection == elevio.MD_Stop {
					newDirection, _ := tools.DivCheck((currentOrder - currentFloor), int(math.Abs(float64(currentOrder-currentFloor))))
					moveInDir(elev_nr, elevio.MotorDirection(newDirection))
				}
			}
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
