package main

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

func handleAuth(w http.ResponseWriter, r *http.Request) {
	// Check if credentials is in database

	// Create JWT and add to DB

	// Return JWT to user
	mySigningKey := []byte("AllYourBase")

	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: 15000,
		Issuer:    "test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Printf("%v %v\n", ss, err)
}

func handleChats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Println(r)
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
