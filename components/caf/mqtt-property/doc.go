package main

import (
	"github.com/cascades-fbp/cascades/library"
)

var registryEntry = &library.Entry{
	Description: "Subscribes to a given MQTT broker's topic and outputs all incoming message to output port",
	Inports: []library.EntryPort{
		library.EntryPort{
			Name:        "OPTIONS",
			Type:        "string",
			Description: "MQTT broker connection options (e.g. tcp://username:password@127.0.0.1:1883/defaultTopic)",
			Required:    true,
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
			Name:        "ERR",
			Type:        "string",
			Description: "Error port for errors reporting",
			Required:    false,
		},
	},
}
