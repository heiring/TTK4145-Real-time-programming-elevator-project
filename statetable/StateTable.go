package statetable

import (
	"fmt"
	"strconv"
	"time"

	"../config"
	"../elevio"
	"../mapsync"
	"../orderdistributor"
)

// StateTables is a map with the elevator IDs as keys and their respective state table as the value.
// A state table describes the correseponding elevators state (last motor direction, last floor visited, ID, network online/offline status, motor functionality status and buttons pressed)
//
// A state table has the following form:
//	-----------------------------------------
//	| Online	|	ID		| Motor working	|
//	| x			| Direction	|		x		|
//	| x			| Position	|		x		|
//	| Hall Up	| Hall Down | 		Cab		| Floor 1
//	| Hall Up	| Hall Down | 		Cab		| Floor 2
//	| Hall Up	| Hall Down | 		Cab		| Floor 3
//	| Hall Up	| Hall Down | 		Cab		| Floor 4
//	-----------------------------------------
//  x - not in use
var StateTables *mapsync.StateTablesSync

// Active lights is a map with the keys being slices with button type and floor as elements, and the values being bools indicating if the lights are shut on or off.
var activeLights *mapsync.ActiveLightsSync
var localID string

const UnknownFloor int = -1

// ElevatorState is the data type which will be broadcasted from each elevator
type ElevatorState struct {
	ID         string
	StateTable [7][3]int
}

// InitStateTable initialized the local state table with appropriate values.
func InitStateTable(port int) {
	var tempStateTable [7][3]int
	for row, cells := range tempStateTable {
		for _, col := range cells {
			tempStateTable[row][col] = 0
		}
	}
	// Set status to active
	tempStateTable[0][0] = 1
	// motor functional
	tempStateTable[0][2] = 1
	// Unknown starting position
	tempStateTable[2][1] = UnknownFloor
	// Set ID = port
	tempStateTable[0][1] = port

	localID = strconv.Itoa(port)

	StateTables = &mapsync.StateTablesSync{
		StateTables: map[string][7][3]int{
			localID: tempStateTable,
		},
	}

	activeLights = &mapsync.ActiveLightsSync{ActiveLights: map[[2]int]bool{}}
}

//UpdateStateTableFromPacket receives the state table transmitted from other elevators and updates their information in StateTables. Hall buttons are also synchronized.
func UpdateStateTableFromPacket(receiveStateCh <-chan ElevatorState, stateTableTransmitCh chan [7][3]int) {
	recoveryCompleted := false
	for {
		select {
		case receivedState := <-receiveStateCh:
			ID := receivedState.ID
			if ID != localID {

				stable, ok := StateTables.ReadWholeMap()[ID]
				if (stable[0][0] != 0 && receivedState.StateTable[0][0] != 0 && ok) || !ok {

					if receivedState.StateTable[0][2] == 0 {
						fmt.Println("DAyumn he ded yo")
					}
					StateTables.Write(ID, receivedState.StateTable)
					updatedLocalState, ok := checkIfExternalOrderCompleted(receivedState.StateTable)
					if ok {
						StateTables.Write(localID, updatedLocalState)
						stateTableTransmitCh <- Get()
					}
					updateHallLightsFromExternalOrders()
					RunOrderDistribution()
				} else {

				}

			} else if !recoveryCompleted {
				// Check if local was dead
				recoveryStatetable := receivedState.StateTable
				wasDead := recoveryStatetable[0][0]
				if wasDead == 0 {
					currentStatetable := Get()
					for row := 0; row < 4; row++ {
						currentStatetable[row+3][2] = recoveryStatetable[row+3][2]
					}
					StateTables.Write(localID, currentStatetable)
					RunOrderDistribution()
					recoveryCompleted = true
				}
			}
		default:
		}
	}
}

//checkIfExternalOrderCompleted checks if one of the external elevators has completed an order linked to hall button being pressed.
func checkIfExternalOrderCompleted(elevState [7][3]int) ([7][3]int, bool) {
	// fmt.Println("Checking external...")
	positionFloor := elevState[2][1]
	elevDirection := elevState[1][1]
	localStateTable := Get()
	updateLocal := false
	for col := 0; col < 2; col++ {
		if (localStateTable[3+positionFloor][col] == 1) && !(col == 0 && elevDirection == int(elevio.MD_Down)) && !(col == 1 && elevDirection == int(elevio.MD_Up)) {
			localStateTable[3+positionFloor][col] = 0
			updateLocal = true
			fmt.Println("ORDER ", positionFloor, " completed externally!")
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
				}
			} else {
				if activeLightsUpdate[[2]int{int(butn), floor}] {
					elevio.SetButtonLamp(butn, floor, false)
					activeLightsUpdate[[2]int{int(butn), floor}] = false
				}
			}

		}
	}
	activeLights.WriteWholeMap(activeLightsUpdate)
}

