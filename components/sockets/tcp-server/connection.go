package main

import (
	"net"

	uuid "github.com/nu7hatch/gouuid"
)

// Connection
// Wraps actual TCP connection and input channel
type Connection struct {
	Id      string
	TCPConn *net.TCPConn
	Input   chan []byte
}

// Constructor function for connection
func NewConnection(tcpConn *net.TCPConn, in chan []byte) *Connection {
	id, _ := uuid.NewV4()
	c := &Connection{
		Id:      id.String(),
		TCPConn: tcpConn,
		Input:   in,
	}
	return c
}

// Closes TCP connection & input channel
func (self *Connection) Close() {
	self.TCPConn.Close()
	close(self.Input)
}
