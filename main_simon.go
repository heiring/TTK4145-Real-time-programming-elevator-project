package main

import (
	"fmt"
	"./elevio"
	"./StateTable"
)

func main() {
	var elev_nr int = 0

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	
	// curState := make(chan FSM.Elev_state)
	dir := make(chan elevio.MotorDirection)
	
	// curState <- FSM.INIT

	// initializeState(<-dir)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)
	// go FSM.FSM(curState)



	// Detect state changes and update StateTable
	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+b\n", a)
			var row int = 5 + a.Floor
			var col int = int(a.Button)
			elevio.SetButtonLamp(a.Button, a.Floor, true)
			StateTable.UpdateStateTable(row, col, elev_nr, 1)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			var row int = 1 + a
			var col int = int(<-dir) + 1
			elevio.SetFloorIndicator(a)
			StateTable.UpdateStateTable(row, col, elev_nr, 1) // Add current floor
			for col := 0; col < 3; col++ {
				StateTable.UpdateStateTable(row - int(<-dir), col, elev_nr, 0) // Remove previous floor
			}
		}
	}
}
