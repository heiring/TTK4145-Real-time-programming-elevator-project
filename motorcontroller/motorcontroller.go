package motorcontroller

import (
	"fmt"
	"math"
	"time"

	"../elevio"
	"../orderdistributor"
	"../statetable"
	"../tools"
)

const InitOrder int = -1

// MotorController updates the state table when buttons are pressed, floor sensors are triggered or when the motor direction changes.
// It receives orders from the order distributor and sets the motor direction accordingly.
// When the state table is changed, it signals that a new state table is to be transmitted.
func MotorController(stateTableTransmitCh chan<- [7][3]int) {
	drvButtonsCh := make(chan elevio.ButtonEvent)
	drvFloorCh := make(chan int)
	drvObstrCh := make(chan bool)
	drvStopCh := make(chan bool)
	newOrderCh := make(chan int)
	startWaitCh := make(chan bool)
	newMotorDirCh := make(chan elevio.MotorDirection)
	motorFunctionalCh := make(chan bool)
	orderReceivedCh := make(chan bool)

	go elevio.PollButtons(drvButtonsCh)
	go elevio.PollFloorSensor(drvFloorCh)
	go elevio.PollObstructionSwitch(drvObstrCh)
	go elevio.PollStopButton(drvStopCh)
	go orderdistributor.PollOrders(newOrderCh)
	go executeNewMotorDirectionOrWait(startWaitCh, newMotorDirCh)
	go monitorMotorStatus(motorFunctionalCh, orderReceivedCh)

	var currentOrder int
	localID := statetable.GetLocalID()
	for {
		select {
		case butn := <-drvButtonsCh:
			var row int = 3 + butn.Floor
			var col int = int(butn.Button)
			elevio.SetButtonLamp(butn.Button, butn.Floor, true)
			statetable.UpdateActiveLights(butn.Button, butn.Floor, true)
			statetable.UpdateStateTableIndex(row, col, localID, 1, true)
			stateTableTransmitCh <- statetable.Get()
		case floor := <-drvFloorCh:
			lastFloor := statetable.GetCurrentFloor()
			curDir := statetable.GetElevDirection(localID)
			stateTableTransmitCh <- statetable.Get()
			if lastFloor != statetable.UnknownFloor {
				elevio.SetFloorIndicator(lastFloor)
			}
			elevio.SetFloorIndicator(floor)
			statetable.UpdateElevLastFLoor(floor)

			if currentOrder == floor {
				moveInDir(elevio.MD_Stop, newMotorDirCh)
				completeCurOrder(startWaitCh, motorFunctionalCh)
				stateTableTransmitCh <- statetable.Get()
			} else if currentOrder == InitOrder {
				moveInDir(elevio.MD_Stop, newMotorDirCh)
				stateTableTransmitCh <- statetable.Get()
				motorFunctionalCh <- true
			} else if currentOrder > floor {
				moveInDir(elevio.MD_Up, newMotorDirCh)
				stateTableTransmitCh <- statetable.Get()
			} else if currentOrder < floor {
				moveInDir(elevio.MD_Down, newMotorDirCh)
				stateTableTransmitCh <- statetable.Get()
			} else if (floor == 0 && curDir == int(elevio.MD_Down)) || (floor == 3 && curDir == int(elevio.MD_Up)) {
				moveInDir(elevio.MD_Stop, newMotorDirCh)
				stateTableTransmitCh <- statetable.Get()
			}

		case obstrEvent := <-drvObstrCh:
			fmt.Println("case drvObstr: obstruction")
			fmt.Printf("%+v\n", obstrEvent)
		case stopEvent := <-drvStopCh:
			fmt.Println("stop button pressed")
			fmt.Printf("%+v\n", stopEvent)
			moveInDir(elevio.MD_Stop, newMotorDirCh)
			stateTableTransmitCh <- statetable.Get()
		case order := <-newOrderCh:
			orderReceivedCh <- true
			currentOrder = order
			if currentOrder == InitOrder {
				moveInDir(elevio.MD_Down, newMotorDirCh)
				stateTableTransmitCh <- statetable.Get()
			} else {
				currentFloor := statetable.GetCurrentFloor()
				currentDirection := statetable.GetElevDirection(localID)
				if currentOrder == currentFloor {
					moveInDir(elevio.MD_Stop, newMotorDirCh)
					completeCurOrder(startWaitCh, motorFunctionalCh)
					stateTableTransmitCh <- statetable.Get()
				} else if currentDirection == elevio.MD_Stop {
					newDirection, _ := tools.DivCheck((currentOrder - currentFloor), int(math.Abs(float64(currentOrder-currentFloor))))
					moveInDir(elevio.MotorDirection(newDirection), newMotorDirCh)
					stateTableTransmitCh <- statetable.Get()
				}
			}
		default:
		}

	}
}

