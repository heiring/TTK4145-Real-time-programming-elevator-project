package packetprocessor

import (
	"time"

	. "../config"

	"../network/bcast"
)

func transmitStateTable(stateTable [7][3]int, ID string, transmitStateCh chan<- ElevatorState) {
	statePacket := ElevatorState{ID: ID, StateTable: stateTable}
	transmitStateCh <- statePacket
}

//BroadcastElevatorState broadcasts elevator state with a given interval.
func BroadcastElevatorState(transmitStateCh <-chan ElevatorState, elevatorStateTxCh chan<- ElevatorState, transmitInterval time.Duration) {
	transmissionTicker := time.NewTicker(transmitInterval)
	elevatorStateTx := <-transmitStateCh
	for {
		select {
		case transmitState := <-transmitStateCh:
			elevatorStateTx = transmitState
		case <-transmissionTicker.C:
			elevatorStateTxCh <- elevatorStateTx
		default:
			//do stuff
		}
	}
}

//ListenElevatorState listens for elevator state packets, sends to update channel if necessary
func ListenElevatorState(elevatorStateRxCh <-chan ElevatorState, receiveStateCh chan<- ElevatorState, lostIDCh chan<- string, lifeSignalIDCh chan<- string, timeout time.Duration, offlineTickerInterval time.Duration) {

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
			receiveStateCh <- receivedPacket
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
func MonitorActiveElevators(lostIDCh <-chan string, lifeSignalIDCh <-chan string, activeElevatorsCh chan<- map[string]bool, activeElevatorsTransmitInterval time.Duration) {
	activeElevators := make(map[string]bool)
	emptyMap := true
	ticker := time.NewTicker(activeElevatorsTransmitInterval)
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
				emptyMap = false
			}

		case lostID := <-lostIDCh:
			if mapValue, ok := activeElevators[lostID]; ok {
				if mapValue {
					activeElevators[lostID] = false
					activeElevatorsCh <- activeElevators
				}
			}
		case <-ticker.C:
			if !emptyMap {
				activeElevatorsCh <- activeElevators
			}

		default:
			//do stuff
		}
	}
}
func PacketInterchange(transmitStateCh <-chan ElevatorState, receiveStateCh chan<- ElevatorState, activeElevatorsCh chan<- map[string]bool, StateTransmissionInterval time.Duration,
	ElevatorTimeout time.Duration, LastUpdateInterval time.Duration, activeElevatorsTransmitInterval time.Duration, TransmissionPort int) {
	elevatorStateTxCh := make(chan ElevatorState)
	elevatorStateRxCh := make(chan ElevatorState)

	lostIDCh := make(chan string)
	lifeSignalIDCh := make(chan string)

	go BroadcastElevatorState(transmitStateCh, elevatorStateTxCh, StateTransmissionInterval)
	go ListenElevatorState(elevatorStateRxCh, receiveStateCh, lostIDCh, lifeSignalIDCh, ElevatorTimeout, LastUpdateInterval)

	go bcast.Transmitter(TransmissionPort, elevatorStateTxCh)
	go bcast.Receiver(TransmissionPort, elevatorStateRxCh)

	go MonitorActiveElevators(lostIDCh, lifeSignalIDCh, activeElevatorsCh, activeElevatorsTransmitInterval)

}
