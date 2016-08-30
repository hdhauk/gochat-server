package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

const myKey = "Password123"

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

func handleLogin(w http.ResponseWriter, r *http.Request) {
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
			ExpiresAt: 15000,
			Issuer:    "test",
		},
	}
	// Create Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(myKey))
	if err != nil {
		fmt.Println(err)
	}
	// Add JWT to store in memory

	// Return JWT to user
	w.Write([]byte(ss))
}
func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(users) // TODO: Fetch data from DB instead
	if err != nil {
		fmt.Println("Error marhsalling JSON")
	}
	w.Write(json)
}

func handleGetChats(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(chats) // TODO: Fetch data from DB instead
	if err != nil {
		fmt.Println("Error marhsalling JSON")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
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
