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

func (c *credentials) isValid() bool {
	return c.Username != "" && c.Password != ""
}

func extractCredentials(r *http.Request) (credentials, error) {
	decoder := json.NewDecoder(r.Body)
	var ret credentials
	err := decoder.Decode(&ret)
	if err != nil && !ret.isValid() {
		return ret, err
	}
	return ret, nil
}

func checkCredentials(c credentials) (bool, error) {
	// FIXME
	return true, nil
}

var handleLogin = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Get credentials
	user, err := extractCredentials(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if credentials is in database
	login, err := checkCredentials(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !login {
		w.WriteHeader(http.StatusUnauthorized)
	}

	// Create Claims
	claims := myCustomClaims{
		user.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	// Create Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString(myKey)
	if err != nil {
		fmt.Println(err)
	}

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
