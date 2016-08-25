package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)

	// Get list of all chats
	r.HandleFunc("/chats/", handleChats).Methods("GET")

	// Get details for specific chat
	r.HandleFunc("/chats/{id}/", handleGetChat).Methods("GET")

	// Get IDs of messages in channel
	r.HandleFunc("/chats/{chatID}/msgs/", handleGetMsgs).Methods("GET")
	// Post a new message to a channel
	r.HandleFunc("/chats/{chatID}/msgs/", handlePostMsg).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}
