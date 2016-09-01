package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var handleGetChats = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(chats) // TODO: Fetch data from DB instead
	if err != nil {
		fmt.Println("Error marhsalling JSON")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
})

var handleGetChat = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	// Get details for specific chat
})
