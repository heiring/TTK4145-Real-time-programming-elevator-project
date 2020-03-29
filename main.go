package main

import (
	"./network"
)

func main() {

	//initialization for simulator
	//numFloors := 3
	//elevio.Init("localhost:15657", numFloors)

	//network test
	elevatorStateTxCh := make(chan network.ElevatorState)
	elevatorStateRxCh := make(chan network.ElevatorState)

	transmitPacketCh := make(chan network.ElevatorState)
	stateUpdateCh := make(chan network.ElevatorState)

	lostIDCh := make(chan string)

	go network.BroadcastElevatorState(transmitPacketCh, elevatorStateTxCh, 500)
	go network.ListenElevatorState(elevatorStateRxCh, stateUpdateCh, 10000, lostIDCh)

	go bcast.Transmitter(10001, elevatorStateTxCh)
	go bcast.Receiver(100001, elevatorStateRxCh)

}
