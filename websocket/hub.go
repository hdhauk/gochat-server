package websocket

import (
	"log"
	"net/http"
	"os"
)

// Hub represnts an endpoint to which external clients my connect to in order
// to establish a secure websocket with the server.
type Hub struct {
	logger *log.Logger
}

// NewHub returns a Hub.
func NewHub() *Hub {
	return &Hub{
		logger: log.New(os.Stderr, "[websocket] ", log.Ltime|log.Lshortfile),
	}
}

// Open start the enpoint and return any sessions on the returned channel.
func (h *Hub) Open(addr string, ret chan<- Session) {
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		s, err := NewSession(w, r)
		if err != nil {
			h.logger.Printf("[WARN] Failed to establish websocket session: %v\n", err.Error())
			return
		}
		ret <- s
	})
	http.ListenAndServeTLS(addr, "server.crt", "server.key", nil)
}
