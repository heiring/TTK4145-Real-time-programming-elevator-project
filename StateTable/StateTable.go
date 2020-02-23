package StateTable

var stateTable [9][9] int

func initStateTable() {
	for row, cells := range stateTable {
		for _, col := range cells {
			stateTable[row][col] = 0
		}
    }
}

func calculateOrders() {
	
}

func UpdateStateTable(row, col, elev_nr, val int) {
	stateTable[row][col + elev_nr*3] = val

	// If button pressed: Calculate orders
}