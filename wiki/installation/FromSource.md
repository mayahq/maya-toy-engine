## Prerequisites

The Cascades framework requires the latest stable Go release. If you don't have Go installed, read the official [Go installation guide](http://golang.org/doc/install). 

Then you need to [download](http://zeromq.org/intro:get-the-software) ZeroMQ or install using the package manager of your choice (e.g. `brew install zmq` under OSX).

## Building runtime

Create a folder for your installation and change to it (all subsequent commands should be executed from this folder):

```
mkdir cascades-fbp
```

Checkout the sources from Github treating the current folder as a GOPATH directory:

```
GOPATH=`pwd` go get -d github.com/cascades-fbp/cascades
```

Compile the runtime CLI:

```
GOPATH=`pwd` go install github.com/cascades-fbp/cascades/cmd/...
```

## Building core components

Create a directory for the components:

```
mkdir components
```

and then for every component you want to build execute the following command as shown for the `ticker` component:

```
GOPATH=`pwd` go build -o components/core/ticker github.com/cascades-fbp/cascades/components/core/ticker
```

## Configuration

The only configuration you need to perform is creating Component Library file, which will be used by the Runtime. A `cascades` executable has a command `library` with subcommand `add` to manage your local library of available components. Run the following command to create the registry:

```
$ ./cascades library add ./components/
Walking components directory: /path/to/my/installation/components
Added "bonjour/discover"
Added "bonjour/register"
Added "core/console"
Added "core/delay"
Added "core/drop"
Added "core/exec"
Added "core/joiner"
Added "core/passthru"
Added "core/readfile"
Added "core/splitter"
Added "core/submatch"
Added "core/switch"
Added "core/template"
Added "core/ticker"
Added "debug/crasher"
Added "debug/oneshot"
Added "fs/walk"
Added "fs/watchdog"
Added "influxdb/write"
Added "mqtt/pub"
Added "mqtt/sub"
Added "sockets/tcp-server"
Added "websocket/client"
Added "websocket/server"
DONE
```

The exact output may very from release to release and depends on the additional components you put into your `components` directory. As a result of executing this command a file `library.json` will be created in the current directory.

## Running examples

From now on you can run the flows which rely on components registered in your local library. The `examples` directory has some of them. 

For instance, a simple ticking flow, which generates s timestamp every given number of seconds, passes the timestamp through a dummy filter, which sends data to a component that writes all received information to the _stdout_.

```
$ cascades run examples/tick.fbp
runtime | Starting processes...
runtime | Activating processes by sending IIPs...
runtime | Sending '5s' to socket 'tcp://127.0.0.1:5000'
Log     | IP: 1416337759
Log     | IP: 1416337764
...
<Ctrl+C pressed>
...
^Cruntime | Shutdown...
runtime | sending SIGTERM to Ticker
runtime | sending SIGTERM to Forward
runtime | sending SIGTERM to Log
Log     | Stopped
Forward | Stopped
Ticker  | Stopped
runtime | Shutdown...
Stopped
```

## What's next?

Now you have a custom build of Cascades framework. The next step is to proceed to [Getting Started](Getting-started) guide and follow the tutorial.