package websocket_test

import (
	"crypto/tls"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	myWS "github.com/hdhauk/gochat-server/websocket"
	"github.com/stretchr/testify/assert"
)

// TestSecureConnection opens a simple websocket connection to the hub.
func TestSecureConnection(t *testing.T) {
	h := myWS.NewHub()
	s := make(chan myWS.Session)

	addr := ":9000"
	go h.Open(addr, s)

	time.Sleep(10 * time.Millisecond)
	u := url.URL{Scheme: "wss", Host: addr, Path: "/websocket"}

	dailer := websocket.Dialer{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c, _, err := dailer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Failed to open websocket: %v", err)
	}
	defer c.Close()

	session := <-s

	assert.IsType(t, *(new(myWS.Session)), session)

}
