package main

import (
	"fmt"
	"os"

	"./config"
	"./elevio"
	"./network/bcast"
	"./network/network2"
	//"flag"
	//"./fsm"
)

func main() {
	//var id int
	//var port string
	//lag.IntVar(&id, "id", 1, "Specify the id nr for the elevator")
	//flag.StringVar(&port, "port", "32001", "Specify a port corresponding to an elevator")
	//flag.Parse()

	//fsm.FSM(id, port)
	//for true {

	//}
	// 	//initialization for simulator
	numFloors := 4
	ID := os.Args[1]
	elevio.Init("localhost:"+ID, numFloors)

	//network test
	elevatorStateTxCh := make(chan network2.ElevatorState)
	elevatorStateRxCh := make(chan network2.ElevatorState)

	transmitPacketCh := make(chan network2.ElevatorState)
	stateUpdateCh := make(chan network2.ElevatorState)

	lostIDCh := make(chan string)
	lifeSignalIDCh := make(chan string)

	activeElevatorsCh := make(chan map[string]bool)

	go network2.BroadcastElevatorState(transmitPacketCh, elevatorStateTxCh, config.TRANSMIT_INTERVAL)
	go network2.ListenElevatorState(elevatorStateRxCh, stateUpdateCh, lostIDCh, lifeSignalIDCh, config.ELEVATOR_TIMEOUT, config.LAST_UPDATE_INTERVAL)

	go bcast.Transmitter(19569, elevatorStateTxCh)
	go bcast.Receiver(19569, elevatorStateRxCh)

	go network2.MonitorActiveElevators(lostIDCh, lifeSignalIDCh, activeElevatorsCh)

	msg := network2.ElevatorState{ID: ID}

	for {
		transmitPacketCh <- msg
		select {
		case y := <-stateUpdateCh:
			//fmt.Println("main : packet received")
			//fmt.Println(y.ID)
			y.ID = "1111"
		case activeElevators := <-activeElevatorsCh:
			for ID, isAlive := range activeElevators {
				fmt.Println(ID)
				fmt.Println(isAlive)
			}
		default:
			//do stuff
		}
	}
}
