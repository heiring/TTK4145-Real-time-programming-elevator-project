package statetable

import (
	"fmt"
	"strconv"
	"time"

	. "../config"
	"../orderdistributor"
)

// var stateTables [7][3]int
var StateTables = make(map[string][7][3]int)

// var stateTable [7][9]int
var localID string

const UnknownFloor int = -1

func InitStateTable(port int) {
	fmt.Println("InitStateTable")
	var tempStateTable [7][3]int
	for row, cells := range tempStateTable {
		for _, col := range cells {
			tempStateTable[row][col] = 0
			// UpdateStateTableIndex(row, col, 0, false)
		}
	}
	tempStateTable[0][0] = 1
	// Unknown starting position
	// stateTable[2][elevNr*3+1] = UnknownFloor
	// UpdateStateTableIndex(2, 1, UnknownFloor, false)
	tempStateTable[2][1] = UnknownFloor
	// Set ID = port
	// UpdateStateTableIndex(0, 1, port, false)
	tempStateTable[0][1] = port

	localID = strconv.Itoa(port)
	StateTables[localID] = tempStateTable
}

// func UpdateEntireStateTable(elevState ElevatorState) {
// 	for row, cells := range elevState.StateTable {
// 		for _, col := range cells {
// 			if !(row <= 2 && col == (elevNr*3+1)) {
// 				UpdateStateTableIndex(row, col, cells[col], false)
// 			}
// 		}
// 	}
// 	runOrderDistribution()
// }

func UpdateStateTableFromPacket(receiveStateCh <-chan ElevatorState) {
	for {
		select {
		case elevState := <-receiveStateCh:
			ID := elevState.ID
			if ID != localID {
				StateTables[ID] = elevState.StateTable
				runOrderDistribution()
			}

			// for row, cells := range elevState.StateTable {
			// 	fmt.Println("CELLS = ", cells)
			// 	for col := range cells {
			// 		if !(row <= 2 && col >= (elevNr*3) && col <= (elevNr*3+2)) {
			// 			// fmt.Printf("pre error, row: ")
			// 			// fmt.Println(row)
			// 			// fmt.Printf("col: ")
			// 			// fmt.Println(col)
			// 			stateTable[row][col] = cells[col]
			// 			fmt.Println("ROW = ", row, "\tCOL = ", col)
			// 			// fmt.Println("post error")
			// 		}
			// 	}
			// }
		default:
			//do stuff
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
		case activeElevators := <-activeElevatorsCh: //pakker kommer regelmessig
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
				// for index := 0; index < 3; index++ {
				// 	if ID == string(stateTable[0][1+index*3]) {
				// 		if isAlive {
				// 			UpdateStateTableIndex(0, (index * 3), 1, true)
				// 		} else {
				// 			UpdateStateTableIndex(0, (index * 3), 0, true)
				// 		}
				// 	}
				// }
			}
		default:
			//do stuff
		}
	}
	runOrderDistribution()
}

func UpdateStateTableIndex(row, col int, port string, val int, runDistribution bool) { // stateTableTransmitCh chan<- [7][9]int) {
	// stateTable[row][col+elevNr*3] = val
	statetable := StateTables[port]
	statetable[row][col] = val
	StateTables[port] = statetable
	if runDistribution {
		runOrderDistribution()
	}

}

func runOrderDistribution() {
	// var orders [4][3]int
	// for row := range orders {
	// 	for col := range orders[row] {
	// 		orders[row][col] = StateTables[3+row][col]
	// 	}
	// }
	// position := stateTable[2][elevNr*3+1]
	// direction := stateTable[1][elevNr*3+1]
	// orderdistributor.DistributeOrders(orders, position, direction)
	orderdistributor.DistributeOrders(string(localID), StateTables)
}

func UpdateElevLastFLoor(val int) {
	//fmt.Println("UpdateElevLastFLoor, val: ", val)
	// stateTable[2][elevNr*3+1] = val
	UpdateStateTableIndex(2, 1, localID, val, false)
}

// UpdateElevDirection comment
func UpdateElevDirection(val int) {
	// stateTable[1][elevNr*3+1] = val
	UpdateStateTableIndex(1, 1, localID, val, false)
}

func ResetElevRow(row int) {
	for col := 0; col < 3; col++ {
		UpdateStateTableIndex(row, col, localID, 0, false)
		// stateTable[row][col+elevNr*3] = 0
	}
}

func ResetRow(row int) {
	ResetElevRow(row)
}

func getPositionRow(port string) int {
	// var position [3]int
	// for i := range position {
	// 	position[i] = stateTable[2][i+3*elev_nr]
	// }
	position := StateTables[port][2][1]
	return position
}

func GetElevDirection(port string) int {
	// dir := stateTable[1][0+3*elev_nr]*(-1) + stateTable[1][1+3*elev_nr]*0 + stateTable[1][1+3*elev_nr]*1
	// check := stateTable[1][0+3*elev_nr] + stateTable[1][1+3*elev_nr] + stateTable[1][1+3*elev_nr]

	// if check == 1 {
	// 	return dir
	// }
	// return 11 // Error!
	direction := StateTables[port][1][1]
	return direction
}

func GetCurrentFloor() int {
	// bin := tools.ArrayToString(getPositionRow(elev_nr))
	// lastFloor, _ := strconv.ParseInt(bin, 2, 64)
	// return int(lastFloor)
	floor := StateTables[localID][2][1]
	return floor
}

func GetCurrentElevFloor(port string) int {
	// bin := tools.ArrayToString(getPositionRow(elev_nr))
	// lastFloor, _ := strconv.ParseInt(bin, 2, 64)
	// return int(lastFloor)
	floor := StateTables[port][2][1]
	return floor
}

func getCurrentID() string { //samme som GetLocalID()?
	return string(StateTables[localID][0][1])
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
