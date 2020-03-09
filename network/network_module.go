package network

import (
	"fmt"
	"os"

	"./utilities/localip"
	"./utilities/peers"
)

type Elevator struct {
	Id      string
	IsAlive bool
}

var elevator0 Elevator
var elevator1 Elevator
var elevator2 Elevator

var elevatorSlice [3]Elevator

func initializeElevatorStructs() {
	elevatorSlice[0].IsAlive = false
	elevatorSlice[1].IsAlive = false
	elevatorSlice[2].IsAlive = false
	//Id??
}

// ElevatorLifeStatusMonitor outputs a struct with elevator structs containing life status
func ElevatorLifeStatusMonitor(elevatorSliceCh chan<- [3]Elevator) {

	initializeElevatorStructs()
	initialUpdateComplete := false

	//slice containing 3 elevator structs
	elevatorSlice := [3]Elevator{elevator0, elevator1, elevator2}

	//initialize peerUpdate communication
	localIP, _ := localip.LocalIP()
	Id := fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, Id, peerTxEnable) //which port???
	go peers.Receiver(15647, peerUpdateCh)

	for {
		select {
		case p := <-peerUpdateCh:

			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			if len(p.Peers) == 3 && initialUpdateComplete == false {
				//assume peer Id always remain the same for every elevator ??
				elevatorSlice[0].IsAlive = true
				elevatorSlice[0].Id = p.Peers[0]

				elevatorSlice[1].IsAlive = true
				elevatorSlice[1].Id = p.Peers[1]

				elevatorSlice[2].IsAlive = true
				elevatorSlice[2].Id = p.Peers[2]

				initialUpdateComplete = true
			}

			if initialUpdateComplete {
				//elevator is no longer alive if its Id is in p.Lost:
				for _, Id := range p.Lost {
					for _, elevator := range elevatorSlice {
						if elevator.Id == Id {
							elevator.IsAlive = false
						}
					}
				}
				//if an elevator Id is in p.New, it's now alive
				for _, elevator := range elevatorSlice {
					if elevator.Id == p.New {
						elevator.IsAlive = true
					}
				}
			}
			elevatorSliceCh <- elevatorSlice

		default:
			// Do nothing

		}

	}
}
