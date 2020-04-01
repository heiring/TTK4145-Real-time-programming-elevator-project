package main

import (
	"fmt"

	"./network/network2"
)

func main() {

	//initialization for simulator
	//numFloors := 3
	//elevio.Init("localhost:15657", numFloors)

	//network test
	elevatorStateTxCh := make(chan network2.ElevatorState)
	elevatorStateRxCh := make(chan network2.ElevatorState)

	transmitPacketCh := make(chan network2.ElevatorState)
	//stateUpdateCh := make(chan network2.ElevatorState)

	//lostIDCh := make(chan string)

	go network2.BroadcastElevatorState(transmitPacketCh, elevatorStateTxCh, 500)
	//go network2.ListenElevatorState(elevatorStateRxCh, stateUpdateCh, 10000, lostIDCh)

	// go bcast.Transmitter(19569, elevatorStateTxCh)
	// go bcast.Receiver(19569, elevatorStateRxCh)

	//go network2.NetworkTest(transmitPacketCh, stateUpdateCh)

	yeet := network2.ElevatorState{ID: "2222", IsAlive: true}

	for {
		//elevatorStateTxCh <- yeet
		transmitPacketCh <- yeet
		select {
		case y := <-elevatorStateRxCh:
			fmt.Println("packet received")
			fmt.Println(y.ID)
		default:
			//do stuff
		}
	}
}
