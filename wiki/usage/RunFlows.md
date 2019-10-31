After you have updated your _Components Library_ as described in [Manage components library](Components-library) documentation and written your first program you can execute the program using Cascades runtime's `run` command.

### Running options

You can execute `cascades run --help` command to see all possible running options as shown below:

```
NAME:
   run - Runs a given graph defined in the .fbp or .json formats

USAGE:
   command run [command options] [arguments...]

OPTIONS:
   --port, -p '5000'	initial port to use for connections between nodes
   --dry		dry run (parses and validates the graph, exits without executing it)
```

Additionally you can specify the global command line flags:

```
GLOBAL OPTIONS:
   --file, -f 'library.json'	components library file
   --debug, -d			enable extra output for debug purposes
```

### Validating a flow

In order to validate the flow without actually executing it you can run the following command (notice `--debug` and `--dry` flags):
```
$ cascades --debug run --dry examples/tick.fbp
--------- Properties ----------
---------- Inports -----------
---------- Outports -----------
---------- Processes ----------
Ticker (core/ticker)
Forward (core/passthru)
Log (core/console)
--------- Connections ---------
'5s' -> INTERVAL Ticker
Ticker OUT -> IN Forward
Forward OUT -> IN Log
-------------------------------
```
This indicates that a flow is valid.

If we try to validate a flow that uses a component not registered in Library an error will be shown:

```
$ cascades --debug run --dry examples/bad-tick.fbp
Failed to load/flatten graph: Component core/newticker not found in the library
```

Or if you have a syntax error in FBL DSL a corresponding error will be shown:

```
$ cascades --debug run --dry examples/tick.fbp
Failed to load/flatten graph: Failed to parse graph definition:
...
```

### Running with/without debug option

Executing a flow with a `--debug` flag usually produces a lot of additional output in the _stdout_ and it can give you a lot of insights on how the Cascades Runtime is manages the processes.

#### Example executing flow without --debug options

```
$ cascades run examples/tick.fbp
runtime | Starting processes...
runtime | Activating processes by sending IIPs...
runtime | Sending '5s' to socket 'tcp://127.0.0.1:5000'
Log     | IP: 1416347159
Log     | IP: 1416347164

...

^Cruntime | Shutdown...
runtime | sending SIGTERM to Forward
runtime | sending SIGTERM to Log
runtime | sending SIGTERM to Ticker
Forward | Stopped
Ticker  | Stopped
Log     | Stopped
runtime | Shutdown...
Stopped
```

#### Example executing flow with --debug options

```
$ cascades --debug run examples/tick.fbp
--------- Properties ----------
---------- Inports -----------
---------- Outports -----------
---------- Processes ----------
Ticker (core/ticker)
Forward (core/passthru)
Log (core/console)
--------- Connections ---------
'5s' -> INTERVAL Ticker
Ticker OUT -> IN Forward
Forward OUT -> IN Log
-------------------------------
Log.IN [tcp://127.0.0.1:5002]
Ticker.INTERVAL [tcp://127.0.0.1:5000]
Ticker.OUT [tcp://127.0.0.1:5001]
Forward.IN [tcp://127.0.0.1:5001]
Forward.OUT [tcp://127.0.0.1:5002]
--------- Executables ---------
Forward: "/Users/alex/Projects/OpenSource/Cascades/components/core/passthru" "-port.in=\"tcp://127.0.0.1:5001\" -port.out=\"tcp://127.0.0.1:5002\" -debug=\"true\" "
Log: "/Users/alex/Projects/OpenSource/Cascades/components/core/console" "-debug=\"true\" -port.in=\"tcp://127.0.0.1:5002\" "
Ticker: "/Users/alex/Projects/OpenSource/Cascades/components/core/ticker" "-debug=\"true\" -port.interval=\"tcp://127.0.0.1:5000\" -port.out=\"tcp://127.0.0.1:5001\" "
-------------------------------
runtime | Starting processes...
Ticker  | Wait for configuration IP...
Log     | Started...
Forward | Started...
runtime | Activating processes by sending IIPs...
runtime | Sending '5s' to socket 'tcp://127.0.0.1:5000'
Ticker  | Started...
Ticker  | Configured to tick with interval: 5s
Ticker  | 1416347216
Log     | IP: 1416347216
Ticker  | 1416347221
Log     | IP: 1416347221

...

^Cruntime | Shutdown...
runtime | sending SIGTERM to Log
runtime | sending SIGTERM to Ticker
runtime | sending SIGTERM to Forward
Log     | Error receiving message: interrupted system call
Ticker  | Give 0MQ time to deliver before stopping...
Log     | Give 0MQ time to deliver before stopping...
Forward | Give 0MQ time to deliver before stopping...
Forward | Error receiving message: interrupted system call
Forward | Stopped
Ticker  | Stopped
Log     | Stopped
Forward | Stopped
Ticker  | Stopped
Log     | Stopped
runtime | Shutdown...
Stopped
```