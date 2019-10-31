package main

import (
	"github.com/sibeshkar/maya-engine/library"
)

var registryEntry = &library.Entry{
	Description: "Creates a TCP server and binds to an address/port received from options",
	Elementary:  true,
	Inports: []library.EntryPort{
		library.EntryPort{
			Name:        "OPTIONS",
			Type:        "all",
			Description: "Configuration port to pass IP with TCP endpoint in the format tcp://ip:port",
			Required:    true,
		},
		library.EntryPort{
			Name:        "IN",
			Type:        "all",
			Description: "Input port for receiving IPs",
			Required:    true,
		},
	},
	Outports: []library.EntryPort{
		library.EntryPort{
			Name:        "OUT",
			Type:        "all",
			Description: "Output port for sending IPs",
			Required:    true,
		},
	},
}
