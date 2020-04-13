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

	var port string
	// flag.IntVar(&elevNr, "elevNr", 1, "Specify the elevator nr")
	flag.StringVar(&port, "port", "32001", "Specify a port corresponding to an elevator")
	flag.Parse()

	numFloors := 4
	ip := "localhost:" + port

	intport, _ := strconv.Atoi(port)
	statetable.InitStateTable(intport)
	fmt.Println("STATETABLE:\n", statetable.StateTables[port])
	// network2.stateTable[row][col+elevNr*3] = valInit(transmitPacketCh)
	elevio.Init(ip, numFloors)

	transmitStateCh := make(chan ElevatorState)
	stateTableTransmitCh := make(chan [7][3]int)
	receiveStateCh := make(chan ElevatorState)
	activeElevatorsCh := make(chan map[string]bool)

	fsm.InitFSM(stateTableTransmitCh)

	go packetprocessor.PacketInterchange(transmitStateCh, receiveStateCh, activeElevatorsCh, StateTransmissionInterval, ElevatorTimeout, LastUpdateInterval, ActiveElevatorsTransmissionInterval, TransmissionPort)

	go statetable.UpdateStateTableFromPacket(receiveStateCh)
	go statetable.TransmitState(stateTableTransmitCh, transmitStateCh)
	go statetable.UpdateActiveElevators(activeElevatorsCh)

	ticker := time.NewTicker(1000 * time.Millisecond)
	stateTables := statetable.GetStateTables()
	for {
		select {
		case <-ticker.C:
			stateTables = statetable.GetStateTables()
			fmt.Print("localID: ")
			fmt.Println(statetable.GetLocalID())
			for i := 0; i < 7; i++ {
				fmt.Print(stateTables["15000"][i])
				fmt.Print("				")
				fmt.Println(stateTables["16000"][i])
			}
			fmt.Printf("\n")
		default:
			//do nothing
		}
	}

}
