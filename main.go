package main

import (
	"os"

	"./elevio"
	"./network/bcast"
	"./network/network2"
	//"./fsm"
)

func main() {
	/*
		fsm.FSM(1)
		for true {

		}
	*/
	//initialization for simulator
	numFloors := 4
	ID := os.Args[1]
	elevio.Init("localhost:"+ID, numFloors)

	//network test
	elevatorStateTxCh := make(chan network2.ElevatorState)
	elevatorStateRxCh := make(chan network2.ElevatorState)

	transmitPacketCh := make(chan network2.ElevatorState)
	stateUpdateCh := make(chan network2.ElevatorState)

	lostIDCh := make(chan string)

	go network2.BroadcastElevatorState(transmitPacketCh, elevatorStateTxCh, 5000)
	go network2.ListenElevatorState(elevatorStateRxCh, stateUpdateCh, 10, lostIDCh, 4)

	go bcast.Transmitter(19569, elevatorStateTxCh)
	go bcast.Receiver(19569, elevatorStateRxCh)

	msg := network2.ElevatorState{ID: ID, IsAlive: true}

	for {
		transmitPacketCh <- msg
		select {
		case y := <-stateUpdateCh:
			//fmt.Println("main : packet received")
			//fmt.Println(y.ID)
			y.ID = "1111"
		default:
			//do stuff
		}
	}
}
