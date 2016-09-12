package socketio09

import "sync"

/*
AckManager processes response listeners for socketio messages where ack is expected,
and the user put a callback handler in place. A listener is a channel to the handler func.
*/
type AckManager struct {
	// counter is ack ID
	counter     int
	counterLock sync.Mutex

	// int is the counter/ID
	// string is the raw message data
	responseListeners     map[int](chan string)
	responseListenersLock sync.RWMutex
}

/*
getNextID provisions the next ACK id, in a thread safe way. counter starts at 0 by virtue of
being an int, and gets iterated before being returned. So first ackID is 1.
*/
func (a *AckManager) getNextID() int {
	a.counterLock.Lock()
	defer a.counterLock.Unlock()

	a.counter++
	return a.counter
}

func (a *AckManager) addListener(id int, w chan string) {
	a.responseListenersLock.Lock()
	a.responseListeners[id] = w
	a.responseListenersLock.Unlock()
}

func (a *AckManager) removeListener(id int) {
	a.responseListenersLock.Lock()
	delete(a.responseListeners, id)
	a.responseListenersLock.Unlock()
}

/*
getListener returns an ack listener and removes it.
*/
func (a *AckManager) getListener(id int) (chan string, error) {
	a.responseListenersLock.RLock()
	defer a.responseListenersLock.RUnlock()

	listener, exists := a.responseListeners[id]
	if exists {
		delete(a.responseListeners, id)
		return listener, nil
	}
	return nil, ErrorAckListenerNotFound
}
