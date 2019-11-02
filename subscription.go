package lp

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
)

var subscriptions map[uuid]*Subscription

func init() {
	subscriptions = make(map[uuid]*Subscription)
}

// Subscription is the object that reppresent a connection
type Subscription struct {
	l            sync.Mutex
	id           uuid
	feeds        map[uuid]*Feed
	handlerIsSet bool
	handler      func(http.ResponseWriter, *http.Request)
	channel      chan State
}

// NewSubscription tries to create a new connection object and returns it
func NewSubscription() *Subscription {
	id := newUUID()
	c := new(Subscription)
	c.id = id
	c.feeds = make(map[uuid]*Feed)
	c.channel = make(chan State)
	c.handlerIsSet = false
	subscriptions[c.id] = c
	return c
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

// SetHandler register the active handler for this connection
func (c *Subscription) SetHandler(h func(http.ResponseWriter, *http.Request)) error {
	if c.handlerIsSet {
		// If an handler is already present, sent an ABORT message
		c.channel <- ABORTED
		// Wait the previous handler sent the ABORT message
		done := <-c.channel
		if done != READY {
			return errors.New("can not send abort response to connection " + string(c.id))
		}
	}
	c.handler = h
	c.handlerIsSet = true
	return nil
}

// Subscribe allows a connection to subscribe to a particular feed
func (c *Subscription) Subscribe(feed *Feed) error {
	err := feed.addSubscription(c)
	if err != nil {
		return err
	}
	c.feeds[feed.id] = feed
	return nil
}

// Unsubscribe allows a connection to unsubscribe from a particular feed
func (c *Subscription) Unsubscribe(feed *Feed) error {
	err := feed.removeSubscription(c)
	if err != nil {
		return err
	}
	delete(c.feeds, feed.id)
	return nil
}

func (c *Subscription) String() string {
	return "C:" + string(c.id)
}

// Log logs connection in STDOUT
func (c *Subscription) Log() {
	fmt.Printf("%s\n", string(c.id))
	for _, f := range c.feeds {
		fmt.Printf("|-- %s\n", f)
	}
	fmt.Print("\n")
}
