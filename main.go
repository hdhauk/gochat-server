package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", Index)

	// Get list of all chats
	r.HandleFunc("/chats/", handleChats).Methods("GET")

	// Get details for specific chat
	r.HandleFunc("/chats/{id}/", handleGetChat).Methods("GET")

	// Get IDs of messages in channel
	r.HandleFunc("/chats/{chatID}/msgs/", handleGetMsgs).Method("GET")
	// Post a new message to a channel
	r.HandleFunc("/chats/{chatID}/msgs/" handlePostMsg).Method("POST")



	log.Fatal(http.ListenAndServe(":8080", r))
}

// Index TODO: Fix text
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func handleChats(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This will return a list over all chats")
}
