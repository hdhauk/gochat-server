package chat

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestAddRoom(t *testing.T) {
	s := NewService()
	myRoom, id := NewRoom(nil, "My Room", "Top Secret stuff...")
	// Add room while service not running.
	err := s.AddRoom(myRoom)
	equals(t, ErrServiceNotRunning, err)

	// Create expected map
	exp := make(map[RoomID]RoomInterface)
	exp[id] = myRoom

	// Start room twice.
	s.Start()
	s.Start()

	// Add room.
	err = s.AddRoom(myRoom)
	ok(t, err)
	equals(t, exp, s.rooms)

	// Stop room twice.
	s.Stop()
	s.Stop()
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
