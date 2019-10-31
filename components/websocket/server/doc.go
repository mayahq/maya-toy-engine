package main

import (
	"github.com/cascades-fbp/cascades/library"
)

var registryEntry = &library.Entry{
	Description: "Creates a Websocket server and binds to an address/port received from options",
	Elementary:  true,
	Inports: []library.EntryPort{
		library.EntryPort{
			Name:        "OPTIONS",
			Type:        "string",
			Description: "Configures the Websocket server with address to bind to (i.e. 0.0.0.0:5000)",
			Required:    true,
		},
		library.EntryPort{
			Name:        "IN",
			Type:        "json",
			Description: "Input port for receiving IPs and forwarding them to the corresponding connection",
			Required:    true,
		},
	},
	Outports: []library.EntryPort{
		library.EntryPort{
			Name:        "OUT",
			Type:        "json",
			Description: "Output port for sending IPs with data received from a specific connection",
			Required:    true,
		},
	},
}
