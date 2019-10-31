cascades-caf
============

A Context-Awareness Framework as components for Cascades

## Examples
Using HttpProperty with the [system agent](https://github.com/patchwork-toolkit/agent-examples/tree/master/system)

```
'3s' -> INT HttpProperty(caf/http-property)

'{
	"url": "http://localhost:9000/rest/System/PS", 
	"method": "GET",
	"content-type": "application/json"
}' -> REQ HttpProperty

'{
	"type": "float",	
	"name": "ProcessCount",
	"template": "{{ .count }}",
	"group": "SystemA"
}' -> TMPL HttpProperty


HttpProperty PROP -> IN Log(core/console)
```