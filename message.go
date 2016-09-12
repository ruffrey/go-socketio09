package socketio09

/*
Message is an internal construct representing a socket.io protocol 0.9 frame message
container (inbound or outbound).
*/
type Message struct {
	// Type is the socket.io 0.9 specificiation message type
	Type string
	// AckID will be present on an ack response only
	AckID int
	// EventName is for events (Type=5), to or from, where this is the `"name": "some event name"`
	EventName string
	// Args will be a JSON array in socket.io protocol
	Args string
}
