package socketio09

import (
	"io/ioutil"
	"math"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const wsDefaultBufferSize = 1024 * 32

// WebsocketTransport is an object representing defaults for a socket.io transport type
type WebsocketTransport struct {
	HeartbeatInterval time.Duration
	HeartbeatTimeout  time.Duration
	ReceiveTimeout    time.Duration
	SendTimeout       time.Duration

	// TODO: implement this
	ConnectionCloseTimeout time.Duration

	BufferSize int
}

// WebsocketConnection represents the web socket client connection
type WebsocketConnection struct {
	socket    *websocket.Conn
	transport *WebsocketTransport
}

// GetNextMsg reads the latest buffered message into a string
func (wsc *WebsocketConnection) GetNextMsg() (text string, err error) {
	wsc.socket.SetReadDeadline(time.Now().Add(wsc.transport.ReceiveTimeout))
	msgType, reader, err := wsc.socket.NextReader()
	if err != nil {
		return "", err
	}

	// support only text messages exchange
	if msgType != websocket.TextMessage {
		return "", ErrorTransportOnlySupportsText
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", ErrorTransportBufferError
	}
	text = string(data)

	if len(text) == 0 {
		return "", ErrorTransportEmptyPacket
	}

	return text, nil
}

// WriteMsg writes the exact message to a web socket (should be in protocol format already).
func (wsc *WebsocketConnection) WriteMsg(message string) error {
	wsc.socket.SetWriteDeadline(time.Now().Add(wsc.transport.SendTimeout))
	writer, err := wsc.socket.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	if _, err := writer.Write([]byte(message)); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return nil
}

// Close just calls close on the underlying websocket
func (wsc *WebsocketConnection) Close() {
	wsc.socket.Close()
}

// GetPingInfo pulls ping information from a web socket connection
func (wsc *WebsocketConnection) GetPingInfo() (interval, timeout time.Duration) {
	return wsc.transport.HeartbeatInterval, wsc.transport.HeartbeatTimeout
}

/*
SocketIOClient is a Socket.io client container
*/
type SocketIOClient struct {
	eventEmitter
	SocketIOConnection
}

/*
Close will properly terminate the web socket connection according to socket.io's preferences.
*/
func (c *SocketIOClient) Close() {
	CloseChannel(&c.SocketIOConnection, &c.eventEmitter)
}

/*
Connect will initiate a socketio style web socket connection. fullURL is the full web socket fullURL.
If you wish to pass additional querystring params, feel free to do so.

	conn, err := socketio09.Connect("https://127.0.0.1:4500/socket.io/1?__sails_io_sdk_version=0.10.0")
*/
func (wst *WebsocketTransport) Connect(fullURL string) (client *SocketIOClient, err error) {
	client = &SocketIOClient{}
	urlWithToken, _ := url.Parse(fullURL)

	// golang url does not support ws:// or wss://, so we hack it later during web socket connect
	var wsScheme string
	if urlWithToken.Scheme == "https" {
		wsScheme = "wss"
	} else {
		wsScheme = "ws"
	}

	hr, err := handshake(fullURL)
	if err != nil {
		return client, err
	}

	wst.HeartbeatTimeout = time.Duration(hr.heartbeatTimeout) * time.Second
	// heartbeat in 3/4 the timeout time
	wst.HeartbeatInterval = time.Duration(math.Floor(float64(hr.heartbeatTimeout/2))) * time.Second
	wst.ConnectionCloseTimeout = time.Duration(hr.connectionTimeout) * time.Second
	// not sure if these next two are right, or apply to socket.io 0.9
	wst.SendTimeout = time.Duration(hr.heartbeatTimeout) * time.Second
	wst.ReceiveTimeout = time.Duration(hr.heartbeatTimeout) * time.Second

	urlWithToken.Path = "/socket.io/1/websocket/" + hr.token
	webSocketURLWithToken := strings.Replace(urlWithToken.String(), urlWithToken.Scheme, wsScheme, 1)
	dialer := websocket.Dialer{}
	socket, _, err := dialer.Dial(webSocketURLWithToken, nil)
	if err != nil {
		return client, err
	}

	client.conn = WebsocketConnection{socket, wst}
	client.initChannel()
	client.initMethods()
	go handleInboundMessages(&client.SocketIOConnection, &client.eventEmitter)
	go handleOutboundMessages(&client.SocketIOConnection, &client.eventEmitter)
	go heartbeatService(&client.SocketIOConnection)

	return client, nil
}

/*
NewConnection returns a new socketio websocket connection transport with default timings
and buffer size. The next step should be to call `wst.Connect(url)`
*/
func NewConnection() (wst *WebsocketTransport) {
	return &WebsocketTransport{
		BufferSize: wsDefaultBufferSize,
	}
}
