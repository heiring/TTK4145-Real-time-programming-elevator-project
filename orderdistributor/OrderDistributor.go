package orderdistributor

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"../config"
	"../elevio"
	"../tools"
)

// prioritizedOrders contains the orders awaiting execution by the local elevator, ordered by priority.
var prioritizedOrders = make([]int, 0)

func DistributeOrders(localID string, allOrders [4][3]int, allMovementDirection map[string]int, elevStatuses map[string][2]int) {
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

	curLocalFloor := elevStatuses[localID][0]
	localMovementDirection := allMovementDirection[localID]
	var nrOfHallDownOrders, nrOfHallUpOrders int
	for row := 0; row < 4; row++ {
		nrOfHallUpOrders += allOrders[row][0]
		nrOfHallDownOrders += allOrders[row][1]
	}

	for row := 0; row < 4; row++ {
		curHallUpOrder := allOrders[row][elevio.BT_HallUp]
		curHallDownOrder := allOrders[row][elevio.BT_HallDown]
		curCabOrder := allOrders[row][elevio.BT_Cab]
		orderDestination := row

		localOrderDir, _ := tools.DivCheck((orderDestination - curLocalFloor), (orderDestination - curLocalFloor))
		localDistance := int(math.Abs(float64(orderDestination - curLocalFloor)))
		var allDistances = make(map[string]int)
		for ID, status := range elevStatuses {
			location := status[0]
			allDistances[ID] = int(math.Abs(float64(orderDestination - location)))
		}

		// Hall buttons
		if (curHallUpOrder != 0 || curHallDownOrder != 0) && !tools.IntInSlice(orderDestination, prioritizedOrders) {
			var butnTypeDir int // +1 = Up, -1 = Down
			if curHallDownOrder != 0 {
				butnTypeDir = elevio.MD_Down
			} else {
				butnTypeDir = int(elevio.MD_Up)
			}

			takeOrder := false
			if len(allMovementDirection) > 1 {
				// fmt.Println("Multiple elevs!")
				// fmt.Println("Statuses: ", elevStatuses)
				// fmt.Println("Directions: ", allMovementDirection)
				for ID := range allMovementDirection {
					// fmt.Println("CURRENT ID: ", ID)
					if ID != localID {
						// fmt.Println("SCHWAM 1")
						isAlive := elevStatuses[ID][1]
						if isAlive == 1 {
							// fmt.Println("SCHWAM 2")
							externalOrderDir, _ := tools.DivCheck((orderDestination - elevStatuses[ID][0]), int(math.Abs(float64(orderDestination-elevStatuses[ID][0]))))
							if localOrderDir == localMovementDirection && externalOrderDir != allMovementDirection[ID] && (localOrderDir == butnTypeDir || localOrderDir == elevio.MD_Stop) {
								// If local only elev in same dir
								// fmt.Println("SCHWAM 3")
								takeOrder = true
							} else if localOrderDir == localMovementDirection && externalOrderDir == allMovementDirection[ID] && (localOrderDir == butnTypeDir || localOrderDir == elevio.MD_Stop) {
								// Else if local has shortest distance
								fmt.Println("SCHWAM 4")
								if localDistance <= allDistances[ID] {
									// fmt.Println("SCHWAM 5")
									takeOrder = true
									fmt.Println("true - 2")
								} else {
									// fmt.Println("SCHWAM 6")
									takeOrder = false
									fmt.Println("false - 0")
									break
								}
							} else if localOrderDir != localMovementDirection && externalOrderDir == allMovementDirection[ID] && externalOrderDir == butnTypeDir {
								// If other elevs is in same dir && local is not
								// fmt.Println("SCHWAM 7")
								takeOrder = false
								break
							} else if localOrderDir != localMovementDirection && externalOrderDir != allMovementDirection[ID] {
								// If no elevs in same dir
								// fmt.Println("SCHWAM 8")
								if localMovementDirection == elevio.MD_Stop && allMovementDirection[ID] != elevio.MD_Stop {
									// fmt.Println("SCHWAM 9")
									takeOrder = true
									fmt.Println("true - 3")
									// fmt.Println("Status: ", elevStatuses)
									// fmt.Println("Directions: ", allMovementDirection)
									// fmt.Println("My location: ", curLocalFloor, "\tHis location: ", elevStatuses[ID][0])
									// fmt.Println("My dir: ", localMovementDirection, "\tHis dir: ", allMovementDirection[ID])
									// fmt.Println("externalOrderDir: ", externalOrderDir, "\tbutnTypeDir: ", butnTypeDir)
									// CORRECT
								} else if localMovementDirection == elevio.MD_Stop && allMovementDirection[ID] == elevio.MD_Stop {
									// Else if local has shortest distance
									// fmt.Println("SCHWAM 10")
									if localDistance < allDistances[ID] {
										// fmt.Println("SCHWAM 11")
										takeOrder = true
										fmt.Println("true - 4")
										// fmt.Println(localID, " localDist: ", localDistance, "\n", ID, " externalDist: ", allDistances[ID])
									} else if localDistance == allDistances[ID] {
										localIDInt, _ := strconv.Atoi(localID)
										externalIDInt, _ := strconv.Atoi(ID)
										if localIDInt < externalIDInt {
											takeOrder = true
											fmt.Println("true - lucky")
										} else {
											takeOrder = false
											break
										}
									} else {
										// fmt.Println("SCHWAM 12")
										takeOrder = false
										// fmt.Println("false - 2, better choice: ", ID, "distances: ", localDistance, " - ", allDistances[ID])
										break
									}
								} else if localMovementDirection != elevio.MD_Stop && allMovementDirection[ID] == elevio.MD_Stop {
									// If other elevs in STOP
									// fmt.Println("SCHWAM 13")
									takeOrder = false
									break
								} else if localMovementDirection != elevio.MD_Stop && allMovementDirection[ID] != elevio.MD_Stop {
									// If no elevs in STOP
									// fmt.Println("SCHWAM 14")
									takeOrder = true
								} else {
									// fmt.Println("SCHWAM 15")
									fmt.Println("Order error 1!")
								}
								// } else if externalOrderDir == butnTypeDir {
								// 	// fmt.Println("SCHWAM 16")
								// 	fmt.Println("Order error 2!")
								// } else if localOrderDir == butnTypeDir {
								// 	// fmt.Println("SCHWAM 17")
								// 	fmt.Println("Order error 3!")
							} else if externalOrderDir != butnTypeDir {
								// fmt.Println("SCHWAM 18")
								if (nrOfHallDownOrders >= 1 && nrOfHallUpOrders == 0) || (nrOfHallDownOrders == 0 && nrOfHallUpOrders >= 1) {
									takeOrder = false
									break
								}

							}

						} else { // Elev is ded
							// make funeral
							// fmt.Println("SCHWAM 19")
							takeOrder = true
							fmt.Println("funeral prepared for ", ID)
						}
						// fmt.Println("SCHWAM 20")
					}
				}
			} else if len(allMovementDirection) == 1 {
				takeOrder = true
			}

			if takeOrder {
				if curHallDownOrder != 0 {
					addOrderToQueue(elevio.BT_HallDown, orderDestination, curLocalFloor, localMovementDirection)
				}
				if curHallUpOrder != 0 {
					addOrderToQueue(elevio.BT_HallUp, orderDestination, curLocalFloor, localMovementDirection)
				}
			}
		}

		// Cab buttons
		if curCabOrder != 0 && !tools.IntInSlice(orderDestination, prioritizedOrders) {
			addOrderToQueue(elevio.BT_Cab, orderDestination, curLocalFloor, localMovementDirection)
		}
	}

}

