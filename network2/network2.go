package network2

import (
	"time"

	. "../config"

	"../network/bcast"
)

func transmitStateTable(stateTable [7][9]int, ID string, transmitPacketCh chan<- ElevatorState) {
	statePacket := ElevatorState{ID: ID, StateTable: stateTable}
	transmitPacketCh <- statePacket
}

//BroadcastElevatorState broadcasts elevator state with a given interval.
func BroadcastElevatorState(transmitPacketCh <-chan ElevatorState, elevatorStateTxCh chan<- ElevatorState, transmitInterval time.Duration) {
	transmissionTicker := time.NewTicker(transmitInterval)
	elevatorStateTx := <-transmitPacketCh
	for {
		select {
		case transmitPacket := <-transmitPacketCh:
			elevatorStateTx = transmitPacket
		case <-transmissionTicker.C:
			elevatorStateTxCh <- elevatorStateTx
		default:
			//do stuff
		}
	}
}

//ListenElevatorState listens for elevator state packets, sends to update channel if necessary
func ListenElevatorState(elevatorStateRxCh <-chan ElevatorState, stateUpdateCh chan<- ElevatorState, lostIDCh chan<- string, lifeSignalIDCh chan<- string, timeout time.Duration, offlineTickerInterval time.Duration) {

	//ticker to check for elevators gone offline
	ticker := time.NewTicker(offlineTickerInterval)

	lastUpdate := make(map[string]time.Time)
	receivedPacket := <-elevatorStateRxCh
	lastUpdate[receivedPacket.ID] = time.Now()

	for {
		select {
		case newPacket := <-elevatorStateRxCh:

			receivedPacket = newPacket
			lastUpdate[receivedPacket.ID] = time.Now()
			stateUpdateCh <- receivedPacket
			lifeSignalIDCh <- receivedPacket.ID
		case <-ticker.C:
			for ID, t := range lastUpdate {
				if time.Now().Sub(t) > timeout {
					//fmt.Println("not to worry, we're still flying half a ship")
					lostIDCh <- ID

				}
			}
		default:
			//do stuff

		}

	}
}

//MonitorActiveElevators outputs (on a channel) a map with elevator IDs as keys and bools, indicating if they're active or not, as values
func MonitorActiveElevators(lostIDCh <-chan string, lifeSignalIDCh <-chan string, activeElevatorsCh chan<- map[string]bool) {
	activeElevators := make(map[string]bool)
	for {
		select {
		case lifeSignalID := <-lifeSignalIDCh:
			if mapValue, ok := activeElevators[lifeSignalID]; ok {
				if !mapValue {
					activeElevators[lifeSignalID] = true
					activeElevatorsCh <- activeElevators
				}
			} else {
				activeElevators[lifeSignalID] = true
				activeElevatorsCh <- activeElevators
			}

		case lostID := <-lostIDCh:
			if mapValue, ok := activeElevators[lostID]; ok {
				if mapValue {
					activeElevators[lostID] = false
					activeElevatorsCh <- activeElevators
				}
			}
		default:
			//do stuff
		}
	}
}
func RunNetwork(transmitPacketCh <-chan ElevatorState, stateUpdateCh chan<- ElevatorState, activeElevatorsCh chan<- map[string]bool, TRANSMIT_INTERVAL time.Duration,
	ELEVATOR_TIMEOUT time.Duration, LAST_UPDATE_INTERVAL time.Duration, TRANSMIT_PORT int) {
	elevatorStateTxCh := make(chan ElevatorState)
	elevatorStateRxCh := make(chan ElevatorState)

	lostIDCh := make(chan string)
	lifeSignalIDCh := make(chan string)

	go BroadcastElevatorState(transmitPacketCh, elevatorStateTxCh, TRANSMIT_INTERVAL)
	go ListenElevatorState(elevatorStateRxCh, stateUpdateCh, lostIDCh, lifeSignalIDCh, ELEVATOR_TIMEOUT, LAST_UPDATE_INTERVAL)

	go bcast.Transmitter(TRANSMIT_PORT, elevatorStateTxCh)
	go bcast.Receiver(TRANSMIT_PORT, elevatorStateRxCh)

	go MonitorActiveElevators(lostIDCh, lifeSignalIDCh, activeElevatorsCh)

}
