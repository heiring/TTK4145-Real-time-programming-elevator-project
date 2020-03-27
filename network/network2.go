package network

import (
	"time"
)

//ElevatorState ...
type ElevatorState struct {
	ID string
	//state table
}

//BroadcastElevatorState broadcasts elevator state. Sends packets to be sent to transmission channel
func BroadcastElevatorState(stateUpdateCh <-chan ElevatorState, elevatorStateTxCh chan<- ElevatorState) {
	//elevatorStateTx := ElevatorState{ID: "null"}
	for {
		select {
		case stateUpdate := <-stateUpdateCh:
			elevatorStateTx := stateUpdate
			elevatorStateTxCh <- elevatorStateTx

		default:
			//elevatorStateTxCh <- elevatorStateTx
		}
	}
}

//ListenElevatorState listens for elevator state packets, sends to update channel if necessary
func ListenElevatorState(elevatorStateRxCh <-chan ElevatorState, stateUpdateCh chan<- ElevatorState, TxFrequency time.Duration) {
	lastUpdate := make(map[string]time.Time) //gjør om til parameter!!!
	for {
		select {
		case packetReceived := <-elevatorStateRxCh:

			if time.Now().Sub(lastUpdate[packetReceived.ID]) > TxFrequency {
				stateUpdateCh <- packetReceived
				lastUpdate[packetReceived.ID] = time.Now()
			}
		default:
			//do stuff
		}

	}
}

//UpdateElevatorLifeStatus checks if an elevator has gone offline
func UpdateElevatorLifeStatus(lastUpdate map[string]time.Time, Timeout time.Duration) {
	for {
		for ID, time := range lastUpdate {
			if time > Timeout {
				//do stuff
			}

		}

	}
}
