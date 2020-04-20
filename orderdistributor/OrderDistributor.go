package orderdistributor

import (
	"fmt"
	"math"
	"time"

	"../config"
	"../elevio"
	"../tools"
)

// prioritizedOrders contains the orders awaiting execution by the local elevator, ordered by priority.
var prioritizedOrders = make([]int, 0)

// DistributeOrders inputs the orders linked to hall lights and cab lights, a map indicating the movement directions of the different elevators, as well their respective positions and life statuses.
// This informations is used to update prioritizedOrders.
func DistributeOrders(localID string, allOrders [4][3]int, allElevDirections map[string]int, elevPositionsAndLifeStatuses map[string][2]int) {
	curLocalFloor := elevPositionsAndLifeStatuses[localID][0]
	localMovementDirection := allElevDirections[localID]

	for row := 0; row < 4; row++ {
		curHallUpOrder := allOrders[row][elevio.BT_HallUp]
		curHallDownOrder := allOrders[row][elevio.BT_HallDown]
		curCabOrder := allOrders[row][elevio.BT_Cab]
		orderDestination := row

		localOrderDir, _ := tools.DivCheck((orderDestination - curLocalFloor), (orderDestination - curLocalFloor))
		localDistance := int(math.Abs(float64(orderDestination - curLocalFloor)))
		var allDistances = make(map[string]int)
		for ID, status := range elevPositionsAndLifeStatuses {
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
			if len(allElevDirections) > 1 {
				for ID := range allElevDirections {
					if ID != localID {
						isAlive := elevPositionsAndLifeStatuses[ID][1]
						if isAlive == 1 {
							externalOrderDir, _ := tools.DivCheck((orderDestination - elevPositionsAndLifeStatuses[ID][0]), int(math.Abs(float64(orderDestination-elevPositionsAndLifeStatuses[ID][0]))))
							if localOrderDir == localMovementDirection && externalOrderDir != allElevDirections[ID] && (localOrderDir == butnTypeDir || localOrderDir == elevio.MD_Stop) {
								// If local only elev in same dir
								takeOrder = true
							} else if localOrderDir == localMovementDirection && externalOrderDir == allElevDirections[ID] && (localOrderDir == butnTypeDir || localOrderDir == elevio.MD_Stop) {
								// Else if local has shortest distance
								if localDistance <= allDistances[ID] {
									takeOrder = true
								}
							} else if localOrderDir != localMovementDirection && externalOrderDir == allElevDirections[ID] && externalOrderDir == butnTypeDir {
								// If other elevs is in same dir && local is not
								takeOrder = false
								break
							} else if localOrderDir != localMovementDirection && externalOrderDir != allElevDirections[ID] {
								// If no elevs in same dir
								if localMovementDirection == elevio.MD_Stop && allElevDirections[ID] != elevio.MD_Stop {
									takeOrder = true
									// CORRECT
								} else if localMovementDirection == elevio.MD_Stop && allElevDirections[ID] == elevio.MD_Stop {
									// Else if local has shortest distance
									if localDistance <= allDistances[ID] {
										takeOrder = true
									} else {
										takeOrder = false
										break
									}
								} else if localMovementDirection != elevio.MD_Stop && allElevDirections[ID] == elevio.MD_Stop {
									// If other elevs in STOP
									takeOrder = false
									break
								} else if localMovementDirection != elevio.MD_Stop && allElevDirections[ID] != elevio.MD_Stop {
									// If no elevs in STOP
									takeOrder = true
								} else {
									fmt.Println("Order error 1!")
								}
							} else if externalOrderDir == butnTypeDir {
								fmt.Println("Order error 2!")
							} else if localOrderDir == butnTypeDir {
								fmt.Println("Order error 3!")
							}
						} else { // Elevator is dead
							takeOrder = true

						}
					}
				}
			} else if len(allElevDirections) == 1 {
				takeOrder = true
			}

			if takeOrder {
				if curHallDownOrder != 0 {
					addOrderToPrioritizedOrders(elevio.BT_HallDown, orderDestination, curLocalFloor, localMovementDirection)
				}
				if curHallUpOrder != 0 {
					addOrderToPrioritizedOrders(elevio.BT_HallUp, orderDestination, curLocalFloor, localMovementDirection)
				}
			}
		}

		// Cab buttons
		if curCabOrder != 0 && !tools.IntInSlice(orderDestination, prioritizedOrders) {
			addOrderToPrioritizedOrders(elevio.BT_Cab, orderDestination, curLocalFloor, localMovementDirection)
		}
	}

}

func addOrderToPrioritizedOrders(button elevio.ButtonType, orderDestination, curLocalFloor, localMovementDirection int) {
	if len(prioritizedOrders) <= 0 {
		prioritizedOrders = append(prioritizedOrders, orderDestination)
	} else {
		for i, lastOrder := range prioritizedOrders {
			if !tools.IntInSlice(orderDestination, prioritizedOrders) {
				lastOrderDirection, _ := tools.DivCheck(lastOrder-curLocalFloor, int(math.Abs(float64(lastOrder-curLocalFloor))))
				curOrderDirection, _ := tools.DivCheck(orderDestination-curLocalFloor, int(math.Abs(float64(orderDestination-curLocalFloor))))

				// if lastOrder not in direction but neworder is
				if (localMovementDirection != elevio.MD_Stop) && (lastOrderDirection != localMovementDirection) && (curOrderDirection == localMovementDirection) {
					if !(button == elevio.BT_HallDown && curOrderDirection == int(elevio.MD_Up)) ||
						!(button == elevio.BT_HallUp && curOrderDirection == elevio.MD_Down) {
						appendToPrioritizedOrdersIndex(i, orderDestination)
						break
					}
				}

				// if both orders in same dir and neworder closer than lastOrder
				if lastOrderDirection == curOrderDirection {
					newOrderDistance := int(math.Abs(float64(orderDestination - curLocalFloor)))
					lastOrderDistance := int(math.Abs(float64(lastOrder - curLocalFloor)))
					if newOrderDistance < lastOrderDistance {
						if button == elevio.BT_HallUp && curOrderDirection == int(elevio.MD_Up) {
							appendToPrioritizedOrdersIndex(i, orderDestination)
							break
						} else if button == elevio.BT_HallDown && curOrderDirection == int(elevio.MD_Down) {
							appendToPrioritizedOrdersIndex(i, orderDestination)
							break
						}
					}
				}

				// Give new order lowest priority
				if i == (len(prioritizedOrders) - 1) {
					prioritizedOrders = append(prioritizedOrders, orderDestination)
				}
			}
		}
	}
}

func appendToPrioritizedOrdersIndex(index, floor int) {
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
		prioritizedOrders = prioritizedOrders[1:]
	}
}

func GetOrderFloor() int {
	if len(prioritizedOrders) > 0 {
		return prioritizedOrders[0]
	}
	return -1
}
