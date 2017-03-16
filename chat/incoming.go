package chat

// IncomingMsg is the data structure clients send to the server.
type IncomingMsg struct {
	RoomID    string `json:"room-id"`
	MsgTxt    string `json:"msg-txt"`
	Language  string `json:"lang"`
	AttData   []byte `json:"att-data"`
	AttFormat string `json:"att-format"`
	Priority  int    `json:"pri"`
	SenderID  string `json:"-"`
}
