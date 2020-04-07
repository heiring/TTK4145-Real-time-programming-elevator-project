package orderdistributor

import (
	"fmt"
	"math"
	"time"

	"../config"
	"../elevio"
	"../tools"
)

var prioritizedOrders = make([]int, 0)

func DistributeOrders(orders [4][3]int, lastFloor, direction int) {

	for row := 0; row < 4; row++ {
		curHallUpOrder := orders[row][elevio.BT_HallUp]
		curHallDownOrder := orders[row][elevio.BT_HallDown]
		curCabOrder := orders[row][elevio.BT_Cab]
		orderDestination := row

		// Hall buttons
		// - To do

		// HallUp buttons
		if curHallUpOrder != 0 && !tools.IntInSlice(orderDestination, prioritizedOrders) {
			prioritizedOrders = append(prioritizedOrders, orderDestination)
			fmt.Println("NEW PRIO (HU): ", prioritizedOrders)
		}

		// HallDown buttons
		if curHallDownOrder != 0 && !tools.IntInSlice(orderDestination, prioritizedOrders) {
			prioritizedOrders = append(prioritizedOrders, orderDestination)
			fmt.Println("NEW PRIO (HD): ", prioritizedOrders)
			// if len(prioritizedOrders) <= 0 {
			// 	fmt.Println("0 len")
			// 	prioritizedOrders = append(prioritizedOrders, curHallDownOrder)
			// 	fmt.Println("NEW ORDER APPENDED (HD): ", curHallDownOrder)
			// } else {
			// 	for index, lastOrder := range prioritizedOrders {
			// 		lastOrderDirection := (lastOrder - lastFloor) / int(math.Abs(float64(lastOrder-lastFloor)))
			// 		curHallOrderDirection := (curHallDownOrder - lastFloor) / int(math.Abs(float64(curCabOrder-lastFloor)))

			// 		// if direction not down
			// 		// Add before lastOrder if closer

			// 		// if lastOrder not in direction but neworder is
			// 		if (direction != elevio.MD_Stop) && (lastOrderDirection != direction) && (curCabOrderDirection == direction) {
			// 			prioritizedOrders = append([]int{curCabOrder}, prioritizedOrders...)
			// 			break
			// 		}

			// 		// if both orders in same dir and neworder closer than lastOrder
			// 		if lastOrderDirection == curCabOrderDirection {
			// 			newOrderDistance := int(math.Abs(float64(lastFloor - curCabOrder)))
			// 			lastOrderDistance := int(math.Abs(float64(lastFloor - lastOrder)))
			// 			if newOrderDistance < lastOrderDistance {
			// 				prioritizedOrders = append([]int{curCabOrder}, prioritizedOrders...)
			// 			}
			// 		}

			// 		// Give new order lowest priority
			// 		if i == (len(prioritizedOrders) - 1) {
			// 			prioritizedOrders = append(prioritizedOrders, curCabOrder)
			// 		}
			// 	}
			// }
		}

		// Cab buttons
		if curCabOrder != 0 && !tools.IntInSlice(orderDestination, prioritizedOrders) {
			if len(prioritizedOrders) <= 0 {
				prioritizedOrders = append(prioritizedOrders, orderDestination)
				fmt.Println("NEW ORDER APPENDED (CAB): ", orderDestination)
			} else {
				for i, lastOrder := range prioritizedOrders {
					if !tools.IntInSlice(orderDestination, prioritizedOrders) {
						lastOrderDirection, _ := tools.DivCheck(lastOrder-lastFloor, int(math.Abs(float64(lastOrder-lastFloor))))
						curCabOrderDirection, _ := tools.DivCheck(curCabOrder-lastFloor, int(math.Abs(float64(curCabOrder-lastFloor))))

						// if lastOrder not in direction but neworder is
						if (direction != elevio.MD_Stop) && (lastOrderDirection != direction) && (curCabOrderDirection == direction) {
							prioritizedOrders = append([]int{orderDestination}, prioritizedOrders...)
							break
						}

						// if both orders in same dir and neworder closer than lastOrder
						if lastOrderDirection == curCabOrderDirection {
							newOrderDistance := int(math.Abs(float64(lastFloor - curCabOrder)))
							lastOrderDistance := int(math.Abs(float64(lastFloor - lastOrder)))
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
			fmt.Println("NEW PRIO (CAB): ", prioritizedOrders)
		}
	}

}

func PollOrders(receiver chan<- int) {
	var prevOrder int
	init := true
	for {
		time.Sleep(config.PollRate)
		order := GetOrderFloor()
		if (order != prevOrder && order != -1) || (order == -1 && init) {
			fmt.Println("Sending order...")
			receiver <- order
			if init {
				init = false
			}
		}
		prevOrder = order
	}
}

func CompleteCurrentOrder() {
	fmt.Println("REMOVED COMPLETED ORDER")
	fmt.Println("OLD PRIO.: ", prioritizedOrders)
	prioritizedOrders = prioritizedOrders[1:]
	fmt.Println("NEW PRIO.: ", prioritizedOrders)
}

func GetOrderFloor() int {
	if len(prioritizedOrders) > 0 {
		return prioritizedOrders[0]
	}
	return -1
}
