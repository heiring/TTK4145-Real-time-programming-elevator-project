package statetable

import (
	"strconv"
)

var stateTable [7][9]int

func initStateTable() {
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
		distributeOrders()
	}
}

func UpdateElevLastFLoor(elev_nr, val int) {
	// val in [0, 3]
	// convert to binary [001, 100]
	binaryVal := strconv.FormatInt(int64(val), 2)
	for col := 0; col < 3; col++ {
		stateTable[2][col] = int(binaryVal[col])
	}

}

func ResetElevRow(row, elev_nr int) {
	for col := 0; col < 3; col++ {
		stateTable[row][col+elev_nr*3] = 0
	}
}
