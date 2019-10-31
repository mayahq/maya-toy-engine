package main

import (
	"github.com/sibeshkar/maya-engine/library"
)

var registryEntry = &library.Entry{
	Description: "Packets dropper. Simply consumes IPs from the input port and 'deletes' them.",
	Elementary:  true,
	Inports: []library.EntryPort{
		library.EntryPort{
			Name:        "IN",
			Type:        "all",
			Description: "Input port for IP",
			Required:    true,
		},
	},
}
