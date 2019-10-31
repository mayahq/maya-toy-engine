package main

import (
	"github.com/cascades-fbp/cascades/library"
)

var registryEntry = &library.Entry{
	Description: "Configures a Websocket client to connect to a given server, forward incoming data to server, output receive server data to output port",
	Elementary:  true,
	Inports: []library.EntryPort{
		library.EntryPort{
			Name:        "OPTIONS",
			Type:        "string",
			Description: "Configures the Websocket server to connect to (i.e. ws://localhost:5000/ws)",
			Required:    true,
		},
		library.EntryPort{
			Name:        "IN",
			Type:        "json",
			Description: "Input port for receiving IPs and forwarding them to server connection",
			Required:    true,
		},
	},
	Outports: []library.EntryPort{
		library.EntryPort{
			Name:        "OUT",
			Type:        "json",
			Description: "Output port for sending IPs with data received from a server",
			Required:    true,
		},
	},
}
