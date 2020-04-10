package main

import (
	"fmt"
	"os"

	. "./config"
	"./elevio"
	"./packetprocessor"
)

func main() {
	/*
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
	*/
	//initialization for simulator
	numFloors := 4
	ID := os.Args[1]
	elevio.Init("localhost:"+ID, numFloors)

	transmitStateCh := make(chan ElevatorState)
	receiveStateCh := make(chan ElevatorState)
	activeElevatorsCh := make(chan map[string]bool)

	go packetprocessor.PacketInterchange(transmitStateCh, receiveStateCh, activeElevatorsCh, StateTransmissionInterval, ElevatorTimeout, LastUpdateInterval, ActiveElevatorsTransmissionInterval, TransmissionPort)

	msg := ElevatorState{ID: ID}

	for {
		transmitStateCh <- msg
		select {
		case y := <-receiveStateCh:
			//fmt.Println("main : packet received")
			//fmt.Println(y.ID)
			y.ID = "1111"
		case activeElevators := <-activeElevatorsCh:
			for ID, isAlive := range activeElevators {
				fmt.Printf(ID + ": ")
				fmt.Println(isAlive)
			}
			fmt.Printf("\n")
		default:
			//do stuff
		}
	}
}
