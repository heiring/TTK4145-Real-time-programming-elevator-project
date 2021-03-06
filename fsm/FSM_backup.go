package fsm

// import (
// 	"fmt"
// 	"math"
// 	"time"

// 	"../elevio"
// 	"../orderdistributor"
// 	"../statetable"
// 	"../tools"
// )

// type elev_state int

// const (
// 	INIT elev_state = iota
// 	IDLE
// 	MOVE
// 	WAIT
// 	EM_STOP
// )

// var currentOrder int

// func InitFSM(transmitStateTableCh chan<- [7][3]int) {
// 	// moveInDir(elevio.MD_Down)
// 	//!!
// 	//elevio.SetMotorDirection(elevio.MD_Down)
// 	//statetable.UpdateElevDirection(int(elevio.MD_Down))
// 	//!!
// 	go PollHardwareActions(transmitStateTableCh)
// }

// func setCurrentOrder(order int) {
// 	currentOrder = order
// }
// func getCurrentOrder() int {
// 	return currentOrder
// }

// func PollHardwareActions(stateTableTransmitCh chan<- [7][3]int) {
// 	drvButtons := make(chan elevio.ButtonEvent)
// 	drvFloors := make(chan int)
// 	drvObstr := make(chan bool)
// 	drvStop := make(chan bool)
// 	newOrder := make(chan int)

// 	go elevio.PollButtons(drvButtons)
// 	go elevio.PollFloorSensor(drvFloors)
// 	go elevio.PollObstructionSwitch(drvObstr)
// 	go elevio.PollStopButton(drvStop)
// 	go orderdistributor.PollOrders(newOrder)

// 	startWaitCh := make(chan bool)
// 	newMotorDirCh := make(chan elevio.MotorDirection)
// 	go executeNewMotorDirectionOrWait(startWaitCh, newMotorDirCh)
// 	//!!
// 	motorFunctionalCh := make(chan bool)
// 	orderReceivedCh := make(chan bool)
// 	go monitorMotorStatus(motorFunctionalCh, orderReceivedCh)
// 	//!!

// 	// var currentOrder int
// 	// localID := statetable.GetLocalID()
// 	for {
// 		select {
// 		case butn := <-drvButtons:
// 			handleButtonPressed(butn, stateTableTransmitCh)
// 		case floor := <-drvFloors:
// 			go handleNewFloor(floor, stateTableTransmitCh, startWaitCh, motorFunctionalCh, newMotorDirCh)
// 		case a := <-drvObstr:
// 			fmt.Println("case drvObstr: obstruction")
// 			fmt.Printf("%+v\n", a)
// 		case a := <-drvStop:
// 			fmt.Println("stop button pressed")
// 			fmt.Printf("%+v\n", a)
// 			moveInDir(elevio.MD_Stop, newMotorDirCh)
// 			stateTableTransmitCh <- statetable.Get()
// 		case order := <-newOrder:
// 			go handleNewOrder(order, stateTableTransmitCh, orderReceivedCh, startWaitCh, motorFunctionalCh, newMotorDirCh)
// 		default:
// 			// fmt.Println("Default yo")
// 		}

// 	}
// }

// func handleButtonPressed(butn elevio.ButtonEvent, stateTableTransmitCh chan<- [7][3]int) {
// 	localID := statetable.GetLocalID()
// 	var row int = 3 + butn.Floor
// 	var col int = int(butn.Button)
// 	elevio.SetButtonLamp(butn.Button, butn.Floor, true)
// 	statetable.UpdateActiveLights(butn.Button, butn.Floor, true)
// 	statetable.UpdateStateTableIndex(row, col, localID, 1, true)
// 	stateTableTransmitCh <- statetable.Get()
// }

// func handleNewOrder(order int, stateTableTransmitCh chan<- [7][3]int, orderReceivedCh, startWaitCh, motorFunctionalCh chan<- bool, newMotorDirCh chan<- elevio.MotorDirection) {
// 	orderReceivedCh <- true
// 	localID := statetable.GetLocalID()
// 	// currentOrder = order
// 	fmt.Println("NEW CURRENT ORDER: ", order)
// 	setCurrentOrder(order)            //!!!
// 	currentOrder := getCurrentOrder() //!!!
// 	if currentOrder == -1 {
// 		moveInDir(elevio.MD_Down, newMotorDirCh)
// 		stateTableTransmitCh <- statetable.Get()
// 	} else {
// 		currentFloor := statetable.GetCurrentFloor()
// 		currentDirection := statetable.GetElevDirection(localID)
// 		if currentOrder == currentFloor {
// 			moveInDir(elevio.MD_Stop, newMotorDirCh)
// 			completeCurOrder(startWaitCh, motorFunctionalCh)
// 			stateTableTransmitCh <- statetable.Get()
// 			// time.Sleep(3 * time.Second) //??????????????????????????????????????????
// 		} else if currentDirection == elevio.MD_Stop {
// 			newDirection, _ := tools.DivCheck((currentOrder - currentFloor), int(math.Abs(float64(currentOrder-currentFloor))))
// 			moveInDir(elevio.MotorDirection(newDirection), newMotorDirCh)
// 			stateTableTransmitCh <- statetable.Get()
// 		}
// 	}
// }

