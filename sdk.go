package lp

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// SDK are the connection parameters
type SDK struct {
	Protocol       string
	Host           string
	Port           int
	Feeds          []string
	Timeout        int
	Debug          bool
	subscriptionID string
}

// LongPollClient is the interface that should be passed to an SDK client
type LongPollClient interface {
	EventsHandler([]EventData, error) bool
}

// Connect main method to interact with SDK
func (sdk *SDK) Connect(lpc LongPollClient) error {
	var events []EventData

	protocol := sdk.Protocol
	if protocol == "" {
		protocol = "http"
	}

	host := sdk.Host
	if host == "" {
		host = "localhost"
	}

	port := sdk.Port
	if port == 0 {
		port = 8080
	}

	timeout := sdk.Timeout
	if timeout == 0 {
		timeout = 30
	}

	// PRepare subscription feed list
	feeds := ""
	for _, feedName := range sdk.Feeds {
		feeds += "feed=" + feedName + "&"
	}

	serverURL := getServerURL(protocol, host, port)

	// 1. Subscribe to one or more feeds
	subscriptionRequestURL := serverURL + "/subscribe?" + feeds
	sdk.log("-> %s\n", subscriptionRequestURL)

	subscriptionID, err := getSubscriptionID(subscriptionRequestURL)
	sdk.subscriptionID = subscriptionID
	sdk.log("<- SubscriptionID=%s\n", subscriptionID)

	if err != nil {
		if lpc.EventsHandler(events, err) == false {
			return err
		}
	}

	// 2. Listen and send the events to the callback. Stop when the callback
	//    returns false
	listenRequestURL := serverURL + "/listen?subscriptionID=" + sdk.subscriptionID + "&timeout=" + strconv.Itoa(timeout)

	for true {
		sdk.log("-> %s\n", listenRequestURL)

		events, timeout, err := getEvents(listenRequestURL)
		if timeout {
			sdk.log("<- Timeout, reconnect...\n")
			continue
		}

		sdk.log("<- %+v\n", events)

		if lpc.EventsHandler(events, err) == false {
			sdk.log("STOP\n")
			break
		}
	}

	return nil
}

func (sdk *SDK) log(format string, i ...interface{}) {
	if sdk.Debug {
		if len(i) > 0 {
			log.Printf(format, i)
		} else {
			log.Printf(format)
		}
	}
}

func getSubscriptionID(subscriptionRequestURL string) (string, error) {
	httpResponse, err := http.Get(subscriptionRequestURL)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return "", err
	}

	type decodedResponse struct{ SubscriptionID string }
	var sr decodedResponse
	err = fromJSON(body, &sr)
	if err != nil {
		return "", err
	}
	if sr.SubscriptionID == "" {
		return "", errors.New("server did not return SubscriptionID")
	}
	return sr.SubscriptionID, nil
}

func getEvents(listenRequestURL string) (events []EventData, timeout bool, err error) {
	httpResponse, err := http.Get(listenRequestURL)
	if err != nil {
		return events, false, err
	}

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return events, false, err
	}

	type decodedResponse struct {
		Error   bool
		Message string
		Events  []EventData
	}
	var resp decodedResponse
	err = fromJSON(body, &resp)
	if err != nil {
		return events, false, err
	}

	// Return timeout
	if resp.Error == true && resp.Message == "timeout" {
		return events, true, nil
	}

	// Extract events
	return resp.Events, false, nil
}

func getServerURL(protocol string, host string, port int) string {
	return protocol +
		"://" +
		host +
		":" +
		strconv.Itoa(port)
}