func addOrderToQueue(button elevio.ButtonType, orderDestination, curLocalFloor, localMovementDirection int) {
	// if curLocalFloor != orderDestination {
	if true {
		if len(prioritizedOrders) <= 0 {

			prioritizedOrders = append(prioritizedOrders, orderDestination)
			fmt.Println("APPENDED TO EMPTY: ", prioritizedOrders)
			//fmt.Println("NEW ORDER APPENDED (CAB): ", orderDestination)
		} else {
			for i, lastOrder := range prioritizedOrders {
				if !tools.IntInSlice(orderDestination, prioritizedOrders) {
					lastOrderDirection, _ := tools.DivCheck(lastOrder-curLocalFloor, int(math.Abs(float64(lastOrder-curLocalFloor))))
					curOrderDirection, _ := tools.DivCheck(orderDestination-curLocalFloor, int(math.Abs(float64(orderDestination-curLocalFloor))))

					// if lastOrder not in direction but neworder is
					if (localMovementDirection != elevio.MD_Stop) && (lastOrderDirection != localMovementDirection) && (curOrderDirection == localMovementDirection) {
						if !(button == elevio.BT_HallDown && curOrderDirection == int(elevio.MD_Up)) ||
							!(button == elevio.BT_HallUp && curOrderDirection == elevio.MD_Down) {
							// prioritizedOrders = append([]int{orderDestination}, prioritizedOrders...)
							appendToIndex(i, orderDestination)
							fmt.Println("Append (1) - ", orderDestination)
							break
						}
					}

					// if both orders in same dir and neworder closer than lastOrder
					if lastOrderDirection == curOrderDirection {
						newOrderDistance := int(math.Abs(float64(orderDestination - curLocalFloor)))
						lastOrderDistance := int(math.Abs(float64(lastOrder - curLocalFloor)))
						if newOrderDistance < lastOrderDistance {
							if button == elevio.BT_HallUp && curOrderDirection == int(elevio.MD_Up) {
								// prioritizedOrders = append([]int{orderDestination}, prioritizedOrders...)
								appendToIndex(i, orderDestination)
								fmt.Println("Append (2) - ", orderDestination)
								break
							} else if button == elevio.BT_HallDown && curOrderDirection == int(elevio.MD_Down) {
								// prioritizedOrders = append([]int{orderDestination}, prioritizedOrders...)
								appendToIndex(i, orderDestination)
								fmt.Println("Append (3) - ", orderDestination)
								break
							}
						}
					}

					// Give new order lowest priority
					if i == (len(prioritizedOrders) - 1) {
						prioritizedOrders = append(prioritizedOrders, orderDestination)
						fmt.Println("Append (4) - ", orderDestination)
					}
				}
			}
		}
	} else {
		fmt.Println("Already at floor ", orderDestination)
	}
}

func appendToIndex(index, floor int) {
	prioritizedOrders = append(prioritizedOrders, 0)
	copy(prioritizedOrders[(index+1):], prioritizedOrders[index:])
	prioritizedOrders[index] = floor
}

// PollOrders repeatedly checks for orders awaiting execution
func PollOrders(receiver chan<- int) {
	var prevOrder int
	init := true
	for {
		time.Sleep(config.PollRate)
		order := GetOrderFloor()
		if (order != prevOrder && order != -1) || (order == -1 && init) {
			receiver <- order
			if init {
				init = false
			}
		}
		prevOrder = order
	}
}

func RemoveOrder() {
	if len(prioritizedOrders) > 0 {
		fmt.Println("COMPLETED ORDER ", prioritizedOrders[0])
		prioritizedOrders = prioritizedOrders[1:]
		fmt.Println("Remaining ORDERS ", prioritizedOrders)
	}
}

func GetOrderFloor() int {
	if len(prioritizedOrders) > 0 {
		return prioritizedOrders[0]
	}
	return -1
}
