package lp

import (
	"errors"
	"fmt"
	"sync"
)

var subscriptions map[uuid]*Subscription

func init() {
	subscriptions = make(map[uuid]*Subscription)
}

// Subscription is the object that reppresent a connection
type Subscription struct {
	l         sync.Mutex
	id        uuid
	feeds     map[uuid]*Feed
	listening bool
	signal    chan state
	events    []*Event
}

// NewSubscription tries to create a new connection object and returns it
func NewSubscription() *Subscription {
	id := newUUID()
	s := new(Subscription)
	s.id = id
	s.feeds = make(map[uuid]*Feed)
	s.signal = make(chan state)
	s.events = make([]*Event, 0)
	s.listening = false
	subscriptions[s.id] = s
	return s
}

// GetSubscription returns a connection object ptr, if exists
func GetSubscription(id uuid) (*Subscription, error) {
	var c *Subscription
	var exists bool

	if c, exists = subscriptions[id]; exists == false {
		return c, errors.New("connection " + string(id) + " does not exists")
	}

	return c, nil
}

// Subscribe allows a connection to subscribe to a particular feed
func (s *Subscription) Subscribe(feed *Feed) error {
	err := feed.addSubscription(s)
	if err != nil {
		return err
	}
	s.feeds[feed.id] = feed
	return nil
}

// Unsubscribe allows a connection to unsubscribe from a particular feed
func (s *Subscription) Unsubscribe(feed *Feed) error {
	err := feed.removeSubscription(s)
	if err != nil {
		return err
	}
	delete(s.feeds, feed.id)
	return nil
}

// NotifyEvent notify the event to this subscription list
func (s *Subscription) NotifyEvent(e *Event) {
	s.l.Lock()
	defer s.l.Unlock()

	s.events = append(s.events, e)

	// If a listener is connected, notify an event is ready
	if s.listening {
		s.signal <- stateReady
	}
}

// CheckForEvents checks if there are events in the subscriber queue,
// in case send them in the communicationChannel
func (s *Subscription) CheckForEvents() {
	s.l.Lock()
	defer s.l.Unlock()

	if len(s.events) == 0 {
		return
	}

	s.signal <- stateReady
}

// GetEvents returns the events for this subscription
func (s *Subscription) GetEvents() []*Event {
	s.l.Lock()
	defer s.l.Unlock()

	events := make([]*Event, 0)
	if len(s.events) == 0 {
		return events
	}

	for _, e := range s.events {
		events = append(events, e)
	}

	// Clean the list for this subscription
	s.events = make([]*Event, 0)

	return events
}

func (s *Subscription) String() string {
	return "C:" + string(s.id)
}

// Log logs connection in STDOUT
func (s *Subscription) Log() {
	fmt.Printf("%s\n", string(s.id))
	for _, f := range s.feeds {
		fmt.Printf("|-- %s\n", f)
	}
	fmt.Print("\n")
}
