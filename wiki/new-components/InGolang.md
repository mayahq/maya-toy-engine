The core components shipped with Cascades are written in Go and use some utility packages to simplify the development. Therefore creating new components in Go is the easy as you can a lot of examples to learn from.

## Prerequisites

The communication between components is implemented using ZeroMQ, which needs to be installed on your machine as described in the [ZeroMQ Golang Bindings documentation](http://zeromq.org/bindings:go). You need to use `github.com/pebbe/zmq4` package, which can installed by executing `go get github.com/pebbe/zmq4` command.

Additionally (if you have installed Cascades binary distribution) you need to installed the Cascades utility packages by executing `github.com/cascades-fbp/cascades/components/utils` and `github.com/cascades-fbp/cascades/runtime` commands.

## Complete example of a Passthru component

The source code of the _Passthru_ component is available from the [Cascades Github repository](https://github.com/cascades-fbp/cascades/tree/master/components/core/passthru).

## Using Cascades utility packages

Include required packages:

```
import (
	"github.com/cascades-fbp/cascades/components/utils"
	"github.com/cascades-fbp/cascades/runtime"
	zmq "github.com/pebbe/zmq4"
)
```

Create & open ports:

```
func openPorts() {
	inPort, err = utils.CreateInputPort(*inputEndpoint)
	utils.AssertError(err)

	outPort, err = utils.CreateOutputPort(*outputEndpoint)
	utils.AssertError(err)
}
```

Set up interrupting handler & input port monitoring:

```
ch := utils.HandleInterruption()
err = runtime.SetupShutdownByDisconnect(inPort, "passthru.in", ch)
utils.AssertError(err)
```

Receive an IP from input port:

```
ip, err := inPort.RecvMessageBytes(0)
if err != nil {
	log.Println("Error receiving message:", err.Error())
	// handle error
}
```

Validate the received message if it is a proper IP:

```
if !runtime.IsValidIP(ip) {
	// handle error
}
```

Forward IP to the output port:

```
outPort.SendMessage(ip)
```