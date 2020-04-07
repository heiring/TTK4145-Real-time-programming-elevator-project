package network2

import (
	"time"
)

//ElevatorState ...
type ElevatorState struct {
	ID string
	//state table
}

//BroadcastElevatorState broadcasts elevator state. Sends packets to be sent to transmission channel
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
	//convert ns to s
	//timeout := timeoutns * time.Millisecond
	//offlineTickerInterval := offlineTickerIntervalns * 1000000000

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
				//fmt.Printf(ID + ": ")
				//fmt.Println(t)
				if time.Now().Sub(t) > timeout {
					//fmt.Println("not to worry, we're still flying half a ship")
					lostIDCh <- ID
					//fmt.Println("lostIDCh")
				}
			}
		default:
			//fmt.Println("ListenElevatorState default case")

		}

	}
}

//MonitorActiveElevators outputs (on a channel) a map with elevator IDs as keys and bools, indicating if they're active or not, as values
func MonitorActiveElevators(lostIDCh <-chan string, lifeSignalIDCh <-chan string, activeElevatorsCh chan<- map[string]bool) {
	activeElevators := make(map[string]bool)
	for {
		select {
		case lifeSignalID := <-lifeSignalIDCh:
			//is the elevator ID a key in the map?
			if mapValue, ok := activeElevators[lifeSignalID]; ok {
				//is an elevator back online?
				if !mapValue {
					activeElevators[lifeSignalID] = true
					//the set of active elevators has changed
					activeElevatorsCh <- activeElevators
				}
			} else { //the elevator ID is not a key in the map
				activeElevators[lifeSignalID] = true
				//the set of active elevators has changed
				activeElevatorsCh <- activeElevators
			}

		case lostID := <-lostIDCh:
			//is the elevator ID a key in the map?
			if mapValue, ok := activeElevators[lostID]; ok {
				if mapValue {
					activeElevators[lostID] = false
					//the set of active elevators has changed
					activeElevatorsCh <- activeElevators
				}
			}
		default:
			//do stuff
		}
	}
}
