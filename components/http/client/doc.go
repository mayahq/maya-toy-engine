package main

import (
	"github.com/cascades-fbp/cascades/library"
)

var registryEntry = &library.Entry{
	Description: "Multi-purpose HTTP client component",
	Elementary:  true,
	Inports: []library.EntryPort{
		library.EntryPort{
			Name:        "REQ",
			Type:        "json",
			Description: "JSON object describing the HTTP request",
			Required:    true,
		},
	},
	Outports: []library.EntryPort{
		library.EntryPort{
			Name:        "RESP",
			Type:        "json",
			Description: "Response JSON object defined in utils of HTTP components library",
			Required:    false,
		},
		library.EntryPort{
			Name:        "BODY",
			Type:        "string",
			Description: "Body of the response",
			Required:    false,
		},
		library.EntryPort{
			Name:        "ERR",
			Type:        "string",
			Description: "Error port for errors while performing requests",
			Required:    false,
		},
	},
}
