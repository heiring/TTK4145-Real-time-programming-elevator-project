package main

import (
	//"./fsm"
	."./network"
	"fmt"
	"time"
)

func main() {
	//fsm.FSM(1)
	elevatorSliceCh := make(chan [3]Elevator)
	var counter = 0
	go ElevatorLifeStatusMonitor(elevatorSliceCh)

	for {
		
		fmt.Printf("iteration ")
		counter ++
		fmt.Println(counter)		

		select {
		case p := <- elevatorSliceCh:
			for _,elevator := range p{
				
				fmt.Printf("elevator with id: ")
				fmt.Printf(elevator.Id)
				fmt.Printf( "is alive? ")
				fmt.Println(elevator.IsAlive)
				
			}
			time.Sleep(1000 * time.Millisecond) 
		}
		
	}
	
	
}
