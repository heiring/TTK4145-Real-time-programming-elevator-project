package FSM

type elev_state int

const (
	INIT elev_state = iota
	IDLE
	MOVE
	WAIT
	EM_STOP
)

func StateTransistions(state elev_state) {
	switch state {
	case INIT:

	case IDLE:

	case MOVE:

	case WAIT:

	case EM_STOP:
	}
}

func FSM(state elev_state) {
	switch state {
	case INIT:
	//elevator reaches first floor, transistion to IDLE
	case IDLE:
	//transitions to MOVE when an order is detected
	//transisiton to EM_STOP
	case MOVE:
	//transisiton to wait when elevator reaches target floor
	//transisiton to EM_stop
	case WAIT:
	//transisiton to MOVE if there are pending orders
	//transition to IDLE if not
	case EM_STOP:
	}
}
