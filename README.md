[![Go Report Card](https://goreportcard.com/badge/github.com/frncscsrcc/lp)](https://goreportcard.com/report/github.com/frncscsrcc/lp)

lp: Go long-poll library
---

Note: this is a WiP and it is mainly a battlefield to improve my Golang skills.
*Use this library at your own risk*. Suggestions, bug reports and feature requests
are more than welcome.

Create a simple server
---

Working examples in the examples folder.

```
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/frncscsrcc/lp"
)

func main() {

	// Registers a event parse function. This function is application domain
	// specific and must be provided. This function receive the body of each
	// new event request (as string) and returns a interface{}. The returned
	// value will be used as payload in the internal event structure.
	// This function is required only if you want to receive new events from
	// the REST end point
	lp.RegisterEventParser(eventParser)

	// Create a new feed
	feed1, _ := lp.NewFeed("feed1")

	// Send events generated from the server
	go simulateServerEvents(feed1)

	http.HandleFunc("/newfeed", lp.CreateFeed)
	http.HandleFunc("/newevent", lp.NotifyEvent)
	http.HandleFunc("/subscribe", lp.SubscribeHandler)
	http.HandleFunc("/listen", lp.ListenHandler)

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var eventParser = func(bodyString string) (interface{}, error) {
	// 1. decode the bodyString (eg: JSON) ...
	// ...

	// 2. Generate the payload for the new event based on business logic
	type genericObject struct {
		A string
		b string
	}

	// 3. Prepare the payload
	return genericObject{
		A: "This is A and is exported",
		b: "This is b and is not exported (you should avoid it!)",
	}, nil
}


// This function simulates new events incoming
func simulateServerEvents(feed *lp.Feed) {
	for {
		time.Sleep(10 * time.Second)
		type payload struct {
			Data string
		}
		lp.NewEvent(feed, payload{"This is some data generated by the server. It will not be parsed by the eventParse function"})
	}
}

```

Create a simple client using the Golang SDK
---

```
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/frncscsrcc/lp"
)

type myClient struct{}

// Connect function
func (c myClient) EventsHandler(events []lp.EventData, err error) bool {
	if err != nil {
		// In this example print the error and, disconnect
		log.Printf("*** [ERROR]: %s\n", err)
		return false
	}
	fmt.Printf("\n\n*** Received data %+v\n\n", events)
	// Wait for new data
	return true
}

func main() {
	SDK := lp.SDK{
		Protocol: "http",
		Host:     "127.0.0.1",
		Port:     8080,
		Feeds:    []string{"feed1"},
		Timeout:  5,
		Debug:    true,
	}

	var client myClient

	// This function will subscribe to the feed "feed1" and cycling until an
	// error happens. In this case EventsHandler() will intercept the error
	// and force SDK.Connect() to return.
	err := SDK.Connect(client)
	if err != nil {
		os.Exit(1)
	}
}
```
