package network

import (
	"fmt"
	"os"
	"time"

	"./utilities/localip"
	"./utilities/peers"
)

type elevator struct {
	id      string
	isAlive bool
}

var elevator0 elevator
var elevator1 elevator
var elevator2 elevator

var elevatorSlice [3]elevator

func initializeElevatorStructs() {
	elevatorSlice[0].isAlive = false
	elevatorSlice[1].isAlive = false
	elevatorSlice[2].isAlive = false
	//id??
}

// AliveTransmission transmits an AliveMsg every second
func AliveTransmission(elevatorSliceCh chan<- [3]elevator) {

	initializeElevatorStructs()
	initialUpdateComplete := false

	//slice containing 3 elevator structs
	elevatorSlice := [3]elevator{elevator0, elevator1, elevator2}

	//initialize peerUpdate communication
	localIP, _ := localip.LocalIP()
	id := fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable) //which port???
	go peers.Receiver(15647, peerUpdateCh)

	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			if len(p.Peers) == 3 && initialUpdateComplete == false {
				elevatorSlice[0].isAlive = true
				elevatorSlice[0].id = p.Peers[0]

				elevatorSlice[1].isAlive = true
				elevatorSlice[1].id = p.Peers[1]

				elevatorSlice[2].isAlive = true
				elevatorSlice[2].id = p.Peers[2]

				initialUpdateComplete = true
			}

			if initialUpdateComplete {
				//elevator is no longer alive if its ID is in p.Lost:
				for _, id := range p.Lost {
					for _, elevator := range elevatorSlice {
						if elevator.id == id {
							elevator.isAlive = false
						}
					}
				}
				//if an elevator ID is in p.New, it's now alive
				for _, elevator := range elevatorSlice {
					if elevator.id == p.New {
						elevator.isAlive = true
					}
				}
			}

		default:
			// Do nothing

		}
		time.Sleep(1000 * time.Millisecond) //??
		elevatorSliceCh <- elevatorSlice

	}
}
