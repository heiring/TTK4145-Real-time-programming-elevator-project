package statetable

import (
	"fmt"
	"strconv"
	"time"

	. "../config"
	"../elevio"
	"../orderdistributor"
)

var StateTables *StateTablesMutex

var activeLights *ActiveLightsMutex

var localID string

const UnknownFloor int = -1

func InitStateTable(port int) {
	fmt.Println("InitStateTable")
	var tempStateTable [7][3]int
	for row, cells := range tempStateTable {
		for _, col := range cells {
			tempStateTable[row][col] = 0
		}
	}
	// Set status to active
	tempStateTable[0][0] = 1
	// Unknown starting position
	tempStateTable[2][1] = UnknownFloor
	// Set ID = port
	tempStateTable[0][1] = port

	localID = strconv.Itoa(port)

	//StateTables[localID] = tempStateTable

	StateTables = &StateTablesMutex{
		Internal: map[string][7][3]int{
			localID: tempStateTable,
		},
	}

	activeLights = &ActiveLightsMutex{Internal: map[[2]int]bool{}} //dette passer kanskje ikke i InitStateTable() ?

}

func UpdateStateTableFromPacket(receiveStateCh <-chan ElevatorState, stateTableTransmitCh chan [7][3]int) {
	for {
		select {
		case elevState := <-receiveStateCh:
			ID := elevState.ID
			if ID != localID {
				StateTables.Write(ID, elevState.StateTable)
				updateHallLightsFromExternalOrders()
				runOrderDistribution()
				updatedLocalState, ok := checkIfExternalOrderCompleted(elevState.StateTable)
				if ok {
					StateTables.Write(localID, updatedLocalState)
					stateTableTransmitCh <- Get()
				}
			}
		default:
			//do stuff
		}
	}
}

func checkIfExternalOrderCompleted(elevState [7][3]int) ([7][3]int, bool) {
	positionFloor := elevState[2][1]
	elevDirection := elevState[1][1]
	localStateTable := Get()
	updateLocal := false
	for col := 0; col < 2; col++ {
		if (localStateTable[3+positionFloor][col] == 1) && !(col == 0 && elevDirection == int(elevio.MD_Down)) && !(col == 1 && elevDirection == int(elevio.MD_Up)) {
			// Hall up/down with reversed elev dir will result in order deleted to soon
			localStateTable[3+positionFloor][col] = 0
			updateLocal = true
			// We also need to notify the external elev that the order is completed now
		}
	}
	return localStateTable, updateLocal
}

func updateHallLightsFromExternalOrders() {
	allOrders, _, _ := GetSyncedOrders()
	activeLightsUpdate := activeLights.ReadWholeMap()
	for floor := range allOrders {
		for butn := elevio.BT_HallUp; butn < elevio.BT_Cab; butn++ {
			if allOrders[floor][butn] == 1 {
				if !activeLightsUpdate[[2]int{int(butn), floor}] {
					elevio.SetButtonLamp(butn, floor, true)

					activeLightsUpdate[[2]int{int(butn), floor}] = true

					fmt.Println("ON - BTN LIGHT FROM PACKET")
				}
			} else {
				if activeLightsUpdate[[2]int{int(butn), floor}] {
					elevio.SetButtonLamp(butn, floor, false)

					activeLightsUpdate[[2]int{int(butn), floor}] = false

					// fmt.Println("OFF - BTN LIGHT FROM PACKET")
				}
			}

		}
	}
	activeLights.WriteWholeMap(activeLightsUpdate)
}

func TransmitState(stateTableTransmitCh <-chan [7][3]int, transmitStateCh chan<- ElevatorState) {
	ticker := time.NewTicker(StateTransmissionInterval)
	//stateTable := StateTables[localID]
	stateTable := ReadStateTable(localID)
	elevatorState := ElevatorState{ID: localID, StateTable: stateTable}
	for {
		select {
		case stateTable = <-stateTableTransmitCh:
			elevatorState.StateTable = stateTable
		case <-ticker.C:
			transmitStateCh <- elevatorState
		default:
			//do nothing
		}
	}
}

