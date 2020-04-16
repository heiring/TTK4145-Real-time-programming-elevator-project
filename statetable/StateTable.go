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

				updateLightsFromPacket()
				runOrderDistribution()
			}
		default:
			//do stuff
		}
	}
}
func updateLightsFromPacket() {
	// allOrders does not include new external orders
	allOrders, _, _ := GetSyncedOrders()
	fmt.Println("allOrders:\n", allOrders)
	for floor := range allOrders {
		for butn := elevio.BT_HallUp; butn < elevio.BT_Cab; butn++ {
			if allOrders[floor][butn] == 1 {
				elevio.SetButtonLamp(butn, floor, true)
				fmt.Println("ON - BTN LIGHT FROM PACKET")
			} else {
				elevio.SetButtonLamp(butn, floor, false)
				// fmt.Println("OFF - BTN LIGHT FROM PACKET")
			}
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
	allOrders, allDirections, allLocations := GetSyncedOrders()
	orderdistributor.DistributeOrders(string(localID), allOrders, allDirections, allLocations)
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

func GetSyncedOrders() ([4][3]int, map[string]int, map[string]int) {
	var allOrders [4][3]int
	var allDirections = make(map[string]int)
	var allLocations = make(map[string]int)
	for ID, statetable := range StateTables {
		isAlive := statetable[0][0]
		if isAlive == 0 {
			break
		}

		for row := 0; row < 4; row++ {
			for col := 0; col < 2; col++ {
				allOrders[row][col] += statetable[row+3][col]
				if statetable[row+3][col] != 0 {
					allOrders[row][col] = (allOrders[row][col] / allOrders[row][col])
				}
			}
			if ID == localID {
				allOrders[row][2] = statetable[row+3][2]
			}
		}
		allDirections[ID] = statetable[1][1]
		allLocations[ID] = statetable[1][2]
	}

	return allOrders, allDirections, allLocations
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
