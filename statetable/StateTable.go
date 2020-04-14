package statetable

import (
	"fmt"
	"strconv"
	"time"

	. "../config"
	"../elevio"
	"../orderdistributor"
)

var StateTables = make(map[string][7][3]int)
var localID string

const UnknownFloor int = -1

func InitStateTable(port int) {
	fmt.Println("InitStateTable")
	var tempStateTable [7][3]int
	for row, cells := range tempStateTable {
		for _, col := range cells {
			tempStateTable[row][col] = 0
		}
	}
	// Set status to active
	tempStateTable[0][0] = 1
	// Unknown starting position
	tempStateTable[2][1] = UnknownFloor
	// Set ID = port
	tempStateTable[0][1] = port

	localID = strconv.Itoa(port)
	StateTables[localID] = tempStateTable
}

func UpdateStateTableFromPacket(receiveStateCh <-chan ElevatorState) {
	for {
		select {
		case elevState := <-receiveStateCh:
			ID := elevState.ID
			if ID != localID {
				StateTables[ID] = elevState.StateTable

				// To do: Update lights
				runOrderDistribution()
			}
		default:
			//do stuff
		}
	}
}
func UpdateLightsFromPacket() {
	for row, cells := range StateTables {
		for col := range cells {
			//
			elevio.SetButtonLamp(butn, curFloor, false)
		}
	}
}

func TransmitState(stateTableTransmitCh <-chan [7][3]int, transmitStateCh chan<- ElevatorState) {
	ticker := time.NewTicker(StateTransmissionInterval)
	stateTable := StateTables[localID]
	elevatorState := ElevatorState{ID: localID, StateTable: stateTable}
	for {
		select {
		case stateTable = <-stateTableTransmitCh:
			elevatorState.StateTable = stateTable
		case <-ticker.C:
			transmitStateCh <- elevatorState
		default:
			//do nothing
		}
	}
}

func UpdateActiveElevators(activeElevatorsCh <-chan map[string]bool) {
	for {
		select {
		case activeElevators := <-activeElevatorsCh: //Packets arrive regularly
			//update state table
			for ID, isAlive := range activeElevators {
				for mapID, _ := range StateTables {
					if mapID == ID {
						if isAlive {
							UpdateStateTableIndex(0, 0, ID, 1, true)
						} else {
							UpdateStateTableIndex(0, 0, ID, 0, true)
							fmt.Println("DANGER")
						}
					}
				}
			}
		default:
			//do stuff
		}
	}
	runOrderDistribution()
}

func UpdateStateTableIndex(row, col int, port string, val int, runDistribution bool) { // stateTableTransmitCh chan<- [7][9]int) {
	statetable := StateTables[port]
	statetable[row][col] = val
	StateTables[port] = statetable
	if runDistribution {
		runOrderDistribution()
	}

}

func runOrderDistribution() {
	orderdistributor.DistributeOrders(string(localID), StateTables)
}

func UpdateElevLastFLoor(val int) {
	UpdateStateTableIndex(2, 1, localID, val, false)
}

func UpdateElevDirection(val int) {
	UpdateStateTableIndex(1, 1, localID, val, false)
}

func ResetElevRow(row int, ID string) {
	for col := 0; col < 3; col++ {
		UpdateStateTableIndex(row, col, ID, 0, false)
	}
}

func ResetRow(row int) {
	for ID, _ := range StateTables {
		ResetElevRow(row, ID)
	}
}

func getPositionRow(port string) int {
	position := StateTables[port][2][1]
	return position
}

func GetElevDirection(port string) int {
	direction := StateTables[port][1][1]
	return direction
}

func GetCurrentFloor() int {
	floor := StateTables[localID][2][1]
	return floor
}

func GetCurrentElevFloor(port string) int {
	floor := StateTables[port][2][1]
	return floor
}

func GetLocalID() string {
	return strconv.Itoa(StateTables[localID][0][1])
}

func Get() [7][3]int {
	return StateTables[localID]
}

func GetStateTables() map[string][7][3]int {
	return StateTables
}