func moveInDir(dir elevio.MotorDirection, newMotorDirCh chan<- elevio.MotorDirection) {
	statetable.UpdateElevDirection(int(dir))
	newMotorDirCh <- dir
}

func completeCurOrder(startWaitCh chan<- bool, motorFunctionalCh chan<- bool) {
	curFloor := statetable.GetCurrentFloor()
	row := curFloor + 3
	statetable.ResetRow(row)
	orderdistributor.RemoveOrder()
	for butn := elevio.BT_HallUp; butn <= elevio.BT_Cab; butn++ {
		elevio.SetButtonLamp(butn, curFloor, false)
	}
	startWaitCh <- true
	motorFunctionalCh <- true
}

// executeNewMotorDirectionOrWait gets a signal when an order is completed. Then, the door open lamp is lit for three seconds.
// After three seconds the elevator may start moving in a potential new direction received on the newMotorDirCh-channel.
func executeNewMotorDirectionOrWait(startWaitCh <-chan bool, newMotorDirCh <-chan elevio.MotorDirection) {
	var motorDir elevio.MotorDirection
	motorDir = elevio.MD_Stop
	startedWaiting := time.Now()
	initCompleted := false
	newDir := true
	doorOpenLightLit := false
	for {
		select {
		case <-startWaitCh:
			startedWaiting = time.Now()
			elevio.SetDoorOpenLamp(true)
			doorOpenLightLit = true
		case newMotorDir := <-newMotorDirCh:
			if newMotorDir != motorDir {
				motorDir = newMotorDir
				newDir = true

			}
		default:
			if time.Now().Sub(startedWaiting) > 3*time.Second {
				if doorOpenLightLit {
					elevio.SetDoorOpenLamp(false)
					doorOpenLightLit = false
				}
				if newDir {
					elevio.SetMotorDirection(motorDir)
					statetable.UpdateElevDirection(int(motorDir))
					newDir = false
				}
			} else {
				if !initCompleted {
					elevio.SetMotorDirection(motorDir)
					statetable.UpdateElevDirection(int(motorDir))
					initCompleted = true
				}
			}
		}
	}
}

// monitorMotorStatus updates StateTables if the motor in the local elevator stops functioning, or if it starts functioning after being non-functional.
// The elevator is given eight seconds to complete an order. If it fails to do so, the motor is labeled as non-functional in the state table.
// When an order is completed, a signal is received, which indicates that the motor is functioning.
func monitorMotorStatus(motorFunctionalCh <-chan bool, orderRecievedCh <-chan bool) {
	motorFunctional := true
	orderCompleted := true
	lastOrderCompleted := time.Now()
	lastOrderReceived := time.Now()

	localID := statetable.GetLocalID()
	ticker := time.NewTicker(1000 * time.Millisecond)
	for {
		select {
		case <-motorFunctionalCh:
			motorFunctional = true
			orderCompleted = true
			lastOrderCompleted = time.Now()
		case <-orderRecievedCh:
			orderCompleted = false
			lastOrderReceived = time.Now()
		case <-ticker.C:
			if time.Now().Sub(lastOrderCompleted) > 8000*time.Millisecond && time.Now().Sub(lastOrderReceived) > 8000*time.Millisecond && !orderCompleted {
				motorFunctional = false
			}
			stateTable := statetable.ReadStateTable(localID)
			if motorFunctional {
				if stateTable[0][2] == 0 {
					stateTable[0][2] = 1
					statetable.StateTables.Write(localID, stateTable)
					statetable.RunOrderDistribution()
				}
			} else {
				if stateTable[0][2] == 1 {
					stateTable[0][2] = 0
					statetable.StateTables.Write(localID, stateTable)
					statetable.RunOrderDistribution()
				}
			}
		default:
		}
	}
}
