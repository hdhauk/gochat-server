package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type myCustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *credentials) IsValid() bool {
	return c.Username != "" && c.Password != ""
}

var handleLogin = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Pick off username and password from request body
	decoder := json.NewDecoder(r.Body)
	var cred credentials
	err := decoder.Decode(&cred)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error decoding JSON")
		return
	}
	if !cred.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Check if credentials is in database

	// Create Claims
	claims := myCustomClaims{
		cred.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "test",
		},
	}
	// Create Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString(myKey)
	if err != nil {
		fmt.Println(err)
	}
	// Add JWT to store in memory

	// Return JWT to user
	w.Write([]byte(ss))
})

var handleGetUsers = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(users) // TODO: Fetch data from DB instead
	if err != nil {
		fmt.Println("Error marhsalling JSON")
	}
	w.Write(json)
})

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

var handleGetMsgs = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Get IDs of messages in channel
	w.WriteHeader(http.StatusNotImplemented)
})

var handlePostMsg = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Post a new message to a channel
	w.WriteHeader(http.StatusNotImplemented)
})
