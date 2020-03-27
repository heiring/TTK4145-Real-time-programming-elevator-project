package network

type ElevatorState struct {
	Id string
	//state table
}

//broadcasts elevator state. Sends packets to be sent to transmission channel
func BroadcastElevatorState(stateUpdateCh <-chan ElevatorState, elevatorStateTxCh chan<- ElevatorState) {
	elevatorStateTx := ElevatorState{Id: "null"}
	for {
		select {
		case stateUpdate := <-stateUpdateCh:
			elevatorStateTx := stateUpdate
			elevatorStateTxCh <- elevatorStateTx

		default:
			elevatorStateTxCh <- elevatorStateTx
		}
	}
}

//Listens for elevator state packets, checks if an elevator has gone offline
func ListenElevatorState() {

}
