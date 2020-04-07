package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"./config"
	"./elevio"
	"./fsm"
	"./network/network2"
	"./statetable"
)

func main() {
	var elevNr int
	var port string
	flag.IntVar(&elevNr, "elevNr", 1, "Specify the elevator nr")
	flag.StringVar(&port, "port", "32001", "Specify a port corresponding to an elevator")
	flag.Parse()

	numFloors := 4
	ip := "localhost:" + port

	intport, _ := strconv.Atoi(port)
	statetable.InitStateTable(elevNr, intport)
	// network2.Init(transmitPacketCh)
	elevio.Init(ip, numFloors)
	fsm.InitFSM(elevNr)
	for true {

	}

	// 	//initialization for simulator
	//numFloors := 4
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
