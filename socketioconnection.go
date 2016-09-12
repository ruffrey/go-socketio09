package socketio09

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/ruffrey/go-socketio09/spec"
)

const queueMaxSize = 500

/*
SocketIOConnection is a socket.io connection handler object.
*/
type SocketIOConnection struct {
	conn WebsocketConnection

	outboundMQ chan string

	alive     bool
	aliveLock sync.Mutex

	acks AckManager

	requestHeader http.Header
}

/*
initChannel create channel, map, and set active
*/
func (c *SocketIOConnection) initChannel() {
	//TODO: queueMaxSize from constant to server or client variable
	c.outboundMQ = make(chan string, queueMaxSize)
	c.acks.responseListeners = make(map[int](chan string))
	c.alive = true
}

/*
IsActive checks that the socket connection is still alive
*/
func (c *SocketIOConnection) IsActive() bool {
	return c.alive
}

/*
CloseChannel closes the respoke signaling channel
*/
func CloseChannel(c *SocketIOConnection, m *eventEmitter, args ...interface{}) error {
	c.aliveLock.Lock()
	defer c.aliveLock.Unlock()

	if !c.alive {
		//already closed
		return nil
	}

	c.conn.Close()
	c.alive = false

	// clean handleOutboundMessages
	for len(c.outboundMQ) > 0 {
		<-c.outboundMQ
	}
	c.outboundMQ <- spec.Disconnect + "::"

	m.fireEvent(c, OnDisconnect)

	overfloodedLock.Lock()
	delete(overflooded, c)
	overfloodedLock.Unlock()

	return nil
}

// handleInboundMessages takes incoming message frames from the web socket and transforms
// them into a meaningful type (json, for example) then bubbles that up to any userland handlers.
func handleInboundMessages(c *SocketIOConnection, m *eventEmitter) error {
	for {
		pkg, err := c.conn.GetNextMsg()
		if err != nil {
			return CloseChannel(c, m, err)
		}
		msg, err := DecodeInboundMessage(pkg)
		if err != nil {
			CloseChannel(c, m, ErrorProtocolReceivedInvalidPacket)
			return err
		}

		switch msg.Type {
		case spec.Noop:
			c.outboundMQ <- spec.Heartbeat + "::"
		default:
			go m.checkAndFireListenersForValidMessage(c, msg)
		}
	}
}

var overflooded = make(map[*SocketIOConnection]struct{})
var overfloodedLock sync.Mutex

// AmountOfOverflooded indicates how many bytes we are over the frame limit
func AmountOfOverflooded() int64 {
	overfloodedLock.Lock()
	defer overfloodedLock.Unlock()

	return int64(len(overflooded))
}

/*
handleOutboundMessages waits for outgoing messages, then sends the messages from this
SocketIOConnection to the web socket transport.
*/
func handleOutboundMessages(c *SocketIOConnection, m *eventEmitter) error {
	for {
		outBufferLen := len(c.outboundMQ)
		if outBufferLen >= queueMaxSize-1 {
			return CloseChannel(c, m, ErrorSocketOverflood)
		} else if outBufferLen > int(queueMaxSize/2) {
			overfloodedLock.Lock()
			overflooded[c] = struct{}{}
			overfloodedLock.Unlock()
		} else {
			overfloodedLock.Lock()
			delete(overflooded, c)
			overfloodedLock.Unlock()
		}
		// pull the message off the outbound channel and write it to the web socket
		msg := <-c.outboundMQ
		if msg[0:1] == spec.Disconnect {
			return nil
		}

		err := c.conn.WriteMsg(msg)
		if err != nil {
			return CloseChannel(c, m, err)
		}

	}
}

/*
send will send an outgoing message packet to the SocketIOConnection.
*/
func send(msg *Message, c *SocketIOConnection, args interface{}) error {
	if args != nil {
		json, err := json.Marshal(&args)
		if err != nil {
			return err
		}

		msg.Args = string(json)
	}

	command, err := EncodeOutboundMessage(msg)
	if err != nil {
		return err
	}

	if len(c.outboundMQ) == queueMaxSize {
		return ErrorSocketOverflood
	}

	c.outboundMQ <- command

	return nil
}

/*
Emit creates a packet based on given data and sends it
*/
func (c *SocketIOConnection) Emit(method string, args interface{}) error {
	msg := &Message{
		Type:      spec.Event,
		EventName: method,
	}
	return send(msg, c, args)
}

/*
EmitWithAck creates an ack frame, then sends it AND waits for a response.
*/
func (c *SocketIOConnection) EmitWithAck(method string, args interface{}) (string, error) {
	timeout := c.conn.transport.ReceiveTimeout
	msg := &Message{
		Type:      spec.Event,
		AckID:     c.acks.getNextID(),
		EventName: method,
	}

	listener := make(chan string)
	c.acks.addListener(msg.AckID, listener)

	err := send(msg, c, args)
	if err != nil {
		c.acks.removeListener(msg.AckID)
	}

	select {
	case result := <-listener:
		return result, nil
	case <-time.After(timeout):
		c.acks.removeListener(msg.AckID)
		return "", ErrorSendTimeout
	}
}

/*
heartbeatService sends ping messages for keeping connection alive
*/
func heartbeatService(c *SocketIOConnection) {
	for {
		time.Sleep(c.conn.transport.HeartbeatInterval)
		if !c.IsActive() {
			return
		}

		c.outboundMQ <- spec.Heartbeat + "::"
	}
}
