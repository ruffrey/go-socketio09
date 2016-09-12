package socketio09

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/ruffrey/go-socketio09/spec"
)

const (
	// OnConnect handler
	OnConnect = "connect"
	// OnDisconnect handler
	OnDisconnect = "disconnect"
	// OnError handler
	OnError = "error" // TODO: use this
)

type internalHandler func(c *SocketIOConnection)

type eventEmitter struct {
	messageHandlers     map[string]*HandlerCaller
	messageHandlersLock sync.RWMutex

	internalOnConnect    internalHandler
	internalOnDisconnect internalHandler

	// TODO: add more
}

func (m *eventEmitter) initMethods() {
	m.messageHandlers = make(map[string]*HandlerCaller)
}

/*
On adds a handler to the specified event.
*/
func (m *eventEmitter) On(event string, fn interface{}) (err error) {
	c, err := NewHandlerCaller(fn)
	if err != nil {
		return err
	}

	m.messageHandlersLock.Lock()
	defer m.messageHandlersLock.Unlock()
	m.messageHandlers[event] = c

	return nil
}

func (m *eventEmitter) findHandlerForEvent(event string) (fn *HandlerCaller, exists bool) {
	m.messageHandlersLock.RLock()
	defer m.messageHandlersLock.RUnlock()
	fn, exists = m.messageHandlers[event]
	return fn, exists
}

func (m *eventEmitter) fireEvent(c *SocketIOConnection, event string) {
	if m.internalOnConnect != nil && event == OnConnect {
		m.internalOnConnect(c)
	}
	if m.internalOnDisconnect != nil && event == OnDisconnect {
		m.internalOnDisconnect(c)
	}

	fn, exists := m.findHandlerForEvent(event)
	if !exists {
		return
	}

	fn.callFunc(c, &struct{}{})
}

func (m *eventEmitter) checkAndFireListenersForValidMessage(c *SocketIOConnection, msg *Message) {
	switch msg.Type {
	case spec.Connect:
		m.fireEvent(c, OnConnect)
		return
	case spec.Disconnect:
		CloseChannel(c, m)
		return
	case spec.Event:
		fn, exists := m.findHandlerForEvent(msg.EventName)
		if !exists {
			return
		}
		if !fn.ArgsPresent {
			fn.callFunc(c, &struct{}{})
			return
		}

		data := fn.getArgs()
		err := json.Unmarshal([]byte(msg.Args), &data)

		if err != nil {
			log.Println(err, msg)
			return
		}

		fn.callFunc(c, data)
		return
	case spec.Ack:
		listener, err := c.acks.getListener(msg.AckID)
		if err != nil {
			log.Println(err, msg)
			return
		}
		listener <- msg.Args
		return
	}
}
