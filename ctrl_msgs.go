package main

import "net/http"

var handleGetMsgs = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Get IDs of messages in channel
	w.WriteHeader(http.StatusNotImplemented)
})

var handlePostMsg = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Post a new message to a channel
	w.WriteHeader(http.StatusNotImplemented)
})
