package main

import (
	"./network/network2"

	"./network/bcast"
)

func main() {

	//initialization for simulator
	//numFloors := 3
	//elevio.Init("localhost:15657", numFloors)

	//network test
	elevatorStateTxCh := make(chan network2.ElevatorState)
	elevatorStateRxCh := make(chan network2.ElevatorState)

	transmitPacketCh := make(chan network2.ElevatorState)
	stateUpdateCh := make(chan network2.ElevatorState)

	lostIDCh := make(chan string)

	go network2.BroadcastElevatorState(transmitPacketCh, elevatorStateTxCh, 500)
	go network2.ListenElevatorState(elevatorStateRxCh, stateUpdateCh, 10000, lostIDCh)

	go bcast.Transmitter(10001, elevatorStateTxCh)
	go bcast.Receiver(100001, elevatorStateRxCh)

	finished := make(chan bool)
	finished <- false

	go network2.NetworkTest(transmitPacketCh, stateUpdateCh, finished)

	<-finished
}
