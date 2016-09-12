package main

import (
	"log"
	"runtime"
	"time"

	"github.com/ruffrey/go-socketio09"
)

// Channel is a thing
type Channel struct {
	Channel string `json:"channel"`
}

// Message is a thing
type Message struct {
	Text     string `json:"text,omitempty"`
	Message  string `json:"message,omitempty"`
	SocketID string `json:"id,omitempty"`
	Time     string `json:"time,omitempty"`
}

func emitTest(c *socketio09.SocketIOClient) {
	err := c.Emit("test", Message{
		Text: "A quick message",
	})
	if err != nil {
		log.Println("Emit failed:", err)
	}
}

func emitTestWithAck(c *socketio09.SocketIOClient) {
	response, err := c.EmitWithAck("test", Message{
		Text: "Yell if you got this",
	})
	if err != nil {
		log.Println("EmitWithAck failed:", err)
	} else {
		log.Println("Got ack result", response)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	wst := socketio09.NewConnection()
	c, err := wst.Connect("http://localhost:4500/socket.io/1")
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("test", func(h *socketio09.SocketIOConnection, args []Message) {
		log.Println("test message received: ", args)
	})
	if err != nil {
		log.Fatal(err)
	}
	err = c.On("welcome", func(h *socketio09.SocketIOConnection, args []Message) {
		log.Println("welcome received: ", args)
	})
	if err != nil {
		log.Fatal(err)
	}
	err = c.On("time", func(h *socketio09.SocketIOConnection, args []Message) {
		log.Println("time received: ", args)
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("connect", func(h *socketio09.SocketIOConnection) {
		log.Println("Connected")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("disconnect", func(h *socketio09.SocketIOConnection, args interface{}) {
		log.Println("Disconnected", args)
	})
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	go emitTest(c)
	go emitTestWithAck(c)

	time.Sleep(10 * time.Second)
	c.Close()

	log.Println("Clean exit")
}
