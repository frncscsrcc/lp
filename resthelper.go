package lp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// ErrorResponse is the generic error struct
type ErrorResponse struct {
	Error     bool
	ErrorCode int
	Message   string
}

// EventData is the exported data reppresentation
type EventData struct {
	TimeStamp time.Time
	Payload   interface{}
}

// EventsData is the final response, in case of event(s)
type EventsData struct {
	Error  bool
	Events []EventData
}

// SendError encode an error as JSON
func SendError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	log.Printf("ERROR: [%d] %s\n", code, message)
	json, err := toJSON(ErrorResponse{true, code, message})
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	fmt.Fprintf(w, json)
}

// SendOK returns a generic OK messages
func SendOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	object := struct {
		Error   bool
		Message string
	}{false, "OK"}
	json, err := toJSON(object)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	fmt.Fprintf(w, json)
}

// SendTimeout returns a timeout message
func SendTimeout(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(408)
	object := struct {
		Error   bool
		Message string
	}{true, "timeout"}
	json, err := toJSON(object)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	fmt.Fprintf(w, json)
}

// SendEvents returns the events
func SendEvents(w http.ResponseWriter, events []*Event) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	eventsResponse := EventsData{
		Error:  false,
		Events: make([]EventData, 0),
	}

	for _, e := range events {
		eventsResponse.Events = append(eventsResponse.Events, EventData{e.ts, e.payload})
	}

	json, err := toJSON(eventsResponse)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	fmt.Fprintf(w, json)
}

// SendResponse returns a generic JSON message
func SendResponse(w http.ResponseWriter, object interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json, err := toJSON(object)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	fmt.Fprintf(w, json)
}

func toJSON(object interface{}) (string, error) {
	json, err := json.Marshal(object)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

func fromJSON(byt []byte, dataStructure interface{}) error {
	if err := json.Unmarshal(byt, &dataStructure); err != nil {
		return err
	}
	return nil
}

// LogRequest logs each request
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Context().Value("sessionID")
		log.Println(r.Method, "-", sessionID, "-", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
