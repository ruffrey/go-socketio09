package socketio09

import (
	"log"
	"strconv"

	"github.com/ruffrey/go-socketio09/spec"
)

/*
EncodeOutboundMessage is used to make an outgoing message. A *Message is transformed into a socket.io frame
message string, to be sent over a web socket connection.

At this point in time, we are only going to be sending hearbeats (2) and events (5).
*/
func EncodeOutboundMessage(m *Message) (msg string, err error) {
	log.Println("encoding:", m)
	switch m.Type {
	case spec.Heartbeat:
		msg = spec.Heartbeat
		return msg, nil
	case spec.Event:
		msg = spec.Event
		if m.AckID != 0 {
			msg += ":" + strconv.Itoa(m.AckID) + `+::{"name":"` + m.EventName + `","args":[` + m.Args + `]}`
		} else {
			msg += `:::{"name":"` + m.EventName + `","args":[` + m.Args + `]}`
		}
		return msg, nil
	}

	// this should not happen
	return msg, ErrorProtocolUnexpectedOutboundMessageType
}
