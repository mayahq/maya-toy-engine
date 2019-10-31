## Prerequisites

The communication between components is implemented using ZeroMQ, which needs to be installed on your machine as described in the [ZeroMQ Python Bindings documentation](http://zeromq.org/bindings:python).

## Step by step tutorial

In this example we will examine how to create a simple _Passthru_ component that can be used later for implementing different filters/transformations of the IP streams.

### Step 1: Define the CLI argument

The recommended way is to use the native Python's module `argparse` to define the flags you can pass to your component. Here's an example for the _Passthru_ component:

```
import argparse
parser = argparse.ArgumentParser(
   description='Cascades component for passing IPs through and logging them to stdout')
parser.add_argument(
   "-debug", "--debug", help="Enable debug mode", action="store_true")
parser.add_argument(
   "-json", "--json", help="Print component registration data in JSON", action="store_true")
parser.add_argument(
   "-port.in", "--port.in", type=str, help="Component input port endpoint")
parser.add_argument(
   "-port.out", "--port.out", type=str, help="Component out port endpoint")
args = parser.parse_args()
```

The `-debug` and `-json` are standard flags used by the _Runtime_. The other flags are used to configure the ports of the component and follow the naming convention `port.<name>`.

### Step 2: Implement JSON registration

Your component needs to be added to _Component Library_ first, before the runtime is able to execute it. For this you need to implement the proper output when your script is executed with `-json` flag:

```
import json

if args.json:
   info = dict(
       description="Forwards received IP to the output without any modifications",
       elementary=True,
       inports=(dict(name="IN", type="all", description="Input port for receiving IPs",
                     required=True, addressable=False),),
       outports=(dict(name="OUT", type="all",
                      description="Output port for sending IPs",
                      required=True, addressable=False),),
   )
   print json.dumps(info)
   sys.exit(0)
```

### Step 3: Create input port

The actual endpoint for the port is passed via CLI flag. In our case the _Runtime_ uses `-port.in` flag to pass a TCP endpoint (e.g. `tcp://127.0.0.1:5000`). You extract this value from arguments, create a ZeroMQ `Context` object, which is used to created a `Socket`. Input ports are sockets of type `zmq.PULL` and you should call `bind` method of them (these ports are passive and wait for connections).

```
...
args = parser.parse_args()
data = vars(args)
in_addr = data.get('port.in')

ctx = zmq.Context()

in_port = ctx.socket(zmq.PULL)
in_port.bind(in_addr)

```

### Step 4: Create output port

The output port is created similar to the input port but with different type `zmq.PUSH` as it used to send (push) data to it. You also need to call `connect` method of this port instead of `bind`.

```
out_addr = data.get('port.out')

out_port = ctx.socket(zmq.PUSH)
out_port.connect(out_addr)
```

### Step 5: Component main loop

Now we are ready to start the main loop of the component, which is responsible for reading from IN port(s), optionally processing the data and sending results to the OUT port(s). For a simple _Passthru_ component we just forward the data from input to output ports. In Cascades the data is sent as multipart ZeroMQ message therefore we need to use corresponding `recv_multipart` and `send_multipart` methods of the `Socket` object.

```
while True:
   if args.debug:
	   print "Waiting for IP..."
   
   parts = in_port.recv_multipart()
   
   if args.debug:
	   print "IN Port: Received: %r" % parts
   
   out_port.send_multipart(parts)
   
	if args.debug:
		print "OUT Port: Forwarded"
```

## Complete example of a Passthru component

```
#!/usr/bin/env python

import argparse
import zmq
import sys
import json


def main():
    parser = argparse.ArgumentParser(
        description='Cascades component for passing IPs through and logging them to stdout')
    parser.add_argument(
        "-debug", "--debug", help="Enable debug mode", action="store_true")
    parser.add_argument(
        "-json", "--json", help="Print component registration data in JSON", action="store_true")
    parser.add_argument(
        "-port.in", "--port.in", type=str, help="Component input port endpoint")
    parser.add_argument(
        "-port.out", "--port.out", type=str, help="Component out port endpoint")
    args = parser.parse_args()

    if args.json:
        info = dict(
            description="Forwards received IP to the output without any modifications",
            elementary=True,
            inports=(dict(name="IN", type="all", description="Input port for receiving IPs",
                          required=True, addressable=False),),
            outports=(dict(name="OUT", type="all",
                           description="Output port for sending IPs",
                           required=True, addressable=False),),
        )
        print json.dumps(info)
        sys.exit(0)

    data = vars(args)
    in_addr = data.get('port.in')
    out_addr = data.get('port.out')
    if not in_addr or not out_addr:
        parser.print_help()
        sys.exit(1)

    ctx = zmq.Context()

    in_port = ctx.socket(zmq.PULL)
    in_port.bind(in_addr)

    out_port = ctx.socket(zmq.PUSH)
    out_port.connect(out_addr)

    while True:
        if args.debug:
            sys.stdout.write("Waiting for IP...\n")
            sys.stdout.flush()

        parts = in_port.recv_multipart()

        if args.debug:
            sys.stdout.write("IN Port: Received: %r" % parts)
            sys.stdout.flush()

        out_port.send_multipart(parts)

        if args.debug:
            sys.stdout.write("OUT Port: Forwarded!\n")
            sys.stdout.flush()

if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        sys.exit(0)

```