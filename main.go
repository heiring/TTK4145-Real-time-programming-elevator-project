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
	// ************************************************************
	// ***						Known bugs						***
	// ************************************************************
	//
	// * When elev is right next to current-order-floor, and a new order is received, the old order is not completed.
	//
	// * FIXED External hall buttons seems to be bugged. When an external order is completed, the lights update perfectly,
	//	 but the order is still active, so when a new local order is received and completed, elev returns to the old
	// 	 hall order.
	// - Potential fix: The extertnal elevator receives news of the order completion and updates its statetable,
	//	 however the old statetable is still being transmitted on stateTableTransmitCh - since FSM only reacts to local
	//	 statetable changes.
	//
	// * FIXED Any Hall btn pressed on local elev, with local elev already at the floor, no elevs take the order, the order is never completed
	// 	 BUT! When the local elev starts moving, the external will move in and complete the order.
	//
	// * Ish 50 % of the time a button has to be pressed two times for the system to notice it
	// * Light door lamp!!
	// * When an elevator dies (i.e. network cable unplugged) the elevator should function on its own
	//
	// * When an elev returns online, externally completed orders must be unlit
	//
	// * Elev0, and 1 at floor 3. Elev2 at floor 0. Hall_Down at floor 1 pushed with elev2. All elevs takes the order (true 4)
	//   Happens only for more than 2 elevs.
	// * Implement stop button!
	//
	// * If new cab order is received when elev is at floor, it begins to move to new order without waiting
	//
	// * elev does not stop for new orders in same direction.

	var port string
	flag.StringVar(&port, "port", "32003", "Specify a port corresponding to an elevator")
	flag.Parse()

	numFloors := 4
	ip := "localhost:" + port

	intport, _ := strconv.Atoi(port)

	//fmt.Println("STATETABLE:\n", statetable.StateTables[port])
	// network2.stateTable[row][col+elevNr*3] = valInit(transmitPacketCh)

	transmitStateCh := make(chan ElevatorState)
	stateTableTransmitCh := make(chan [7][3]int)
	receiveStateCh := make(chan ElevatorState)
	activeElevatorsCh := make(chan map[string]bool)
	elevio.Init(ip, numFloors)
	statetable.InitStateTable(intport)
	go statetable.UpdateStateTableFromPacket(receiveStateCh, stateTableTransmitCh)
	go statetable.TransmitState(stateTableTransmitCh, transmitStateCh)
	fsm.InitFSM(stateTableTransmitCh)

	go packetprocessor.PacketInterchange(transmitStateCh, receiveStateCh, activeElevatorsCh, StateTransmissionInterval, ElevatorTimeout, LastUpdateInterval, ActiveElevatorsTransmissionInterval, TransmissionPort)

	go statetable.UpdateActiveElevators(activeElevatorsCh)

	// ticker := time.NewTicker(1000 * time.Millisecond)
	// stateTables := statetable.GetStateTables()
	// localID := statetable.GetLocalID()
	for true {
		time.Sleep(10 * time.Second)
		// select {
		// case <-ticker.C:
		// stateTables := statetable.GetStateTables()
		// fmt.Println(stateTables[localID][0][0])
		// time.Sleep(1 * time.Second)
		// 	fmt.Print("localID: ")
		// 	fmt.Println(statetable.GetLocalID())
		// 	for i := 0; i < 7; i++ {
		// 		fmt.Print(stateTables["32000"][i])
		// 		fmt.Print("			")
		// 		fmt.Println(stateTables["32001"][i])
		// 	}

		// default:
		// 	// do nothing
		// }
	}

}
