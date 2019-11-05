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
	cc        chan communicationChannel
	events    []*Event
}

// NewSubscription tries to create a new connection object and returns it
func NewSubscription() *Subscription {
	id := newUUID()
	s := new(Subscription)
	s.id = id
	s.feeds = make(map[uuid]*Feed)
	s.cc = make(chan communicationChannel)
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
