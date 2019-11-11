package lp

import (
	"encoding/json"
	"errors"
	"log"
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

// EventParserFunction is the signature of the function must be provided in
// order to parse an incoming JSON to an internal Event.paload
type EventParserFunction func(JSON string) (interface{}, error)

var parserFunction EventParserFunction

func init() {
	events = new(eventList)
	events.list = make([]*Event, 0)
	parserFunction = func(bodyString string) (interface{}, error) {
		return nil, errors.New("Parser function not registered")
	}
}

// RegisterEventParser sets the event parse logic function
func RegisterEventParser(f EventParserFunction) {
	parserFunction = f
}

// NewEvent generate a new event and prepare the internal data model
func NewEvent(feed *Feed, payload interface{}) (*Event, error) {
	ev := new(Event)

	// Check if they are registered listeners
	if len(feed.subscriptions) == 0 {
		return ev, errors.New("no subscribers, this event will be lost")
	}

	// Prepare the event
	ev.id = newUUID()
	ev.feed = feed
	ev.ts = time.Now().UTC()
	ev.payload = payload

	// Append the event to the global list
	events.Append(ev)

	log.Printf("New event received, broadcasting to %d clients.\n", len(feed.subscriptions))

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
