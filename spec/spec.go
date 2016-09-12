package spec

// Some socket.io protocol message types are missing; this is because Respoke does
// not use those message types.

const (
	// Disconnect signals disconnection from the server. Server->Client.
	Disconnect = "0"
	// Connect comes from the server as the first message, indicating all is well and ready.
	Connect = "1"
	// Heartbeat means the socket is still alive. Bidirectional.
	Heartbeat = "2"
	// TextMessage is a regular message (TODO: implement this)
	TextMessage = "3"
	// JSONMessage is a regular message that is JSON (TODO: implement this)
	JSONMessage = "4"
	// Event has a name and attached data. Bidirectional.
	Event = "5"
	// Ack acknowledges a request, and maybe has some data attached to it. Server->Client.
	Ack = "6"
	// Error is apparently an error (TODO: implement this)
	Error = "7"
	// Noop mean dont do anything, I guess
	Noop = "8"
)
