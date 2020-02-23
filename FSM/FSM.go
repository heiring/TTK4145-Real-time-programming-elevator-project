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

	case IDLE:

	case MOVE:

	case WAIT:

	case EM_STOP:
	}
}
