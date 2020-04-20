package main

import (
	"flag"
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
	flag.StringVar(&port, "port", "32001", "Specify a port corresponding to an elevator")
	flag.Parse()

	numFloors := 4
	ip := "localhost:" + port

	intport, _ := strconv.Atoi(port)

	transmitStateCh := make(chan statetable.ElevatorState)
	stateTablesTransmitCh := make(chan map[string][7][3]int)
	receiveStateCh := make(chan statetable.ElevatorState)
	activeElevatorsCh := make(chan map[string]bool)

	elevio.Init(ip, numFloors)
	statetable.InitStateTable(intport)

	go statetable.UpdateStateTableFromPacket(receiveStateCh, stateTablesTransmitCh)
	go statetable.TransmitState(stateTablesTransmitCh, transmitStateCh)
	// fsm.InitFSM(stateTableTransmitCh)
	go fsm.PollHardwareActions(stateTablesTransmitCh)
	go packetprocessor.PacketInterchange(transmitStateCh, receiveStateCh, activeElevatorsCh, StateTransmissionInterval, ElevatorTimeout, LastUpdateInterval, ActiveElevatorsTransmissionInterval, TransmissionPort)

	go statetable.UpdateActiveElevators(activeElevatorsCh)

	for true {
		time.Sleep(10 * time.Second)
	}

}
