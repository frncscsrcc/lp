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
