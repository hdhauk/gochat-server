package chat_test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/hdhauk/gochat-server/chat"
)

func TestAddUser(t *testing.T) {
	users := map[chat.UserID]chat.AuthLevel{"a": 1, "b": 1, "c": 1}
	r, _ := chat.NewRoom(
		users,
		"TestRoom1",
		"TestRoom1-Topic",
	)
	equals(t, "TestRoom1", r.Name)
	equals(t, "TestRoom1-Topic", r.Topic)

	r.Start()

	err := r.AddUser("d")
	ok(t, err)

}

func TestAddUserWitoutStarting(t *testing.T) {
	r, _ := chat.NewRoom(
		map[chat.UserID]chat.AuthLevel{"a": 1, "b": 1, "c": 1},
		"TestRoom1",
		"TestRoom1-Topic",
	)
	equals(t, "TestRoom1", r.Name)
	equals(t, "TestRoom1-Topic", r.Topic)

	r.Start()
	err := r.AddUser("d")
	equals(t, nil, err)

	r.Stop()
	time.Sleep(2 * time.Nanosecond)
	err = r.AddUser("d")
	equals(t, chat.ErrRoomNotRunning, err)

}

type mockSession struct {
}

func (m *mockSession) ReadJSON(i interface{}) error  { return nil }
func (m *mockSession) WriteJSON(i interface{}) error { return nil }
func (m *mockSession) Close() error                  { return nil }

func TestAddSession(t *testing.T) {
	users := map[chat.UserID]chat.AuthLevel{"a": 1, "b": 1, "c": 1}
	r, _ := chat.NewRoom(
		users,
		"TestRoom1",
		"TestRoom1-Topic",
	)
	r.Start()
	s := &mockSession{}

	// Add authenticated user.
	id := chat.UserID("a")
	err := r.AddSession(id, s)
	ok(t, err)

	// Add unauthenticated user.
	id = chat.UserID("q")
	err = r.AddSession(id, s)
	equals(t, chat.ErrUserNotAuth, err)

	// Add when service not running
	r.Stop()
	time.Sleep(2 * time.Microsecond)
	err = r.AddSession(id, s)
	equals(t, chat.ErrRoomNotRunning, err)
}

func TestListAuthorizedUsers(t *testing.T) {
	expList := map[chat.UserID]chat.AuthLevel{"a": 1, "b": 1, "c": 1}
	r, _ := chat.NewRoom(
		expList,
		"TestRoom1",
		"TestRoom1-Topic",
	)
	_, err := r.ListUsers()
	equals(t, chat.ErrRoomNotRunning, err)

	r.Start()
	gotList, err := r.ListUsers()
	ok(t, err)
	equals(t, expList, gotList)

	// Test Removing users
	err = r.KickUser("a")
	ok(t, err)
	gotList, err = r.ListUsers()
	ok(t, err)
	equals(t, map[chat.UserID]chat.AuthLevel{"b": 1, "c": 1}, gotList)

	// Test removing non-existing user.
	err = r.KickUser("unknownUser")
	equals(t, chat.ErrUserNotFound, err)

	// Test kicking when not running.
	r.Stop()
	time.Sleep(2 * time.Microsecond)
	err = r.KickUser("b")
	equals(t, chat.ErrRoomNotRunning, err)
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
