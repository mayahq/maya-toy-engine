package main

import (
	"github.com/cascades-fbp/cascades/library"
)

var registryEntry = &library.Entry{
	Description: "Cascades-CAF Context component that evaluates the context based on input Properties and Contexts",
	Elementary:  true,
	Inports: []library.EntryPort{
		library.EntryPort{
			Name:        "DATA",
			Type:        "json",
			Description: "Data input: Properties or Contexts",
			Required:    true,
		},
		library.EntryPort{
			Name:        "TMPL",
			Type:        "string",
			Description: "Context template",
			Required:    true,
		},
	},
	Outports: []library.EntryPort{
		library.EntryPort{
			Name:        "UPDATE",
			Type:        "json",
			Description: "Port for Context state when updated",
			Required:    false,
		},
		library.EntryPort{
			Name:        "MATCH",
			Type:        "json",
			Description: "Port for Context state when matched",
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
