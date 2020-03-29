package network2

import (
	"fmt"
	"time"

	"../localip"
)

//ElevatorState ...
type ElevatorState struct {
	ID      string
	isAlive bool
	//state table
}

//BroadcastElevatorState broadcasts elevator state. Sends packets to be sent to transmission channel
func BroadcastElevatorState(transmitPacketCh <-chan ElevatorState, elevatorStateTxCh chan<- ElevatorState, transmitInterval time.Duration) {
	ticker := time.NewTicker(transmitInterval * time.Millisecond)
	elevatorStateTx := <-transmitPacketCh
	fmt.Println("elevatorStateTx := <-transmitPacketCh")
	for {
		select {
		case transmitPacket := <-transmitPacketCh:
			elevatorStateTx = transmitPacket
			fmt.Println("case transmitPacket := <-transmitPacketCh")
		case <-ticker.C:
			elevatorStateTxCh <- elevatorStateTx
		default:
			//do stuff
		}
	}
}

//ListenElevatorState listens for elevator state packets, sends to update channel if necessary
func ListenElevatorState(elevatorStateRxCh <-chan ElevatorState, stateUpdateCh chan<- ElevatorState, Timeout time.Duration, lostIDCh chan<- string) {
	lastUpdate := make(map[string]time.Time)
	receivedPacket := <-elevatorStateRxCh
	lastUpdate[receivedPacket.ID] = time.Now()
	for {
		select {
		case newPacket := <-elevatorStateRxCh:
			receivedPacket = newPacket
			lastUpdate[receivedPacket.ID] = time.Now()
			//is elevator back online?
			if receivedPacket.isAlive == false {
				receivedPacket.isAlive = true
			}

		default:
			//is an elevator offline?
			for ID, t := range lastUpdate {
				if time.Now().Sub(t) > Timeout {
					lostIDCh <- ID //HÅNDTER DETTE
				}

			}
			stateUpdateCh <- receivedPacket
		}

	}
}

func NetworkTest(transmitPacketCh chan<- ElevatorState, stateUpdateCh <-chan ElevatorState) {
	ip, _ := localip.LocalIP()
	localElevator := ElevatorState{ID: ip, isAlive: true}
	transmitPacketCh <- localElevator

	ticker := time.NewTicker(1000 * time.Millisecond)
	for {

		select {
		case stateUpdate := <-stateUpdateCh:
			fmt.Println(stateUpdate.ID)
			fmt.Println("%v", stateUpdate.isAlive)
			//finished <- true

		case <-ticker.C:
			fmt.Println("nothing recieved")
			transmitPacketCh <- localElevator
		default:
			//do nothing
		}
	}

}

/*
//UpdateElevatorLifeStatus checks if an elevator has gone offline and if so, outputs lost ID on the channel lostIDCh
func UpdateElevatorLifeStatus(lastUpdate map[string]time.Time, Timeout time.Duration, lostIDCh chan<- string) {
	for {
		for ID, t := range lastUpdate {
			if time.Now().Sub(t) > Timeout {
				lostIDCh <- ID
			}

		}

	}
}

func IsElevatorBackOnline (lostIDCh <-chan string, stateUpdateCh <-chan ElevatorState){
	lostIDs := []string
	for {
		select{
		case lostID := <- lostIDCh
			lostIDs = append(lostIDs, lostID)

		default:

		}
	}
}

*/
