package orderdistributor

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"../elevio"
)

var prioritizedOrders = make([]int, 0)

func distributeOrders(orders, positions [][]int) {
	direction, lastFloor := getDirectionAndFloor(positions)

	for row := 0; row < 4; row++ {
		// Hall buttons
		// - To do

		// Cab buttons
		curCabOrder := orders[row][elevio.BT_Cab]
		if len(prioritizedOrders) <= 0 {
			prioritizedOrders = append(prioritizedOrders, orders[row][elevio.BT_Cab])
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

func getDirectionAndFloor(positions [][]int) (int, int) {
	bin := arrayToString(positions[1])
	lastFloor, _ := strconv.ParseInt(bin, 0, 64)
	direction := positions[0][0]*(-1) + positions[0][1]*(0) + positions[0][2]*(1)
	if direction == 0 && positions[0][1] == 0 {
		direction = 11
	}

	return direction, int(lastFloor)
}

func arrayToString(a []int) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", "", -1), "[]")
}
