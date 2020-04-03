package network2

import (
	"fmt"
	"time"

	"../localip"
)

//ElevatorState ...
type ElevatorState struct {
	ID      string
	IsAlive bool
	//state table
}

//BroadcastElevatorState broadcasts elevator state. Sends packets to be sent to transmission channel
func BroadcastElevatorState(transmitPacketCh <-chan ElevatorState, elevatorStateTxCh chan<- ElevatorState, transmitInterval time.Duration) {
	transmissionTicker := time.NewTicker(transmitInterval * time.Millisecond)
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
func ListenElevatorState(elevatorStateRxCh <-chan ElevatorState, stateUpdateCh chan<- ElevatorState, timeoutns time.Duration, lostIDCh chan<- string, offlineTickerIntervalns time.Duration) {
	//convert ns to s
	timeout := timeoutns * 1000000000
	offlineTickerInterval := offlineTickerIntervalns * 1000000000

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
			//is elevator back online?
			if receivedPacket.IsAlive == false {
				receivedPacket.IsAlive = true
			}
			stateUpdateCh <- receivedPacket
		case <-ticker.C:
			for ID, t := range lastUpdate {
				fmt.Printf(ID + ": ")
				fmt.Println(t)
				if time.Now().Sub(t) > timeout {
					fmt.Println("not to worry, we're still flying half a ship")
					//lostIDCh <- ID
					fmt.Println("lostIDCh")
				}
			}
		default:
			//fmt.Println("ListenElevatorState default case")

		}

	}
}

//NetworkTest yeet
func NetworkTest(transmitPacketCh chan<- ElevatorState, stateUpdateCh <-chan ElevatorState) {
	ip, _ := localip.LocalIP()
	localElevator := ElevatorState{ID: ip, IsAlive: true}
	transmitPacketCh <- localElevator

	ticker := time.NewTicker(1000 * time.Millisecond)
	for {

		select {
		case stateUpdate := <-stateUpdateCh:
			fmt.Println(stateUpdate.ID)
			fmt.Println("%v", stateUpdate.IsAlive)
			//finished <- true

		case <-ticker.C:
			fmt.Println("nothing recieved")
			transmitPacketCh <- localElevator
		default:
			//do nothing
		}
	}

}

//NetworkTest2 brrra
func NetworkTest2(transmitPacketCh chan<- ElevatorState, stateUpdateCh <-chan ElevatorState) {
	yeet := ElevatorState{ID: "2222", IsAlive: true}

	for {
		transmitPacketCh <- yeet
		select {
		case y := <-stateUpdateCh:
			fmt.Println("main : packet received")
			fmt.Println(y.ID)
		default:
			//do stuff
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
