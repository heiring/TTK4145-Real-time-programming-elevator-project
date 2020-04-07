package main

import (
	// "os"

	// "./elevio"
	// "./network/bcast"
	// "./network/network2"
	"flag"

	"./fsm"
)

func main() {
	var id int
	var port string
	flag.IntVar(&id, "id", 1, "Specify the id nr for the elevator")
	flag.StringVar(&port, "port", "32001", "Specify a port corresponding to an elevator")
	flag.Parse()

	fsm.FSM(id, port)
	for true {

	}
	//initialization for simulator
	// numFloors := 3
	// ID := os.Args[1]
	// elevio.Init("localhost:"+ID, numFloors)

	// //network test
	// elevatorStateTxCh := make(chan network2.ElevatorState)
	// elevatorStateRxCh := make(chan network2.ElevatorState)

	// transmitPacketCh := make(chan network2.ElevatorState)
	// stateUpdateCh := make(chan network2.ElevatorState)

	// lostIDCh := make(chan string)

	// go network2.BroadcastElevatorState(transmitPacketCh, elevatorStateTxCh, 5000)
	// go network2.ListenElevatorState(elevatorStateRxCh, stateUpdateCh, 10, lostIDCh, 4)

	// go bcast.Transmitter(19569, elevatorStateTxCh)
	// go bcast.Receiver(19569, elevatorStateRxCh)

	// msg := network2.ElevatorState{ID: ID, IsAlive: true}

	// for {
	// 	transmitPacketCh <- msg
	// 	select {
	// 	case y := <-stateUpdateCh:
	// 		//fmt.Println("main : packet received")
	// 		//fmt.Println(y.ID)
	// 		y.ID = "1111"
	// 	default:
	// 		//do stuff
	// 	}
	// }
}
