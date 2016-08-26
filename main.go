package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	// Print welcome message and runtime
	go printRunTime()

	// Connect to database
	db, err := gorm.Open("sqlite3", "datastore.db")
	defer db.Close()
	if err != nil {
		log.Fatal("Unable to access database")
	}

	// Set up tables if none exists
	if !db.HasTable(&User{}) {
		fmt.Println("No Users-table found --> Creating new")
		db.CreateTable(&User{})
	}
	if !db.HasTable(&Chat{}) {
		fmt.Println("No Chats-table found --> Creating new")
		db.CreateTable(&Chat{})
	}
	if !db.HasTable(&Message{}) {
		fmt.Println("No Messages-table found --> Creating new")
		db.CreateTable(&Message{})
	}

	// Set up muxing
	r := mux.NewRouter().StrictSlash(false)

	//----------------------------------------------------------------------------
	//	Authentication Controller
	//----------------------------------------------------------------------------
	// Check credentials and return JWT if they check out
	r.HandleFunc("/users/", handleAuth).Methods("POST")

	//----------------------------------------------------------------------------
	//	Chats Controller
	//----------------------------------------------------------------------------
	// Get list of all chats
	r.HandleFunc("/chats/", handleChats).Methods("GET")

	// Get details for specific chat
	r.HandleFunc("/chats/{id}/", handleGetChat).Methods("GET")

	//----------------------------------------------------------------------------
	//	Message Controller
	//----------------------------------------------------------------------------
	// Get IDs of messages in channel
	r.HandleFunc("/chats/{chatID}/msgs/", handleGetMsgs).Methods("GET")
	// Post a new message to a channel
	r.HandleFunc("/chats/{chatID}/msgs/", handlePostMsg).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}
