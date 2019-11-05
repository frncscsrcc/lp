package lp

type communicationChannel struct {
	events []*Event
	state  state
}
