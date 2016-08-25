package main

import (
	"fmt"
	"net/http"
)

func handleChats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	fmt.Fprintln(w, "This should return a list of all the chats available")
}

func handleGetChat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	// Get details for specific chat
}

func handleGetMsgs(w http.ResponseWriter, r *http.Request) {
	// Get IDs of messages in channel
	w.WriteHeader(http.StatusNotImplemented)
}

func handlePostMsg(w http.ResponseWriter, r *http.Request) {
	// Post a new message to a channel
	w.WriteHeader(http.StatusNotImplemented)
}
