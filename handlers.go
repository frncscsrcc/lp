package lp

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func sendInternalError(w http.ResponseWriter) {
	if recovered := recover(); recovered != nil {
		fmt.Println(recovered)
		SendError(w, 500, "internal server error")
	}
}

// CreateFeed creates a new feed in the system
func CreateFeed(w http.ResponseWriter, r *http.Request) {
	// Send an internal error in case of panic.
	defer sendInternalError(w)

	// Check feeds
	feeds := extractFeeds(r)
	if len(feeds) == 0 {
		SendError(w, 400, "missing valid feed(s)")
		return
	}
	if len(feeds) > 1 {
		SendError(w, 400, "to many feeds")
		return
	}

	_, err := NewFeed(feeds[0])
	if err != nil {
		SendError(w, 500, "can not create feed "+feeds[0])
		return
	}

	SendOK(w)
	return
}

// SubscribeHandler is the handler to be use to listen for subscriptions
func SubscribeHandler(w http.ResponseWriter, r *http.Request) {

	// Send an internal error in case of panic.
	defer sendInternalError(w)

	// Check feeds
	feeds := getFeeds(r)
	if len(feeds) == 0 {
		SendError(w, 400, "missing valid feed(s)")
		return
	}

	// Create a new connection
	subscription := NewSubscription()

	// Subscribe the feeds
	for _, feed := range feeds {
		subscription.Subscribe(feed)
	}

	// Prepare the feed uid list
	feedUUIDs := make([]string, 0)
	for _, f := range feeds {
		feedUUIDs = append(feedUUIDs, f.name)
	}

	resp := struct {
		Feeds        []string
		ConnectionID string
	}{
		feedUUIDs,
		string(subscription.id),
	}
	SendResponse(w, resp)

	return
}

// ListenHandler is the handler to be use to listen for events
func ListenHandler(w http.ResponseWriter, r *http.Request) {

	// Send an internal error in case of panic.
	defer sendInternalError(w)

	subscriptionID := extractSubscription(r)
	subscription, err := GetSubscription(subscriptionID)
	if err != nil {
		SendError(w, 403, "not valid subscriptionID")
		return
	}

	// Check if there are listening connections
	subscription.l.Lock()
	if subscription.listening {
		// Send an abort signal to previous listening connection
		subscription.signal <- stateAbort
	}
	// Set there is an active listening connection
	subscription.listening = true
	subscription.l.Unlock()

	// Check one or more event are ready in the queue. It uses a goroutine
	// so the business logic (wait for events) could be unified
	go subscription.CheckForEvents()

	// Timeout
	timeout := extractTimeout(r)
	if timeout == 0 {
		timeout = 30
	}

	// Wait for some signal...
	select {

	// A message is sent in the communication channel
	case signal := <-subscription.signal:

		// Events are sent in the communication channel
		if signal == stateReady {
			events := subscription.GetEvents()

			subscription.l.Lock()
			subscription.listening = false
			subscription.l.Unlock()

			SendEvents(w, events)
			return
		}

		// An abort signal is sent to the communication channel
		if signal == stateAbort {
			SendError(w, 500, "ABORTED")
			return
		}

	// Timeout is triggered
	case <-time.After(time.Duration(timeout) * time.Second):
		subscription.l.Lock()
		subscription.listening = false
		SendTimeout(w)
		subscription.l.Unlock()
		return
	}

	return
}

// NotifyEvent notify a new event
func NotifyEvent(w http.ResponseWriter, r *http.Request) {

	// Send an internal error in case of panic.
	defer sendInternalError(w)

	// Check feeds
	feeds := getFeeds(r)
	if len(feeds) == 0 {
		SendError(w, 400, "missing valid feed(s)")
		return
	}
	if len(feeds) > 1 {
		SendError(w, 400, "too many feeds")
		return
	}

	e := struct {
		x int
		y int
	}{1, 2}

	_, err := NewEvent(feeds[0], e)
	if err != nil {
		SendError(w, 500, "can not register event")
		return
	}

	SendOK(w)
	return
}

func getFeeds(r *http.Request) []*Feed {
	var feeds = make([]*Feed, 0)

	feedNames := extractFeeds(r)
	if len(feedNames) > 0 {
		for _, feedName := range feedNames {
			feed, err := GetFeedFromName(feedName)
			if err == nil {
				feeds = append(feeds, feed)
			}
		}
	}
	return feeds
}

func extractFeeds(r *http.Request) []string {
	var feedsNames = make([]string, 0)

	var ok bool
	// Search in URL
	extractedFeedNames, ok := r.URL.Query()["feed"]

	if ok == true && len(extractedFeedNames) > 0 {
		return extractedFeedNames
	}
	return feedsNames

	// Search in body
	// TODO
}

func extractSubscription(r *http.Request) (subscriptionID uuid) {
	var ok bool

	// Search in URL
	subscriptionIDs, ok := r.URL.Query()["subscriptionID"]
	if ok == true && len(subscriptionIDs) > 0 {
		return uuid(subscriptionIDs[0])
	}
	return subscriptionID

	// Search in body
	// TODO
}

func extractTimeout(r *http.Request) int {
	var ok bool

	// Search in URL
	timeoutStrings, ok := r.URL.Query()["timeout"]
	if ok == true && len(timeoutStrings) == 1 {
		if timeoutInt, err := strconv.Atoi(timeoutStrings[0]); err == nil {
			return timeoutInt
		}
	}
	return 0

	// Search in body
	// TODO
}
