package lp

import (
	"encoding/json"
	"sync"
	"time"
)

// Event is the exported datamodel for an emitted event
type Event struct {
	id      uuid
	feed    *Feed
	ts      time.Time
	payload interface{}
}

type eventList struct {
	l    sync.Mutex
	list []*Event
}

var events *eventList

func init() {
	events = new(eventList)
	events.list = make([]*Event, 0)
}

// NewEvent generate a new event and prepare the internal data model
func NewEvent(feedID uuid, payload interface{}) (*Event, error) {
	ev := new(Event)
	feed, err := GetFeed(feedID)
	if err != nil {
		return ev, err
	}

	// Prepare the event
	ev.id = newUUID()
	ev.feed = feed
	ev.ts = time.Now().UTC()
	ev.payload = payload

	// Append the event to the global list
	events.Append(ev)

	// Notify listener
	for _, s := range feed.subscriptions {
		// Event subscription can happen in parallel, because it is thread safe
		go s.NotifyEvent(ev)
	}

	return ev, nil
}

// Append adds an event to the global event list
// It is thread safe
func (el *eventList) Append(ev *Event) {
	el.l.Lock()
	defer el.l.Unlock()

	el.list = append(el.list, ev)
}

// ToJSON returns a json encoded reppresentation of an Event object
func (ev Event) ToJSON() (string, error) {
	exported := struct {
		Feed    string
		Stamp   time.Time
		Payload interface{}
	}{
		ev.feed.name,
		ev.ts,
		ev.payload,
	}
	json, err := json.Marshal(exported)
	if err != nil {
		return "", err
	}
	return string(json), nil
}
