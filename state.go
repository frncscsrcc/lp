package lp

type state int

const (
	stateWaiting state = iota + 1
	stateAbort
	stateTimeout
	stateOk
	stateReady
)

func (s state) String() string {
	switch s {
	case 1:
		return "Waiting for events"
	case 2:
		return "Aborted connection due new connection"
	case 3:
		return "Aborted connection due timeout"
	case 4:
		return "Sent event(s)"
	case 5:
		return "Handler can be destroyed"
	}
	return "Unknown"
}
