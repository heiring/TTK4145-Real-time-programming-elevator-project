package packetprocessor

import (
	"time"

	"../network/bcast"
	"../statetable"
)

// This module inputs the elevator states to be broadcasted and outputs the elevator states that are sent from other elevator as well as a map indicating which elevators are online/offline

//BroadcastElevatorState transmits the local elevator state repeatedly (with a given interval).
func broadcastElevatorState(transmitStateCh <-chan statetable.ElevatorState, elevatorStateTxCh chan<- statetable.ElevatorState, transmitInterval time.Duration) {
	transmissionTicker := time.NewTicker(transmitInterval)
	elevatorStateTx := <-transmitStateCh
	for {
		select {
		case transmitState := <-transmitStateCh:
			elevatorStateTx = transmitState
		case <-transmissionTicker.C:
			elevatorStateTxCh <- elevatorStateTx
		default:
		}
	}
}

//ListenElevatorState is the first to receive info from other elevators and this goroutine evaluates the amount of time passed since an elevator last sent a packet.
// If this time exceeds a given limit, it is labeled as an inactive elevator.
func checkTimePassedSinceLastUpdate(elevatorStateRxCh <-chan statetable.ElevatorState, receiveStateCh chan<- statetable.ElevatorState, lostIDCh chan<- string, lifeSignalIDCh chan<- string, timeout time.Duration, offlineTickerInterval time.Duration) {
	ticker := time.NewTicker(offlineTickerInterval)
	lastUpdate := make(map[string]time.Time)
	receivedPacket := <-elevatorStateRxCh
	lastUpdate[receivedPacket.ID] = time.Now()
	for {
		select {
		case receivedPacket = <-elevatorStateRxCh:
			lastUpdate[receivedPacket.ID] = time.Now()
			receiveStateCh <- receivedPacket
			lifeSignalIDCh <- receivedPacket.ID
		case <-ticker.C:
			for ID, t := range lastUpdate {
				if time.Now().Sub(t) > timeout {
					lostIDCh <- ID
				}
			}
		default:
		}

	}
}

//MonitorActiveElevators outputs (on a channel) a map with elevator IDs as keys and bools, indicating if they're active or not, as values
func monitorActiveElevators(lostIDCh <-chan string, lifeSignalIDCh <-chan string, activeElevatorsCh chan<- map[string]bool, activeElevatorsTransmitInterval time.Duration) {
	activeElevators := make(map[string]bool)
	emptyMap := true
	ticker := time.NewTicker(activeElevatorsTransmitInterval)
	for {
		select {
		case lifeSignalID := <-lifeSignalIDCh:
			if mapValue, ok := activeElevators[lifeSignalID]; ok {
				if !mapValue {
					activeElevators[lifeSignalID] = true
				}
			} else {
				activeElevators[lifeSignalID] = true
				emptyMap = false
			}

		case lostID := <-lostIDCh:
			if mapValue, ok := activeElevators[lostID]; ok {
				if mapValue {
					activeElevators[lostID] = false
				}
			}
		case <-ticker.C:
			if !emptyMap {
				activeElevatorsCh <- activeElevators
			}

		default:
		}

	}
}

// PacketInterchange has the input transmitStateCh which is a channel containing the info about the local elevator we want to broadcast.
// The first output (receiveStateCh) is a channel containing info sent from other elevator in the network.
// The second output (activeElevatorsCh) is a channel containing a map with the elevator IDs as keys and a bool indicating their life status as values.
func PacketInterchange(transmitStateCh <-chan statetable.ElevatorState, receiveStateCh chan<- statetable.ElevatorState, activeElevatorsCh chan<- map[string]bool, StateTransmissionInterval time.Duration,
	ElevatorTimeout time.Duration, LastUpdateInterval time.Duration, activeElevatorsTransmitInterval time.Duration, TransmissionPort int) {
	elevatorStateTxCh := make(chan statetable.ElevatorState)
	elevatorStateRxCh := make(chan statetable.ElevatorState)

	lostIDCh := make(chan string)
	lifeSignalIDCh := make(chan string)

	go broadcastElevatorState(transmitStateCh, elevatorStateTxCh, StateTransmissionInterval)
	go checkTimePassedSinceLastUpdate(elevatorStateRxCh, receiveStateCh, lostIDCh, lifeSignalIDCh, ElevatorTimeout, LastUpdateInterval)

	go bcast.Transmitter(TransmissionPort, elevatorStateTxCh)
	go bcast.Receiver(TransmissionPort, elevatorStateRxCh)

	go monitorActiveElevators(lostIDCh, lifeSignalIDCh, activeElevatorsCh, activeElevatorsTransmitInterval)
}
