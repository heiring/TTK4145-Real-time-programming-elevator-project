package orderdistributor

import (
	"math"
	"time"

	"../config"
	"../elevio"
	"../tools"
)

var prioritizedOrders = make([]int, 0)

func DistributeOrders(localID string, stateTables map[string][7][3]int) {
	// Don't care which elev got the order
	// Active orders should not be moved unless the elev looses connection,
	// which should happen automatically, need to use the isAlive value
	// CAB buttons should only be taken by this elev

	// Loop through all orders(!)
	// Create an allorders table
	// - OR all hall orders into the table
	// - Only consider cab orders from this elev
	// Get all current elev locations and directions

	// ****** Discuss, is this necessary? ************************************************
	// If one elevator is not alive:
	// Remove all orders except the current one, in order to take extra orders as well
	// ***********************************************************************************

	var allOrders [4][3]int
	var allDirections = make(map[string]int)
	var allLocations = make(map[string]int)

	for ID, statetable := range stateTables {
		// isAlive := statetable[0][0]
		// if isAlive == 0 {
		// 	activeOrder := prioritizedOrders[0]
		// 	prioritizedOrders = nil
		// 	prioritizedOrders = append([]int{activeOrder}, prioritizedOrders...)
		// }

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

	curLocalFloor := allLocations[localID]
	curLocalDirection := allDirections[localID]

	for row := 0; row < 4; row++ {
		// curHallUpOrder := allOrders[row][elevio.BT_HallUp]
		// curHallDownOrder := allOrders[row][elevio.BT_HallDown]
		curCabOrder := allOrders[row][elevio.BT_Cab]
		orderDestination := row

		// Cab buttons
		if curCabOrder != 0 && !tools.IntInSlice(orderDestination, prioritizedOrders) {
			if len(prioritizedOrders) <= 0 {
				prioritizedOrders = append(prioritizedOrders, orderDestination)
				//fmt.Println("NEW ORDER APPENDED (CAB): ", orderDestination)
			} else {
				for i, lastOrder := range prioritizedOrders {
					if !tools.IntInSlice(orderDestination, prioritizedOrders) {
						lastOrderDirection, _ := tools.DivCheck(lastOrder-curLocalFloor, int(math.Abs(float64(lastOrder-curLocalFloor))))
						curCabOrderDirection, _ := tools.DivCheck(curCabOrder-curLocalFloor, int(math.Abs(float64(curCabOrder-curLocalFloor))))

						// if lastOrder not in direction but neworder is
						if (curLocalDirection != elevio.MD_Stop) && (lastOrderDirection != curLocalDirection) && (curCabOrderDirection == curLocalDirection) {
							prioritizedOrders = append([]int{orderDestination}, prioritizedOrders...)
							break
						}

						// if both orders in same dir and neworder closer than lastOrder
						if lastOrderDirection == curCabOrderDirection {
							newOrderDistance := int(math.Abs(float64(curLocalFloor - curCabOrder)))
							lastOrderDistance := int(math.Abs(float64(curLocalFloor - lastOrder)))
							if newOrderDistance < lastOrderDistance {
								prioritizedOrders = append([]int{orderDestination}, prioritizedOrders...)
							}
						}

						// Give new order lowest priority
						if i == (len(prioritizedOrders) - 1) {
							prioritizedOrders = append(prioritizedOrders, orderDestination)
						}
					}
				}
			}
			//fmt.Println("NEW PRIO (CAB): ", prioritizedOrders)
		}
	}

}

// func DistributeOrders(port string, stateTables map[string][7][3]int) {

// 	// for port, statetable := range statetable.StateTables {
// 	// Must be implemented whith cost function
// 	if true {
// 		sTable := stateTables[port]
// 		var orders [4][3]int
// 		for row := range orders {
// 			for col := range orders[row] {
// 				orders[row][col] = sTable[3+row][col]
// 			}
// 		}
// 		lastFloor := sTable[2][1]
// 		direction := sTable[1][1]
// 		for row := 0; row < 4; row++ {
// 			curHallUpOrder := orders[row][elevio.BT_HallUp]
// 			curHallDownOrder := orders[row][elevio.BT_HallDown]
// 			curCabOrder := orders[row][elevio.BT_Cab]
// 			orderDestination := row

// 			// Hall buttons
// 			// - To do

// 			// HallUp buttons
// 			if curHallUpOrder != 0 && !tools.IntInSlice(orderDestination, prioritizedOrders) {
// 				prioritizedOrders = append(prioritizedOrders, orderDestination)
// 				//fmt.Println("NEW PRIO (HU): ", prioritizedOrders)
// 			}

// 			// HallDown buttons
// 			if curHallDownOrder != 0 && !tools.IntInSlice(orderDestination, prioritizedOrders) {
// 				prioritizedOrders = append(prioritizedOrders, orderDestination)
// 				//fmt.Println("NEW PRIO (HD): ", prioritizedOrders)
// 				// if len(prioritizedOrders) <= 0 {
// 				// 	fmt.Println("0 len")
// 				// 	prioritizedOrders = append(prioritizedOrders, curHallDownOrder)
// 				// 	fmt.Println("NEW ORDER APPENDED (HD): ", curHallDownOrder)
// 				// } else {
// 				// 	for index, lastOrder := range prioritizedOrders {
// 				// 		lastOrderDirection := (lastOrder - lastFloor) / int(math.Abs(float64(lastOrder-lastFloor)))
// 				// 		curHallOrderDirection := (curHallDownOrder - lastFloor) / int(math.Abs(float64(curCabOrder-lastFloor)))

// 				// 		// if direction not down
// 				// 		// Add before lastOrder if closer

// 				// 		// if lastOrder not in direction but neworder is
// 				// 		if (direction != elevio.MD_Stop) && (lastOrderDirection != direction) && (curCabOrderDirection == direction) {
// 				// 			prioritizedOrders = append([]int{curCabOrder}, prioritizedOrders...)
// 				// 			break
// 				// 		}

// 				// 		// if both orders in same dir and neworder closer than lastOrder
// 				// 		if lastOrderDirection == curCabOrderDirection {
// 				// 			newOrderDistance := int(math.Abs(float64(lastFloor - curCabOrder)))
// 				// 			lastOrderDistance := int(math.Abs(float64(lastFloor - lastOrder)))
// 				// 			if newOrderDistance < lastOrderDistance {
// 				// 				prioritizedOrders = append([]int{curCabOrder}, prioritizedOrders...)
// 				// 			}
// 				// 		}

// 				// 		// Give new order lowest priority
// 				// 		if i == (len(prioritizedOrders) - 1) {
// 				// 			prioritizedOrders = append(prioritizedOrders, curCabOrder)
// 				// 		}
// 				// 	}
// 				// }
// 			}

// 			// Cab buttons
// 			if curCabOrder != 0 && !tools.IntInSlice(orderDestination, prioritizedOrders) {
// 				if len(prioritizedOrders) <= 0 {
// 					prioritizedOrders = append(prioritizedOrders, orderDestination)
// 					//fmt.Println("NEW ORDER APPENDED (CAB): ", orderDestination)
// 				} else {
// 					for i, lastOrder := range prioritizedOrders {
// 						if !tools.IntInSlice(orderDestination, prioritizedOrders) {
// 							lastOrderDirection, _ := tools.DivCheck(lastOrder-lastFloor, int(math.Abs(float64(lastOrder-lastFloor))))
// 							curCabOrderDirection, _ := tools.DivCheck(curCabOrder-lastFloor, int(math.Abs(float64(curCabOrder-lastFloor))))

// 							// if lastOrder not in direction but neworder is
// 							if (direction != elevio.MD_Stop) && (lastOrderDirection != direction) && (curCabOrderDirection == direction) {
// 								prioritizedOrders = append([]int{orderDestination}, prioritizedOrders...)
// 								break
// 							}

// 							// if both orders in same dir and neworder closer than lastOrder
// 							if lastOrderDirection == curCabOrderDirection {
// 								newOrderDistance := int(math.Abs(float64(lastFloor - curCabOrder)))
// 								lastOrderDistance := int(math.Abs(float64(lastFloor - lastOrder)))
// 								if newOrderDistance < lastOrderDistance {
// 									prioritizedOrders = append([]int{orderDestination}, prioritizedOrders...)
// 								}
// 							}

// 							// Give new order lowest priority
// 							if i == (len(prioritizedOrders) - 1) {
// 								prioritizedOrders = append(prioritizedOrders, orderDestination)
// 							}
// 						}
// 					}
// 				}
// 				//fmt.Println("NEW PRIO (CAB): ", prioritizedOrders)
// 			}
// 		}
// 	}

// }

func PollOrders(receiver chan<- int) {
	var prevOrder int
	init := true
	for {
		time.Sleep(config.PollRate)
		order := GetOrderFloor()
		if (order != prevOrder && order != -1) || (order == -1 && init) {
			//fmt.Println("Sending order...")
			receiver <- order
			if init {
				init = false
			}
		}
		prevOrder = order
	}
}

func CompleteCurrentOrder() {
	//fmt.Println("REMOVED COMPLETED ORDER")
	//fmt.Println("OLD PRIO.: ", prioritizedOrders)
	prioritizedOrders = prioritizedOrders[1:]
	//fmt.Println("NEW PRIO.: ", prioritizedOrders)
}

func GetOrderFloor() int {
	if len(prioritizedOrders) > 0 {
		return prioritizedOrders[0]
	}
	return -1
}
