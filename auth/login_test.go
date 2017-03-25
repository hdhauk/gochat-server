package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
)

func TestStartingStoppingAddingUser(t *testing.T) {
	var testUsers = make(map[Username]PassSHA512)
	var testConfig = TokenServiceConfig{
		Secret:       []byte("test-secret"),
		PreAuthUsers: testUsers,
	}

	// Password = test-password
	expUsers := make(map[Username]PassSHA512)
	expUsers["test-user"] = "c638833f69bbfb3c267afa0a74434812436b8f08a81fd263c6be6871de4f1265"
	s := NewTokenService(testConfig)

	// Add when server is not running.
	err := s.addUser("test-user", "test-password")
	equals(t, ErrNotRunning, err)

	// Add testuser
	s.Start()
	err = s.addUser("test-user", "test-password")
	ok(t, err)
	equals(t, expUsers, s.users)

	// Add existing user with different password.
	err = s.addUser("test-user", "totally-different-password")
	equals(t, ErrUsernameTaken, err)

	// Add existing user with identical password.
	err = s.addUser("test-user", "test-password")
	equals(t, ErrUsernameTaken, err)

	// Stop then add existin users
	s.Stop()
	err = s.addUser("new-user", "password")
	equals(t, ErrNotRunning, err)
	err = s.addUser("test-user", "test-password")
	equals(t, ErrNotRunning, err)

	// Stop again to ensure nothing panics
	s.Stop()

	// Start again
	s.Start()
	equals(t, expUsers, s.users)
}

func TestDefaultService(t *testing.T) {
	s := NewTokenService(TokenServiceConfig{})
	expUsers := make(map[Username]PassSHA512)
	equals(t, expUsers, s.users)

}

func TestLoginHandlerInvalidLogin(t *testing.T) {
	var testUsers = make(map[Username]PassSHA512)
	var testConfig = TokenServiceConfig{
		Secret:       []byte("test-secret"),
		PreAuthUsers: testUsers,
	}

	// Set up service
	s := NewTokenService(testConfig)
	s.Start()

	// Test with an invalid user-name/password combination.
	payload := LoginReq{Username: "fake-user", Password: "fake-password"}
	reqBody := new(bytes.Buffer)
	json.NewEncoder(reqBody).Encode(payload)
	req, _ := http.NewRequest("GET", "/login", reqBody)
	w := httptest.NewRecorder()
	s.loginHandler(w, req)
	resp := w.Result()
	respBody, _ := ioutil.ReadAll(resp.Body)
	equals(t, http.StatusUnauthorized, resp.StatusCode)
	equals(t, []byte{}, respBody)

}

func TestLoginHandlerValidLogin(t *testing.T) {
	// Set up service
	testUsers := make(map[Username]PassSHA512)
	testConfig := TokenServiceConfig{
		Secret:       []byte("test-secret"),
		PreAuthUsers: testUsers,
	}
	s := NewTokenService(testConfig)
	s.Start()

	// Test with registered user/pass combination.
	err := s.addUser("test-user", "test-password")
	ok(t, err)
	payload := LoginReq{Username: "test-user", Password: "test-password"}
	reqBody := new(bytes.Buffer)
	json.NewEncoder(reqBody).Encode(payload)
	req, _ := http.NewRequest("GET", "/login", reqBody)
	w := httptest.NewRecorder()
	s.loginHandler(w, req)
	resp := w.Result()
	respBody, _ := ioutil.ReadAll(resp.Body)
	equals(t, http.StatusOK, resp.StatusCode)

	// Parse and validate JWT
	tokenStr := string(respBody)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("test-secret"), nil
	},
	)
	ok(t, err)
	ok(t, token.Claims.Valid())
	usr, ok := token.Claims.(jwt.MapClaims)["usr"]
	if !ok {
		t.Errorf("Unable to get user from claims\n")
	}
	equals(t, "test-user", usr)
}

func TestBadLoginRequest(t *testing.T) {
	// Set up service
	testUsers := make(map[Username]PassSHA512)
	testConfig := TokenServiceConfig{
		Secret:       []byte("test-secret"),
		PreAuthUsers: testUsers,
	}
	s := NewTokenService(testConfig)
	s.Start()

	// Test with an arbitrary json-formatted object
	payload := struct {
		a string
		b int
	}{a: "randomShit", b: 444}
	reqBody := new(bytes.Buffer)
	json.NewEncoder(reqBody).Encode(payload)
	req, _ := http.NewRequest("POST", "/login", reqBody)
	w := httptest.NewRecorder()
	s.loginHandler(w, req)
	resp := w.Result()
	equals(t, http.StatusBadRequest, resp.StatusCode)

	// Test with empty request
	req, _ = http.NewRequest("POST", "/login", new(bytes.Buffer))
	w = httptest.NewRecorder()
	s.loginHandler(w, req)
	resp = w.Result()
	equals(t, http.StatusBadRequest, resp.StatusCode)

}

// Helper functions
// =============================================================================

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
