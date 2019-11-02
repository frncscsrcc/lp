package lp

type State int

const (
	WAITING State = iota + 1
	ABORTED
	TIMEOUT
	OK
	READY
)

func (s State) String() string {
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
