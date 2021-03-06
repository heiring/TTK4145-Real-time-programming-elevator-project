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
	stateTableTransmitCh := make(chan [7][3]int)
	receiveStateCh := make(chan statetable.ElevatorState)
	activeElevatorsCh := make(chan map[string]bool)

	saveStateForRecoveryCh := make(chan statetable.ElevatorState)
	recoveryIDCh := make(chan string)
	transmitRecoveryStateCh := make(chan statetable.ElevatorState)

	elevio.Init(ip, numFloors)
	statetable.InitStateTable(intport)

	go statetable.StateTableRecovery(saveStateForRecoveryCh, recoveryIDCh, transmitRecoveryStateCh)
	go statetable.UpdateStateTableFromPacket(receiveStateCh, stateTableTransmitCh)
	go statetable.TransmitState(stateTableTransmitCh, transmitStateCh, transmitRecoveryStateCh)
	// fsm.InitFSM(stateTableTransmitCh)
	go fsm.PollHardwareActions(stateTableTransmitCh)
	go packetprocessor.PacketInterchange(transmitStateCh, receiveStateCh, activeElevatorsCh, StateTransmissionInterval, ElevatorTimeout, LastUpdateInterval, ActiveElevatorsTransmissionInterval, TransmissionPort)

	go statetable.UpdateActiveElevators(activeElevatorsCh, saveStateForRecoveryCh, recoveryIDCh)

	for true {
		time.Sleep(2 * time.Second)
		// st := statetable.StateTables.ReadWholeMap()
		// fmt.Println("StateTables:")
		// fmt.Println("________")
		// for row := 0; row < 7; row++ {
		// 	fmt.Print("|", st["32000"][row][0], "|")
		// 	fmt.Print("|", st["32000"][row][1], "|")
		// 	fmt.Print("|", st["32000"][row][2], "|")
		// 	fmt.Print("\t\t")
		// 	fmt.Print("|", st["32001"][row][0], "|")
		// 	fmt.Print("|", st["32001"][row][1], "|")
		// 	fmt.Print("|", st["32001"][row][2], "|\n")
		// }
		// fmt.Println("________")
	}

}
