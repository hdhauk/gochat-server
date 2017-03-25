package auth

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type (
	// Username is the loginname used for authentication
	Username string
	// PassSHA512 is a SHA512 hashed password.
	PassSHA512 string
)

// LoginReq is a loginrequest to the server.
type LoginReq struct {
	Username string `json:"usr"`
	Password string `json:"pwd"`
}

type userCheck struct {
	user Username
	ret  chan PassSHA512
}
type userUpdate struct {
	user Username
	hash PassSHA512
	err  chan error
}

// TokenService is an endpoint that provide users that authenticate with JWTs.
type TokenService struct {
	secret       []byte
	users        map[Username]PassSHA512
	endpointPort string
	endpointRoot string

	// Internal communication
	checkUserCh chan userCheck
	addUserCh   chan userUpdate
	stopCh      chan struct{}
}

// TokenServiceConfig provides configuration for a TokenService.
type TokenServiceConfig struct {
	Secret       []byte
	PreAuthUsers map[Username]PassSHA512
	EndpointPort string
	EndpointRoot string
}

// NewTokenService returns a new valid TokenService
func NewTokenService(c TokenServiceConfig) *TokenService {
	stopCh := make(chan struct{})
	close(stopCh)
	t := TokenService{
		checkUserCh: make(chan userCheck),
		addUserCh:   make(chan userUpdate),
		stopCh:      stopCh,
	}
	if c.PreAuthUsers == nil {
		c.PreAuthUsers = make(map[Username]PassSHA512)
	}
	if c.EndpointPort == "" {
		c.EndpointPort = ":2000"
	}
	if c.EndpointRoot == "" {
		c.EndpointRoot = "/login"
	}
	t.secret = c.Secret
	t.users = c.PreAuthUsers
	t.endpointPort = c.EndpointPort
	t.endpointRoot = c.EndpointRoot
	return &t
}

// Start serving JWTs to users that can authenticate.
func (t *TokenService) Start() {
	select {
	case <-t.stopCh:
		t.stopCh = make(chan struct{})
	default:
	}
	go t.authenticator()
	go func() {
		r := mux.NewRouter()
		r.HandleFunc(t.endpointRoot, t.loginHandler)
		r.Methods("POST")
		srv := &http.Server{
			Handler:      r,
			Addr:         t.endpointPort,
			WriteTimeout: 5 * time.Second,
			ReadTimeout:  5 * time.Second,
		}
		srv.ListenAndServe()
		//srv.ListenAndServeTLS("server.crt", "server.key")
		<-t.stopCh
		srv.Close()
	}()
}

// Stop serving JWTs.
func (t *TokenService) Stop() {
	select {
	case <-t.stopCh:
		return
	default:
		close(t.stopCh)
	}
}

func (t *TokenService) addUser(username, password string) error {
	if t.notRunning() {
		return ErrNotRunning
	}
	sha := sha256.Sum256([]byte(password))
	errCh := make(chan error)
	update := userUpdate{
		user: Username(username),
		hash: PassSHA512(fmt.Sprintf("%x", sha)),
		err:  errCh,
	}

	t.addUserCh <- update
	return <-errCh
}

func (t *TokenService) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Password == "" || req.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashCh := make(chan PassSHA512)
	t.checkUserCh <- userCheck{user: Username(req.Username), ret: hashCh}

	// Compare request with password store.
	serverHash := <-hashCh
	reqHash := fmt.Sprintf("%x", sha256.Sum256([]byte(req.Password)))
	if serverHash != PassSHA512(reqHash[:len(reqHash)]) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usr": string(req.Username),
		"iat": time.Now().Unix(),
	})
	ss, err := token.SignedString(t.secret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Return JWT
	w.Write([]byte(ss))

}

func (t *TokenService) authenticator() {
	for {
		select {
		case check := <-t.checkUserCh:
			if hash, ok := t.users[check.user]; ok {
				check.ret <- hash
			} else {
				check.ret <- ""
			}
		case add := <-t.addUserCh:
			if _, ok := t.users[add.user]; ok {
				add.err <- ErrUsernameTaken
			} else {
				t.users[add.user] = add.hash
				add.err <- nil
			}
		case <-t.stopCh:
			return
		}
	}
}

func (t *TokenService) notRunning() bool {
	select {
	case <-t.stopCh:
		return true
	default:
		return false
	}
}

var (
	// ErrNotRunning is returned whenever an action fail due to the service not
	// not running.
	ErrNotRunning = errors.New("service not running")
	// ErrUsernameTaken is returned when trying to add a new user with an existing
	// username.
	ErrUsernameTaken = errors.New("username already taken")
)
