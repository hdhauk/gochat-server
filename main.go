package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var myKey = []byte("Password123")

func initDatabase() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "datastore.db")
	defer db.Close()
	if err != nil {
		log.Fatal("Unable to access database")
		return db, err
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
	return db, nil
}

func main() {
	// Print welcome message and runtime
	go printRunTime()

	// Set up database
	initDatabase() //TODO: keep the database to make later calls to it..

	// Set up muxing
	r := mux.NewRouter().StrictSlash(false)

	// Set up JWT Middleware
	jwtMiddleware := jwtmiddleware.New(
		jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return myKey, nil
			},
			SigningMethod: jwt.SigningMethodHS512,
			Debug:         false,
		})

	//	Authentication Controller
	//----------------------------------------------------------------------------
	// Check credentials and return JWT if they check out
	r.Handle("/users/", handleLogin).Methods("POST")
	r.Handle("/users/", jwtMiddleware.Handler(handleGetUsers)).Methods("GET")

	//	Chats Controller
	//----------------------------------------------------------------------------
	// Get list of all chats
	r.Handle("/chats/", jwtMiddleware.Handler(handleGetChats)).Methods("GET")

	// Get details for specific chat
	r.Handle("/chats/{id}/", jwtMiddleware.Handler(handleGetChat)).Methods("GET")

	//	Message Controller
	//----------------------------------------------------------------------------
	// Get IDs of messages in channel
	r.Handle("/chats/{chatID}/msgs/", jwtMiddleware.Handler(handleGetMsgs)).Methods("GET")
	// Post a new message to a channel
	r.Handle("/chats/{chatID}/msgs/", jwtMiddleware.Handler(handlePostMsg)).Methods("POST")

	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}
