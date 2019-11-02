package lp

import (
	"fmt"
	"net/http"
)

func sendInternalError(w http.ResponseWriter) {
	if recovered := recover(); recovered != nil {
		fmt.Println(recovered)
		SendError(w, 500, "internal server error")
	}
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
	i := 0
	fmt.Printf("***** %v\n\n", (1/i - 1))
	// Create a new connection
	connection := NewSubscription()

	// Prepare a new listener handler and save in the connection object
	listenHandler := func(w http.ResponseWriter, r *http.Request) {
		SendResponse(w, "YES!")
	}

	// Returns an error if it can not return a new listener handler
	err := connection.SetHandler(listenHandler)
	if err != nil {
		SendError(w, 500, "can not return a new connection")
		return
	}

	feedUUIDs := make([]string, 0)
	for _, f := range feeds {
		feedUUIDs = append(feedUUIDs, string(f.id))
	}

	resp := struct {
		Feeds        []string
		ConnectionID string
	}{
		feedUUIDs,
		string(connection.id),
	}
	SendResponse(w, resp)

	return
}

// ListenHandler is the handler to be use to listen for events
func ListenHandler(w http.ResponseWriter, r *http.Request) {

	// Send an internal error in case of panic.
	defer sendInternalError(w)

	subscriptionID := getSubscription(r)
	connection, err := GetSubscription(subscriptionID)
	if err != nil {
		SendError(w, 403, "not valid subscriptionID")
		return
	}
	connection.handler(w, r)
	return
}

func getFeeds(r *http.Request) []*Feed {
	var feeds = make([]*Feed, 0)

	var ok bool
	// Search in URL
	feedNames, ok := r.URL.Query()["feed"]

	if ok == true && len(feedNames) > 0 {
		for _, feedName := range feedNames {
			feed, err := GetFeedFromName(feedName)
			if err == nil {
				feeds = append(feeds, feed)
			}
		}
	}
	return feeds

	// Search in body
	// TODO
}

func getSubscription(r *http.Request) (subscriptionID uuid) {
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
