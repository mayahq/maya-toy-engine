
## Picking the right distribution

Installing Cascades from the binary distribution is probably the fastest way to get started. Please, visit the [Github Releases Page](https://github.com/cascades-fbp/cascades/releases) and select the latest stable release of the framework. Find there the .zip or a .tar.gz distribution that correspond to the platform you are planning to use Cascades on. The distribution name includes the suffix _linux_, _darwin_ or _windows_, and an architecture type _amd64_, _arm_ or _i386_. If you don't see the distribution for your platform let us know.

## Installation

Download the distribution and unarchive it into the prefered location on your target machine. The distribution folder will have the following structure:

```
cascades OR cascades.exe
components/
examples/
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

You can start exploring the available components from the distribution and try to execute other examples or, which is highly recommendable, proceed to [Getting Started](Getting-started) guide and follow the tutorial.
