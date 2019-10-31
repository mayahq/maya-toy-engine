package main

import (
	"github.com/cascades-fbp/cascades/library"
)

var registryEntry = &library.Entry{
	Description: "Receives a series JSON structure to be written directly to a configured (through options port) InfluxDB",
	Elementary:  true,
	Inports: []library.EntryPort{
		library.EntryPort{
			Name:        "IN",
			Type:        "json",
			Description: "Event structure in JSON",
			Required:    true,
		},
		library.EntryPort{
			Name:        "OPTIONS",
			Type:        "string",
			Description: "InfluxDB connection options",
			Required:    true,
		},
	},
	Outports: []library.EntryPort{
		library.EntryPort{
			Name:        "ERR",
			Type:        "string",
			Description: "Error port for errors reporting",
			Required:    false,
		},
	},
}
