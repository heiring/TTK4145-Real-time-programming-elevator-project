package statetable

import (
	"fmt"
	"strconv"

	"../orderdistributor"
	"../tools"
)

var stateTable [7][9]int

func InitStateTable() {
	fmt.Println("InitStateTable")
	for row, cells := range stateTable {
		for _, col := range cells {
			stateTable[row][col] = 0
		}
	}
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
		var position [3]int
		for i := range position {
			position[i] = stateTable[2][i+3*elev_nr]
		}
		orderdistributor.DistributeOrders(orders, position)
	}
}

func UpdateElevLastFLoor(elev_nr, val int) {
	ResetElevRow(elev_nr, 2)
	fmt.Println("UpdateElevLastFLoor, val: ", val)
	// val in [0, 3]
	// convert to binary [001, 100]
	binaryVal := strconv.FormatInt(int64(val), 2)
	for len(binaryVal) < 3 {
		binaryVal = "0" + binaryVal
	}
	for col := 0; col < 3; col++ {
		stateTable[2][col] = int(binaryVal[col])
	}
}

func UpdateElevDirection(elev_nr, val int) {
	ResetElevRow(elev_nr, 1)
	// val = -1, 0, 1
	if val == -1 {
		stateTable[2][0] = 1
	} else if val == 0 {
		stateTable[2][1] = 1
	} else if val == 1 {
		stateTable[2][2] = 1
	} else {
		fmt.Println("ERROR! Could not update elev direction")
	}

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

func getPositionRow(elev_nr int) [3]int {
	var position [3]int
	for i := range position {
		position[i] = stateTable[2][i+3*elev_nr]
	}
	return position
}

func GetElevDirection(elev_nr int) int {
	dir := stateTable[1][0+3*elev_nr]*(-1) + stateTable[1][1+3*elev_nr]*0 + stateTable[1][1+3*elev_nr]*1
	check := stateTable[1][0+3*elev_nr] + stateTable[1][1+3*elev_nr] + stateTable[1][1+3*elev_nr]

	if check == 1 {
		return dir
	}
	return 11 // Error!
}

func GetCurrentFloor(elev_nr int) int {
	bin := tools.ArrayToString(getPositionRow(elev_nr))
	lastFloor, _ := strconv.ParseInt(bin, 0, 64)
	return int(lastFloor)
}
