package lp

import (
	"errors"
	"log"
	"sync"
)

var feedLock sync.Mutex
var feeds map[uuid]*Feed
var feedNameToUUID map[string]uuid

func init() {
	feeds = make(map[uuid]*Feed)
	feedNameToUUID = make(map[string]uuid)
}

// Feed is the object that reppresent a feed
type Feed struct {
	l             sync.Mutex
	name          string
	id            uuid
	subscriptions map[uuid]*Subscription
}

// NewFeed tries to create a new feed and returns it
func NewFeed(feedName string) (*Feed, error) {
	// Locl to be sure feedName is uniq
	feedLock.Lock()
	defer feedLock.Unlock()

	f := new(Feed)
	if _, exists := feedNameToUUID[feedName]; exists {
		return f, errors.New("feed " + feedName + " exists")
	}
	id := newUUID()
	feedNameToUUID[feedName] = id
	f.name = feedName
	f.id = id
	f.subscriptions = make(map[uuid]*Subscription)
	feeds[f.id] = f
	return f, nil
}

// GetFeed returns a feed object ptr, if exists
func GetFeed(id uuid) (*Feed, error) {
	var f *Feed
	var exists bool
	if f, exists = feeds[id]; exists == false {
		return f, errors.New("feed " + string(id) + " does not exists")
	}
	return f, nil
}

// GetFeedFromName returns a feed ptr from a feed name, if exists
func GetFeedFromName(feedName string) (*Feed, error) {
	var id uuid
	var exists bool
	if id, exists = feedNameToUUID[feedName]; exists == false {
		return new(Feed), errors.New("feed " + feedName + " does not exists")
	}

	f := feeds[id]
	return f, nil
}

// addSubscription add a connection to a specific feed
func (f *Feed) addSubscription(c *Subscription) error {
	c.l.Lock()
	defer c.l.Unlock()

	if _, exists := f.subscriptions[c.id]; exists {
		return errors.New("connection " + string(c.id) + " already subscribed feed " + f.name)
	}
	f.subscriptions[c.id] = c
	return nil
}

// removeSubscription remove a connection to a specific feed
func (f *Feed) removeSubscription(c *Subscription) error {
	c.l.Lock()
	defer c.l.Unlock()

	delete(f.subscriptions, c.id)
	return nil
}

func (f *Feed) String() string {
	return f.name + "( F:" + string(f.id) + ")"
}

// Log logs connection in STDOUT
func (f *Feed) Log() {
	log.Printf("%s (%v)\n", f.name, f.id)
	for _, c := range f.subscriptions {
		log.Printf("|-- %s\n", c)
	}
	log.Print("\n")
}
