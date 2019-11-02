package lp

import (
	"encoding/json"
	"time"
)

// Event is the exported datamodel for an emitted event
type Event struct {
	id      uuid
	feed    *Feed
	ts      time.Time
	payload interface{}
}

type eventList []string

// NewEvent generate a new event and prepare the internal data model
func NewEvent(feedID uuid, payload interface{}) (*Event, error) {
	ev := new(Event)
	feed, err := GetFeed(feedID)
	if err != nil {
		return ev, err
	}
	ev.id = newUUID()
	ev.feed = feed
	ev.ts = time.Now().UTC()
	ev.payload = payload
	return ev, nil
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
