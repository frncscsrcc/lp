package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/frncscsrcc/lp"
)

func main() {

	feed1, _ := lp.NewFeed("feed1")

	go addEvent(feed1)

	http.HandleFunc("/newfeed", lp.CreateFeed)
	http.HandleFunc("/newevent", lp.NotifyEvent)
	http.HandleFunc("/subscribe", lp.SubscribeHandler)
	http.HandleFunc("/listen", lp.ListenHandler)
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

// Simulate new events (1 every 10 seconds)
func addEvent(feed *lp.Feed) {
	for {
		time.Sleep(5 * time.Second)
		type genericObject struct {
			A string
			B string
		}
		fmt.Println("New Event")
		lp.NewEvent(feed, genericObject{"A", "B"})
	}
}
