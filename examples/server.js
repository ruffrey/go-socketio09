var PORT = process.env.PORT || 4500

var io = require('socket.io').listen(PORT, function () {
  console.log('\n  Socket.IO server is listening at http://localhost:' + PORT + '/socket.io/1')
  console.log('  You can change the port with the environment variable `PORT`.')
  console.log(
    '  Clients can emit `test` event and the data will be relayed back over the `test` event.'
  )
  console.log('  The `time` event will fire every 5 seconds.')
  console.log('\n  Hit control+C to stop\n')
})

// Send current time to all connected clients
function sendTime() {
    io.sockets.emit('time', { time: new Date().toJSON() })
}

setInterval(sendTime, 6000)

// Emit welcome message on connection
io.sockets.on('connection', function (socket) {
    console.log('Socket connected', socket.id)

    socket.emit('welcome', {
      message: 'Welcome! Emit the `test` event to recieve your',
      id: socket.id
    })
    socket.on('test', function (data, callback) {
      if (callback) {
        callback(data)
      } else {
        socket.emit('test', data);
      }
    })
    socket.on('disconnect', function () {
      console.log('Socket disconnected', socket.id)
    })
});
