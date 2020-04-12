package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	. "./config"
	"./elevio"
	"./fsm"
	"./packetprocessor"
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
	// network2.stateTable[row][col+elevNr*3] = valInit(transmitPacketCh)
	elevio.Init(ip, numFloors)

	transmitStateCh := make(chan ElevatorState)
	stateTableTransmitCh := make(chan [7][9]int)
	receiveStateCh := make(chan ElevatorState)
	activeElevatorsCh := make(chan map[string]bool)

	fsm.InitFSM(elevNr, stateTableTransmitCh)

	go packetprocessor.PacketInterchange(transmitStateCh, receiveStateCh, activeElevatorsCh, StateTransmissionInterval, ElevatorTimeout, LastUpdateInterval, ActiveElevatorsTransmissionInterval, TransmissionPort)

	go statetable.UpdateStateTableFromPacket(receiveStateCh)
	go statetable.TransmitState(stateTableTransmitCh, port, transmitStateCh)
	go statetable.UpdateActiveElevators(activeElevatorsCh)

	ticker := time.NewTicker(1000 * time.Millisecond)
	stateTable := statetable.Get()
	for {
		select {
		case <-ticker.C:
			stateTable = statetable.Get()
			for i := 0; i < 7; i++ {
				fmt.Println(stateTable[i])
			}
			fmt.Printf("\n")
		default:
			//do nothing
		}
	}

	/*
		for {
			select {
			case receivedState := <-receiveStateCh:
				for i := 0; i < 7; i++ {
					fmt.Println(receivedState.StateTable[i])
				}
				fmt.Printf("\n")
			}
		}
	*/
}
