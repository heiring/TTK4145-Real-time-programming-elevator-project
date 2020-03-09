package main

import (
	//"./fsm"
	//"fmt"
	//"time"

	. "./network"
)

func main() {
	//fsm.FSM(1)
	elevatorSliceCh := make(chan [3]Elevator)
	//var counter = 0
	go ElevatorLifeStatusMonitor(elevatorSliceCh)

	for {
		/*
			select {
			case p := <-elevatorSliceCh:

				fmt.Printf("iteration ")
				counter++
				fmt.Println(counter)

				for _, elevator := range p {

					fmt.Printf("elevator with id: ")
					fmt.Printf(elevator.Id)
					fmt.Printf("is alive? ")
					fmt.Println(elevator.IsAlive)

				}
				//time.Sleep(1000 * time.Millisecond)
			default:
				//do nothing
			}
		*/

	}

}
