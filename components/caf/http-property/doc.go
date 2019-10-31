package main

import (
	"github.com/cascades-fbp/cascades/library"
)

var registryEntry = &library.Entry{
	Description: "Cascades-CAF HTTP property component that fills properties with data pulled from HTTP endpoints",
	Elementary:  true,
	Inports: []library.EntryPort{
		library.EntryPort{
			Name:        "INT",
			Type:        "string",
			Description: "Pulling interval",
			Required:    true,
		},
		library.EntryPort{
			Name:        "REQ",
			Type:        "json",
			Description: "JSON object describing the HTTP request",
			Required:    false,
		},
		library.EntryPort{
			Name:        "TMPL",
			Type:        "json",
			Description: "JSON object describing the property template",
			Required:    true,
		},
	},
	Outports: []library.EntryPort{
		library.EntryPort{
			Name:        "PROP",
			Type:        "json",
			Description: "Property filled with data",
			Required:    false,
		},
		library.EntryPort{
			Name:        "RESP",
			Type:        "json",
			Description: "HTTP response object defined in utils of HTTP components library",
			Required:    false,
		},
		library.EntryPort{
			Name:        "BODY",
			Type:        "string",
			Description: "Body of the HTTP response",
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
