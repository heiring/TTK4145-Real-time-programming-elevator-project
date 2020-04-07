package main

import (
	"fmt"
	"os"

	"./config"
	"./elevio"

	//"./network/bcast"
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

	transmitPacketCh := make(chan network2.ElevatorState)
	stateUpdateCh := make(chan network2.ElevatorState)
	activeElevatorsCh := make(chan map[string]bool)

	go network2.RunNetwork(transmitPacketCh, stateUpdateCh, activeElevatorsCh, config.TRANSMIT_INTERVAL, config.ELEVATOR_TIMEOUT, config.LAST_UPDATE_INTERVAL, config.TRANSMIT_PORT)

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
