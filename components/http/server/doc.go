package main

import "github.com/cascades-fbp/cascades/library"

var registryEntry = &library.Entry{
	Description: "Create a HTTP server and binds to an address/port received from options",
	Elementary:  true,
	Inports: []library.EntryPort{
		library.EntryPort{
			Name:        "OPTIONS",
			Type:        "string",
			Description: "Configuration port to pass IP with TCP endpoint in the format i.e. 127.0.0.1:8080",
			Required:    true,
		},
		library.EntryPort{
			Name:        "IN",
			Type:        "json",
			Description: "Input port for receiving responses in predefined JSON format",
			Required:    true,
		},
	},
	Outports: []library.EntryPort{
		library.EntryPort{
			Name:        "OUT",
			Type:        "json",
			Description: "Output port for emitting requests in predefined JSON format",
			Required:    true,
		},
	},
}
