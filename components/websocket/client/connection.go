package main

import (
	"log"

	"golang.org/x/net/websocket"
)

// Connection container
type Connection struct {
	WS      *websocket.Conn
	Send    chan []byte
	Receive chan []byte
}

// NewConnection constructs Connection structures
func NewConnection(wsConn *websocket.Conn) *Connection {
	connection := &Connection{
		WS:      wsConn,
		Send:    make(chan []byte, 256),
		Receive: make(chan []byte, 256),
	}
	return connection
}

// Reader is reading from socket and writes to Receive channel
func (c *Connection) Reader() {
	for {
		var data []byte
		err := websocket.Message.Receive(c.WS, &data)
		if err != nil {
			log.Println("Reader: error receiving message")
			continue
		}
		log.Println("Reader: received from websocket", data)
		c.Receive <- data
	}
	c.WS.Close()
}

// Writer is reading from channel and writes to socket
func (c *Connection) Writer() {
	for data := range c.Send {
		log.Println("Writer: will send to websocket", data)
		err := websocket.Message.Send(c.WS, data)
		if err != nil {
			break
		}
	}
	c.WS.Close()
}

// Close closes the channels and socket
func (c *Connection) Close() {
	close(c.Send)
	close(c.Receive)
	c.WS.Close()
}
