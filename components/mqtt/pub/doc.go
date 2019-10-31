package main

import (
	"github.com/cascades-fbp/cascades/library"
)

var registryEntry = &library.Entry{
	Description: "Publishes received payload to MQTT broker configured using options port",
	Elementary:  true,
	Inports: []library.EntryPort{
		library.EntryPort{
			Name:        "IN",
			Type:        "all",
			Description: "Data to publish (single IPs go to default options topic, substreams should include topic as first IP)",
			Required:    true,
		},
		library.EntryPort{
			Name:        "OPTIONS",
			Type:        "string",
			Description: "MQTT broker connection options (e.g. tcp://username:password@127.0.0.1:1883/defaultTopic)",
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
