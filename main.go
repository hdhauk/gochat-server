package main

import "github.com/hdhauk/gochat-server/chat"

func main() {

	room := chat.NewRoom("test")
	room.Start()
	select {}
}
