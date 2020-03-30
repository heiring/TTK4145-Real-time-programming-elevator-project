package main

import (
	"./fsm"
	// "fmt"
	// "time"
	// "./elevio"
	// "./network"
)

func main() {
	fsm.FSM(1)
	// elevatorSliceCh := make(chan [3]network.Elevator)
	//var counter = 0
	// go network.ElevatorLifeStatusMonitor(elevatorSliceCh)

	//initialization for simulator
	// numFloors := 3
	// elevio.Init("localhost:15657", numFloors)

	// for {
	// 	/*
	// 		select {
	// 		case p := <-elevatorSliceCh:

	// 			fmt.Printf("iteration ")
	// 			counter++
	// 			fmt.Println(counter)

	// 			for _, elevator := range p {

	// 				fmt.Printf("elevator with id: ")
	// 				fmt.Printf(elevator.Id)
	// 				fmt.Printf("is alive? ")
	// 				fmt.Println(elevator.IsAlive)

	// 			}
	// 			//time.Sleep(1000 * time.Millisecond)
	// 		default:
	// 			//do nothing
	// 		}
	// 	*/

	// }

}