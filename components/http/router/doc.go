package main

import "github.com/cascades-fbp/cascades/library"

var registryEntry = &library.Entry{
	Description: `Matches a URI and method from incoming JSON requests from http/server and forwards it either
to matching or failing output ports. Each PATTERN[index] port is routed to the corresponding SUCCESS[index]
or a single FAIL output port.`,
	Elementary: true,
	Inports: []library.EntryPort{
		library.EntryPort{
			Name:        "PATTERN",
			Type:        "string",
			Description: "Input array port for matching pattern configuration",
			Required:    true,
			Addressable: true,
		},
		library.EntryPort{
			Name:        "REQUEST",
			Type:        "json",
			Description: "Input port for JSON requests in predefined format",
			Required:    true,
		},
	},
	Outports: []library.EntryPort{
		library.EntryPort{
			Name:        "SUCCESS",
			Type:        "json",
			Description: "Output array port for emitting JSON responses with matched pattern",
			Required:    true,
			Addressable: true,
		},
		library.EntryPort{
			Name:        "FAIL",
			Type:        "json",
			Description: "Output port for emitting responses when URI/method didn't match",
			Required:    true,
		},
	},
}
