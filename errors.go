package socketio09

import "errors"

var (
	/* Client Errors */

	// ErrorAckListenerNotFound indicates firing the listener on an ACK failed, due to the
	// listener being missing, or maybe already called
	ErrorAckListenerNotFound = errors.New("ACK listener not found")
	// ErrorCallerShouldBeTypeFunc is an error
	ErrorCallerShouldBeTypeFunc = errors.New("type error: expected a func in handler arg")
	// ErrorCallerShouldHaveTwoArgs is an error
	ErrorCallerShouldHaveTwoArgs = errors.New("func fn must have one or two args only")
	// ErrorCallerFunctionReturnsTooMuch is an error
	ErrorCallerFunctionReturnsTooMuch = errors.New("func fn must return one value")
	// ErrorSendTimeout is an error
	ErrorSendTimeout = errors.New("Timeout")
	// ErrorSocketOverflood is an error
	ErrorSocketOverflood = errors.New("Socket is flooded")

	/* Protocol Errors */

	// ErrorProtocolUnexpectedInboundMessageType is an error
	ErrorProtocolUnexpectedInboundMessageType = errors.New("Protocol Error: unexpected inbound message type")
	// ErrorProtocolUnexpectedOutboundMessageType is an error
	ErrorProtocolUnexpectedOutboundMessageType = errors.New("Protocol Error: unexpected outbound message type")
	// ErrorProtocolReceivedInvalidPacket is an error
	ErrorProtocolReceivedInvalidPacket = errors.New("Protocol Error: invalid packet type received")

	/* Web Socket Errors */

	// ErrorTransportOnlySupportsText is an error
	ErrorTransportOnlySupportsText = errors.New("Received non-TextMessage at websocket")
	// ErrorTransportBufferError is an error
	ErrorTransportBufferError = errors.New("Buffer error (buffer may not be a buffer)")
	// ErrorTransportEmptyPacket is an error
	ErrorTransportEmptyPacket = errors.New("Web socket message is empty and that is not allowed")
	// ErrorHTTPUpgradeFailed is an error
	ErrorHTTPUpgradeFailed = errors.New("Failure during HTTP upgrade attempt")
)