func UpdateActiveElevators(activeElevatorsCh <-chan map[string]bool) {
	for {
		select {
		case activeElevators := <-activeElevatorsCh:
			for ID, isAlive := range activeElevators {
				stateTableUpdate := ReadStateTable(ID)

				if isAlive {
					if stateTableUpdate[0][0] == 0 {
						//UpdateStateTableIndex(0, 0, ID, 1, true)
						stateTableUpdate[0][0] = 1
						StateTables.Write(ID, stateTableUpdate)
						runOrderDistribution()
					}
				} else {
					if stateTableUpdate[0][0] == 1 {
						//UpdateStateTableIndex(0, 0, ID, 0, true)
						stateTableUpdate[0][0] = 0
						StateTables.Write(ID, stateTableUpdate)
						runOrderDistribution()
						fmt.Println("DANGER")
					}
				}
			}
		default:
			//do stuff
		}
	}
}

func UpdateStateTableIndex(row, col int, port string, val int, runDistribution bool) { // stateTableTransmitCh chan<- [7][9]int) {
	stateTable := ReadStateTable(port)
	stateTable[row][col] = val
	StateTables.Write(port, stateTable)
	if runDistribution {
		runOrderDistribution()
	}

	/*statetable, ok := StateTables.Read(localID)
	if !ok {
		fmt.Println("read error")
	}statetable
		if runDistribution {
			runOrderDistribution()
		}
	*/
}

func runOrderDistribution() {
	allOrders, allDirections, elevStatuses := GetSyncedOrders()
	orderdistributor.DistributeOrders(string(localID), allOrders, allDirections, elevStatuses) //string(localID)
}

func UpdateElevLastFLoor(val int) {
	UpdateStateTableIndex(2, 1, localID, val, false)
}

func UpdateElevDirection(val int) {
	UpdateStateTableIndex(1, 1, localID, val, false)
}

func ResetElevRow(row int, ID string) {
	for col := 0; col < 3; col++ {
		UpdateStateTableIndex(row, col, ID, 0, false)
	}
}

func ResetRow(row int) {
	stateTables := StateTables.ReadWholeMap()
	for ID, _ := range stateTables {
		ResetElevRow(row, ID)
	}
}

func getPositionRow(port string) int {
	stateTable := ReadStateTable(port)

	position := stateTable[2][1]
	return position

	//position := StateTables[port][2][1]
	//return position
}

func GetSyncedOrders() ([4][3]int, map[string]int, map[string][2]int) { //omdøpe til noe som SyncOrdersDirectionsLocations (positions?)
	var allOrders [4][3]int
	var allDirections = make(map[string]int)
	var elevStatuses = make(map[string][2]int)
	stateTables := StateTables.ReadWholeMap()
	for ID, statetable := range stateTables {
		var status [2]int
		isAlive := statetable[0][0]
		status[1] = isAlive

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
		status[0] = statetable[2][1]
		elevStatuses[ID] = status
	}
	return allOrders, allDirections, elevStatuses
}

func GetElevDirection(port string) int {
	stateTable := ReadStateTable(port)

	direction := stateTable[1][1]
	return direction

	//direction := StateTables[port][1][1]
	//return direction
}

func GetCurrentFloor() int {
	stateTable := ReadStateTable(localID)

	floor := stateTable[2][1]
	return floor

	//floor := StateTables[localID][2][1]
	//return floor
}

func GetCurrentElevFloor(port string) int {
	stateTable := ReadStateTable(port)

	floor := stateTable[2][1]
	return floor

	//floor := StateTables[port][2][1]
	//return floor
}

func GetLocalID() string {
	//stateTable := ReadStateTable(localID)
	//return strconv.Itoa(stateTable[0][1])

	return localID
}

func Get() [7][3]int {
	stateTable := ReadStateTable(localID)
	return stateTable

	//return StateTables[localID]
}

func GetStateTables() map[string][7][3]int {
	return StateTables.ReadWholeMap()

	//return StateTables
}
func ReadStateTable(ID string) [7][3]int {
	stateTable, ok := StateTables.Read(ID)
	if !ok {
		fmt.Println("read error")
	}
	return stateTable
}

func UpdateActiveLights(butn elevio.ButtonType, floor int, active bool) {
	activeLights.Write(int(butn), floor, active)
}

func checkActiveLights(butn elevio.ButtonType, floor int) bool {
	isActive, ok := activeLights.Read(int(butn), floor)
	if !ok {
		fmt.Println("active lights: read error")
	}
	return isActive
}
