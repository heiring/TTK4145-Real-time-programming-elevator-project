package statetable

import (
	"fmt"

	. "../config"
	"../orderdistributor"
)

var stateTable [7][9]int
var elevNr int

const UnknownFloor int = -1

func InitStateTable(elevnr, port int) {
	elevNr = elevnr
	fmt.Println("InitStateTable")
	for row, cells := range stateTable {
		for _, col := range cells {
			// stateTable[row][col] = 0
			UpdateStateTableIndex(row, col, 0, false)
		}
	}
	// Unknown starting position
	// stateTable[2][elevNr*3+1] = UnknownFloor
	UpdateStateTableIndex(2, 1, UnknownFloor, false)
	// Set ID = port
	UpdateStateTableIndex(0, 1, port, false)
}

func UpdateEntireStateTable(elevState ElevatorState) {
	for row, cells := range elevState.StateTable {
		for _, col := range cells {
			if !(row <= 2 && col == (elevNr*3+1)) {
				UpdateStateTableIndex(row, col, cells[col], false)
			}
		}
	}
	runOrderDistribution()
}

func UpdateStateTableFromPacket(receiveStateCh <-chan ElevatorState) {
	for {
		select {
		case elevState := <-receiveStateCh:
			for row, cells := range elevState.StateTable {
				for _, col := range cells {
					if !(row <= 2 && col == (elevNr*3+1)) {
						stateTable[row][col+elevNr*3] = cells[col]
					}
				}
			}
		default:
			//do stuff
		}
	}
	runOrderDistribution()
}

func TransmitState(stateTableTransmitCh <-chan [7][9]int, ID string, transmitStateCh chan<- ElevatorState) {
	for {
		select {
		case stateTable := <-stateTableTransmitCh:
			elevatorState := ElevatorState{ID: ID, StateTable: stateTable}
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

		default:
			//do nothing
		}
	}
	runOrderDistribution()
}

func UpdateStateTableIndex(row, col, val int, runDistribution bool) { // stateTableTransmitCh chan<- [7][9]int) {
	stateTable[row][col+elevNr*3] = val
	if runDistribution {
		runOrderDistribution()
		//stateTableTransmitCh <- stateTable

	}

}

func runOrderDistribution() {
	var orders [4][3]int
	for row := range orders {
		for col := range orders[row] {
			orders[row][col] = stateTable[3+row][col+3*elevNr]
		}
	}
	position := stateTable[2][elevNr*3+1]
	direction := stateTable[1][elevNr*3+1]
	orderdistributor.DistributeOrders(orders, position, direction)
}

func UpdateElevLastFLoor(val int) {
	fmt.Println("UpdateElevLastFLoor, val: ", val)
	// stateTable[2][elevNr*3+1] = val
	UpdateStateTableIndex(2, 1, val, false)
}

// UpdateElevDirection comment
func UpdateElevDirection(val int) {
	// stateTable[1][elevNr*3+1] = val
	UpdateStateTableIndex(1, 1, val, false)
}

func ResetElevRow(row int) {
	for col := 0; col < 3; col++ {
		UpdateStateTableIndex(row, col, 0, false)
		// stateTable[row][col+elevNr*3] = 0
	}
}

func ResetRow(row int) {
	for col := 0; col < 9; col++ {
		UpdateStateTableIndex(row, col, 0, false)
		// stateTable[row][col] = 0
	}
}

func getPositionRow(elev_nr int) int {
	// var position [3]int
	// for i := range position {
	// 	position[i] = stateTable[2][i+3*elev_nr]
	// }
	position := stateTable[2][elev_nr*3+1]
	return position
}

func GetElevDirection(elev_nr int) int {
	// dir := stateTable[1][0+3*elev_nr]*(-1) + stateTable[1][1+3*elev_nr]*0 + stateTable[1][1+3*elev_nr]*1
	// check := stateTable[1][0+3*elev_nr] + stateTable[1][1+3*elev_nr] + stateTable[1][1+3*elev_nr]

	// if check == 1 {
	// 	return dir
	// }
	// return 11 // Error!
	direction := stateTable[1][elev_nr*3+1]
	return direction
}

func GetCurrentFloor(elev_nr int) int {
	// bin := tools.ArrayToString(getPositionRow(elev_nr))
	// lastFloor, _ := strconv.ParseInt(bin, 2, 64)
	// return int(lastFloor)
	floor := stateTable[2][elev_nr*3+1]
	return floor
}

func getCurrentID() string {
	return string(stateTable[0][elevNr*3+1])
}
