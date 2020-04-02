package statetable

import (
	"fmt"

	"../orderdistributor"
)

var stateTable [7][9]int

const UnknownFloor int = -1

func InitStateTable(elev_nr int) {
	fmt.Println("InitStateTable")
	for row, cells := range stateTable {
		for _, col := range cells {
			stateTable[row][col] = 0
		}
	}
	// Unknown starting position
	stateTable[2][elev_nr*3+1] = UnknownFloor

}

func UpdateStateTableIndex(row, col, elev_nr, val int) {
	stateTable[row][col+elev_nr*3] = val

	// If orders have changed calculate again
	if row > 2 && row < 7 {
		// To do: figure out slicing, current code is wrong. This is not python
		var orders [4][3]int
		for row := range orders {
			for col := range orders[row] {
				orders[row][col] = stateTable[3+row][col+3*elev_nr]
			}
		}
		position := stateTable[2][elev_nr*3+1]
		direction := stateTable[1][elev_nr*3+1]
		orderdistributor.DistributeOrders(orders, position, direction)
	}
}

func UpdateElevLastFLoor(elev_nr, val int) {
	// ResetElevRow(elev_nr, 2)
	fmt.Println("UpdateElevLastFLoor, val: ", val)
	stateTable[2][elev_nr*3+1] = val
	// // val in [0, 3]
	// // convert to binary [001, 100]
	// binaryVal := strconv.FormatInt(int64(val), 2)
	// for len(binaryVal) < 3 {
	// 	binaryVal = "0" + binaryVal
	// }
	// for col := 0; col < 3; col++ {
	// 	stateTable[2][elev_nr*3+col] = int(binaryVal[col])
	// }
	// fmt.Println("binaryVal: ", binaryVal)
}

func UpdateElevDirection(elev_nr, val int) {
	// ResetElevRow(elev_nr, 1)
	stateTable[1][elev_nr*3+1] = val
	// val = -1, 0, 1
	// if val == -1 {
	// 	stateTable[2][0] = 1
	// } else if val == 0 {
	// 	stateTable[2][1] = 1
	// } else if val == 1 {
	// 	stateTable[2][2] = 1
	// } else {
	// 	fmt.Println("ERROR! Could not update elev direction")
	// }

}

func ResetElevRow(elev_nr, row int) {
	for col := 0; col < 3; col++ {
		stateTable[row][col+elev_nr*3] = 0
	}
}

func ResetRow(row int) {
	for col := 0; col < 9; col++ {
		stateTable[row][col] = 0
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
