package orderdistributor

import (
	"fmt"
	"math"
	"strconv"

	"../elevio"
	"../tools"
)

var prioritizedOrders = make([]int, 0)

func DistributeOrders(orders [4][3]int, positions [3]int) {
	fmt.Println("DistributeOrders")
	fmt.Println("prioritizedOrders: ", prioritizedOrders)
	fmt.Println("len: ", len(prioritizedOrders))
	fmt.Println("orders: ", orders)
	direction, lastFloor := getDirectionAndFloor(positions)

	for row := 0; row < 4; row++ {
		// Hall buttons
		// - To do

		// Cab buttons
		curCabOrder := orders[row][elevio.BT_Cab]
		if orders[row][elevio.BT_Cab] != 0 {
			if len(prioritizedOrders) <= 0 {
				fmt.Println("0 len")

				prioritizedOrders = append(prioritizedOrders, orders[row][elevio.BT_Cab])
				fmt.Println("APPENDED!")
				fmt.Println("prioritizedOrders appended 0: ", prioritizedOrders)
				fmt.Println("orders[row][elevio.BT_Cab] ", orders[row][elevio.BT_Cab], row, elevio.BT_Cab)
			} else {
				for i, lastOrder := range prioritizedOrders {
					lastOrderDirection := (lastOrder - lastFloor) / int(math.Abs(float64(lastOrder-lastFloor)))
					curCabOrderDirection := (curCabOrder - lastFloor) / int(math.Abs(float64(curCabOrder-lastFloor)))

					// if lastOrder not in direction but neworder is
					if (direction != elevio.MD_Stop) && (lastOrderDirection != direction) && (curCabOrderDirection == direction) {
						prioritizedOrders = append([]int{curCabOrder}, prioritizedOrders...)
						break
					}

					// if both orders in same dir and neworder closer than lastOrder
					if lastOrderDirection == curCabOrderDirection {
						newOrderDistance := int(math.Abs(float64(lastFloor - curCabOrder)))
						lastOrderDistance := int(math.Abs(float64(lastFloor - lastOrder)))
						if newOrderDistance < lastOrderDistance {
							prioritizedOrders = append([]int{curCabOrder}, prioritizedOrders...)
						}
					}

					// Give new order lowest priority
					if i == (len(prioritizedOrders) - 1) {
						prioritizedOrders = append(prioritizedOrders, curCabOrder)
					}
				}
			}
		}
	}

	
}

func CompleteCurrentOrder() {
	prioritizedOrders = prioritizedOrders[1:]
}

func getDirectionAndFloor(positions [3]int) (int, int) {
	bin := tools.ArrayToString(positions)
	lastFloor, _ := strconv.ParseInt(bin, 0, 64)
	direction := positions[0]*(-1) + positions[1]*(0) + positions[2]*(1)
	if direction == 0 && positions[1] == 0 {
		direction = 11
	}

	return direction, int(lastFloor)
}

func GetCurrentOrder() int {
	if len(prioritizedOrders) > 0 {
		return prioritizedOrders[0]
	}
	return -1
}
