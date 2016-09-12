# Golang Socket.IO 0.9.x Client

Socket.IO is protocol for framing and acking messages over a web socket. The
[Socket.IO server](/socketio/socket.io/) is written in Node.js. 0.9 is an older version,
with a different protocol than v1 and beyond.

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

## Implemented

- emit events, and receive ack
- listen for events

## Unimplemented

- send text, json messages
- emit events to "message endpoints"

## Known Issues

- protocol level bugs might get buried

## Development

Before working on the `socketio09/` portion of this library, it is best to read the short
spec at `SOCKETIO-SPEC-0.9.x.md`.

## License

MIT. Copyright (c) Jeff H. Parrish 2016. See LICENSE file in this repository.