// func handleNewFloor(floor int, stateTableTransmitCh chan<- [7][3]int, startWaitCh, motorFunctionalCh chan<- bool, newMotorDirCh chan<- elevio.MotorDirection) {
// 	localID := statetable.GetLocalID()
// 	currentOrder := getCurrentOrder()
// 	lastFloor := statetable.GetCurrentFloor()
// 	curDir := statetable.GetElevDirection(localID)
// 	stateTableTransmitCh <- statetable.Get()
// 	if lastFloor != statetable.UnknownFloor {
// 		elevio.SetFloorIndicator(lastFloor)
// 	}
// 	elevio.SetFloorIndicator(floor)
// 	statetable.UpdateElevLastFLoor(floor)

// 	if currentOrder == floor {
// 		moveInDir(elevio.MD_Stop, newMotorDirCh)
// 		completeCurOrder(startWaitCh, motorFunctionalCh)
// 		stateTableTransmitCh <- statetable.Get()
// 		// time.Sleep(3 * time.Second)
// 	} else if currentOrder == -1 {
// 		moveInDir(elevio.MD_Stop, newMotorDirCh)
// 		stateTableTransmitCh <- statetable.Get()
// 		motorFunctionalCh <- true
// 	} else if currentOrder > floor {
// 		moveInDir(elevio.MD_Up, newMotorDirCh)
// 		stateTableTransmitCh <- statetable.Get()
// 	} else if currentOrder < floor {
// 		moveInDir(elevio.MD_Down, newMotorDirCh)
// 		stateTableTransmitCh <- statetable.Get()
// 	} else if (floor == 0 && curDir == int(elevio.MD_Down)) || (floor == 3 && curDir == int(elevio.MD_Up)) {
// 		moveInDir(elevio.MD_Stop, newMotorDirCh)
// 		stateTableTransmitCh <- statetable.Get()
// 	}
// }

// func moveInDir(dir elevio.MotorDirection, newMotorDirCh chan<- elevio.MotorDirection) {
// 	// elevio.SetMotorDirection(dir)
// 	statetable.UpdateElevDirection(int(dir))

// 	newMotorDirCh <- dir
// }

// func completeCurOrder(startWaitCh chan<- bool, motorFunctionalCh chan<- bool) {
// 	curFloor := statetable.GetCurrentFloor()
// 	row := curFloor + 3
// 	statetable.ResetRow(row)
// 	orderdistributor.CompleteCurrentOrder()
// 	for butn := elevio.BT_HallUp; butn <= elevio.BT_Cab; butn++ {
// 		elevio.SetButtonLamp(butn, curFloor, false)
// 	}
// 	startWaitCh <- true
// 	//!!
// 	motorFunctionalCh <- true

// }

// func executeNewMotorDirectionOrWait(startWaitCh <-chan bool, newMotorDirCh <-chan elevio.MotorDirection) {
// 	var motorDir elevio.MotorDirection
// 	motorDir = elevio.MD_Stop
// 	startedWaiting := time.Now()
// 	initCompleted := false
// 	newDir := true
// 	doorOpenLightLit := false
// 	for {
// 		select {
// 		case <-startWaitCh:
// 			startedWaiting = time.Now()
// 			elevio.SetDoorOpenLamp(true)
// 			doorOpenLightLit = true
// 		case newMotorDir := <-newMotorDirCh: //kan flere enn en direction vente i køen??
// 			if newMotorDir != motorDir {
// 				motorDir = newMotorDir
// 				newDir = true

// 			}
// 		default:
// 			if time.Now().Sub(startedWaiting) > 3*time.Second {
// 				if doorOpenLightLit {
// 					elevio.SetDoorOpenLamp(false)
// 					doorOpenLightLit = false
// 				}
// 				if newDir {
// 					elevio.SetMotorDirection(motorDir)
// 					statetable.UpdateElevDirection(int(motorDir))
// 					newDir = false
// 				}
// 			} else {
// 				if !initCompleted {
// 					elevio.SetMotorDirection(motorDir)
// 					statetable.UpdateElevDirection(int(motorDir))
// 					initCompleted = true
// 				}
// 			}
// 		}
// 	}
// }

// func monitorMotorStatus(motorFunctionalCh <-chan bool, orderRecievedCh <-chan bool) {
// 	motorFunctional := true
// 	orderCompleted := true
// 	lastOrderCompleted := time.Now()
// 	lastOrderReceived := time.Now()

// 	localID := statetable.GetLocalID()
// 	ticker := time.NewTicker(1000 * time.Millisecond)
// 	for {
// 		select {
// 		case <-motorFunctionalCh:
// 			motorFunctional = true
// 			orderCompleted = true
// 			lastOrderCompleted = time.Now()
// 		case <-orderRecievedCh: //kan en order komme før den forrige er fullført?
// 			orderCompleted = false
// 			lastOrderReceived = time.Now()
// 		case <-ticker.C:
// 			if time.Now().Sub(lastOrderCompleted) > 8000*time.Millisecond && time.Now().Sub(lastOrderReceived) > 8000*time.Millisecond && !orderCompleted {
// 				motorFunctional = false
// 			}
// 			stateTable := statetable.ReadStateTable(localID)
// 			if motorFunctional {
// 				if stateTable[0][2] == 0 {
// 					stateTable[0][2] = 1
// 					statetable.StateTables.Write(localID, stateTable)
// 					statetable.RunOrderDistribution()
// 				}
// 			} else {
// 				if stateTable[0][2] == 1 {
// 					stateTable[0][2] = 0
// 					statetable.StateTables.Write(localID, stateTable)
// 					statetable.RunOrderDistribution()
// 				}
// 			}
// 		default:
// 		}
// 	}
// }
