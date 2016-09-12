package socketio09

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/ruffrey/go-socketio09/spec"
)

type socketioEventMessage struct {
	Name string      `json:"name"`
	Args interface{} `json:"args"`
}

/*
Get ack id of current packet, if present.
Message strings start with the message type, followed by a `:`, so we will always start at index 2.
If there is
*/
func getAckAndDataFromIncomingMessageText(text string) (ack int, restText string, err error) {
	// example of shortest possible packet is `1::`
	if len(text) < 3 {
		return 0, "", ErrorProtocolReceivedInvalidPacket
	}

	if text[0:1] == spec.Event {
		restText = text[4:]
	}
	if text[0:1] == spec.Ack {
		argsStartPosition := strings.IndexByte(text, '+')
		if argsStartPosition == -1 {
			return ack, restText, ErrorProtocolReceivedInvalidPacket
		}

		ack, err = strconv.Atoi(text[4:argsStartPosition])
		if err != nil {
			return 0, restText, err
		}
		restText = text[argsStartPosition:]
	}

	return ack, restText, nil
}

// This may no longer be necessary
func getInboundMessageType(data string) (string, error) {
	if len(data) == 0 {
		return "", ErrorProtocolUnexpectedInboundMessageType
	}
	switch data[0:1] {
	case spec.Connect:
		return spec.Connect, nil
	case spec.Disconnect:
		return spec.Disconnect, nil
	case spec.Heartbeat:
		return spec.Heartbeat, nil
	case spec.Event:
		if len(data) == 1 {
			return "", ErrorProtocolReceivedInvalidPacket
		}
		return spec.Event, nil
	case spec.Ack:
		return spec.Ack, nil
	}
	return "", ErrorProtocolUnexpectedInboundMessageType
}

// DecodeInboundMessage takes the socketio encoded frame and turns it into a client *Message
func DecodeInboundMessage(data string) (*Message, error) {
	var err error
	m := &Message{}

	m.Type, err = getInboundMessageType(data)
	if err != nil {
		return nil, err
	}

	if m.Type == spec.Connect || m.Type == spec.Disconnect ||
		m.Type == spec.Heartbeat || m.Type == spec.Noop {
		return m, nil
	}

	ackID, rest, err := getAckAndDataFromIncomingMessageText(data)
	if err != nil {
		log.Println("inbound msg decode failed:", err)
		return m, err
	}
	if m.Type == spec.Event {
		msgJSON := socketioEventMessage{}
		err := json.Unmarshal([]byte(rest), &msgJSON)
		if err != nil {
			return m, err
		}
		m.EventName = msgJSON.Name
		argsAsStringAgain, err := json.Marshal(msgJSON.Args)
		if err != nil {
			return m, err
		}
		m.Args = string(argsAsStringAgain)
	}
	if m.Type == spec.Ack {
		m.Args = rest
		m.AckID = ackID
	}
	return m, nil
}