func toggleOffAllBtnLights() {
	for floor := 0; floor < 4; floor++ {
		for butn := elevio.BT_HallUp; butn <= elevio.BT_Cab; butn++ {
			elevio.SetButtonLamp(butn, floor, false)
		}
	}
}

// TransmitState repeatedly outputs the state table to be transmitted. The state table to be transmitted is updated via the stateTableTransmitCh-channel
func TransmitState(stateTableTransmitCh <-chan [7][3]int, transmitStateCh chan<- ElevatorState, receiveRecoveryStateCh <-chan ElevatorState) {
	ticker := time.NewTicker(config.StateTransmissionInterval)
	stateTable := ReadStateTable(localID)
	recoverSent := 0
	elevatorState := ElevatorState{ID: localID, StateTable: stateTable}
	for {
		select {
		case stateTable = <-stateTableTransmitCh:
			elevatorState.StateTable = stateTable
			elevatorState.ID = localID
		case <-ticker.C:
			// fmt.Println("Sending: ", elevatorState.ID)
			transmitStateCh <- elevatorState
			if recoverSent == 1 {
				elevatorState.StateTable = ReadStateTable(localID)
				elevatorState.ID = localID
				recoverSent--
			} else if recoverSent > 1 {
				recoverSent--
			}
		case recoveryState := <-receiveRecoveryStateCh:
			elevatorState.StateTable = recoveryState.StateTable
			elevatorState.ID = recoveryState.ID
			recoverSent = 5
		default:
		}
	}
}

//UpdateActiveElevators updates StateTables if an elevator goes offline or comes back online.
func UpdateActiveElevators(activeElevatorsCh <-chan map[string]bool, saveStateForRecoveryCh chan<- ElevatorState, recoveryIDCh chan<- string) {
	for {
		select {
		case activeElevators := <-activeElevatorsCh:
			for ID, isAlive := range activeElevators {
				stateTableUpdate := ReadStateTable(ID)

				if isAlive {
					if stateTableUpdate[0][0] == 0 {
						fmt.Println("Funeral ceremony /over")
						stateTableUpdate[0][0] = 1
						StateTables.Write(ID, stateTableUpdate)
						RunOrderDistribution()
						recoveryIDCh <- ID
					}
				} else {
					// fmt.Println("Funeral ceremony /begin ", stateTableUpdate[0][0])
					if stateTableUpdate[0][0] == 1 {
						// fmt.Println("Funeral ceremony /continue")
						stateTableUpdate[0][0] = 0
						StateTables.Write(ID, stateTableUpdate)
						RunOrderDistribution()
						saveStateForRecoveryCh <- ElevatorState{ID: ID, StateTable: ReadStateTable(ID)}
					}
				}
			}
		default:
		}
	}
}

func UpdateStateTableIndex(row, col int, ID string, val int, runDistribution bool) {
	stateTable := ReadStateTable(ID)
	stateTable[row][col] = val
	StateTables.Write(ID, stateTable)
	if runDistribution {
		RunOrderDistribution()
	}
}

func RunOrderDistribution() {
	allOrders, allDirections, elevStatuses := GetSyncedOrders()
	// fmt.Println("SyncedOrders!")
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

func GetSyncedOrders() ([4][3]int, map[string]int, map[string][2]int) { //omdøpe til noe som SyncOrdersDirectionsLocations (positions?)
	var allOrders [4][3]int
	var allDirections = make(map[string]int)
	var elevStatuses = make(map[string][2]int) // 1: Position, 2: Alive
	stateTables := StateTables.ReadWholeMap()
	for ID, statetable := range stateTables {
		var status [2]int
		isAlive := statetable[0][0] * statetable[0][2]
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
}

func GetCurrentFloor() int {
	stateTable := ReadStateTable(localID)

	floor := stateTable[2][1]
	return floor

	//floor := StateTables[localID][2][1]
	//return floor
}

func GetLocalID() string {
	return localID
}

func Get() [7][3]int {
	stateTable := ReadStateTable(localID)
	return stateTable

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

func StateTableRecovery(saveStateForRecoveryCh <-chan ElevatorState, recoveryIDCh <-chan string, transmitRecoveryStateCh chan<- ElevatorState) {
	recoveryStateTables := make(map[string][7][3]int)
	for {
		select {
		case recoveryElevatorState := <-saveStateForRecoveryCh:
			recoveryStateTables[recoveryElevatorState.ID] = recoveryElevatorState.StateTable
		case recoveryID := <-recoveryIDCh:
			transmitRecoveryState := ElevatorState{ID: recoveryID, StateTable: recoveryStateTables[recoveryID]}
			// for i := 0; i < 10; i++ {
			transmitRecoveryStateCh <- transmitRecoveryState
			// time.Sleep(100 * time.Millisecond)
			// }
		default:
		}
	}
}
