API is unstable! Beware.

# Golang Socket.IO 0.9.x Client

Socket.IO is a protocol for framing and acking messages over a web socket.
[Socket.IO server](https://github.com/socketio/socket.io/) is implemented in Node.js. Socket.IO 0.9 is an older
version, with a different protocol than Socket.IO 1 and beyond.

This library implements a limited subset of the Socket.IO spec. It is a work in progress.
There were not any other Socket.IO 0.9 compatible libraries at the time I needed this,
which written in go.

## Examples

See the `examples/` folder. Since this client library requires a Socket.IO 0.9 server, there
is one in the `examples/` folder. It requires Node.js (>=0.10) and npm (>=2) to be installed.
To run the example server:

```bash
cd examples
npm install
npm start
```

To run the example go client:

```bash
cd examples
go run basic-client.go
```

## Implemented

- emit json events, and receive json ack
- listen for events

## Unimplemented

- send regular text, json messages (not emitted events) to other socket.io endpoints
- handling some message types: 3, 4, 7

## Known Issues

- protocol level bugs might get buried

## License

MIT. Copyright (c) Jeff H. Parrish 2016. See LICENSE file in this repository.
