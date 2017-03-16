package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hdhauk/gochat-server/chat"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	room := chat.NewRoom("test")
	go room.Start()

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		usr := chat.NewUser(conn, "test-name")
		go usr.Listen()
		fmt.Println("adding user")
		room.AddMember(usr)
	})
	fmt.Println("Staring server")
	http.ListenAndServe(":9000", nil)
	select {}

}
