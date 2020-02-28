package main

import (
	"fmt"

	"./elevio"
)

func main() {
	var elev_nr int = 0

	drvButtons := make(chan elevio.ButtonEvent)
	drvFloors := make(chan int)
	drvObstr := make(chan bool)
	drvStop := make(chan bool)

	// curState := make(chan FSM.Elev_state)
	dir := make(chan elevio.MotorDirection)

	// curState <- FSM.INIT

	// initializeState(<-dir)

	go elevio.PollButtons(drvButtons)
	go elevio.PollFloorSensor(drvFloors)
	go elevio.PollObstructionSwitch(drvObstr)
	go elevio.PollStopButton(drvStop)
	// go FSM.FSM(curState)

	// Detect state changes and update StateTable
	for {
		select {
		case a := <-drvButtons:
			fmt.Printf("%+b\n", a)
			var row int = 3 + a.Floor
			var col int = int(a.Button)
			elevio.SetButtonLamp(a.Button, a.Floor, true)
			StateTable.UpdateStateTableIndex(row, col, elev_nr, 1)

		case a := <-drvFloors:
			fmt.Printf("%+v\n", a)
			var row int = 2
			var col int = int(<-dir) + 1
			elevio.SetFloorIndicator(a)
			StateTable.ResetElevRow(row, elev_nr)
			StateTable.UpdateElevLastFLoor(elev_nr, a)
		}
	}
}
