If you are new to programming ZeroMQ with Node.js the following resources are recommended for review:

 * [http://www.slideshare.net/fedario/zero-mq-with-nodejs](http://www.slideshare.net/fedario/zero-mq-with-nodejs)
 * [http://naviger.com/?p=39](http://naviger.com/?p=39)

## Prerequisites

The communication between components is implemented using ZeroMQ, which needs to be installed on your machine as described in the [ZeroMQ Node.JS Bindings documentation](http://zeromq.org/bindings:node-js).

## Step by step tutorial

In this example we will examine how to create a simple _Passthru_ component that can be used later for implementing different filters/transformations of the IP streams.

### Step 1: Define the CLI argument

In our example we use `nomnom` package for simplifying our work with CLI arguments. You can install this package simply using `npm install nomnom`.

```
var opts = require("nomnom")
   .option('debug', {
      flag: true,
      default: false,
      help: 'Enable debug mode'
   })
   .option('json', {
      flag: true,
      help: 'Print component registration data in JSON'
   })
   .option('port.in', {
      help: "Component input port endpoint"
   })
   .option('port.out', {
      help: "Component out port endpoint"
   })
   .parse();
```

The `--debug` and `--json` are standard flags used by the _Runtime_. The other flags are used to configure the ports of the component and follow the naming convention `port.<name>`.

### Step 2: Implement JSON registration

Your component needs to be added to _Component Library_ first, before the runtime is able to execute it. For this you need to implement the proper output when your script is executed with `--json` flag:

```
if (opts.json) {
  console.log(JSON.stringify({
    "description":"Forwards received IP to the output without any modifications",
    "elementary":true,
    "inports":[
      {"name":"IN","type":"all","description":"Input port for receiving IPs","required":true,"addressable":false}
    ],
    "outports":[
      {"name":"OUT","type":"all","description":"Output port for sending IPs","required":true,"addressable":false}
    ]
  }));
  process.exit(0);
}
```

### Step 3: Create input & output ports

The actual endpoint for the port is passed via CLI flag. In our case the _Runtime_ uses `-port.in` flag to pass a TCP endpoint (e.g. `tcp://127.0.0.1:5000`). Input ports are sockets of type `pull` and you should call their`bindSync` method (these ports are passive and wait for connections).
The output port is created similar to the input port but with different type `push` as it used to send (push) data to it. You also need to call `connect` method of this port instead of `bindSync`.

```
var zmq  = require('zmq')
  , receiver = zmq.socket('pull')
  , sender = zmq.socket('push');
 
...

receiver.bindSync(opts['port.in']);
sender.connect(opts['port.out']);
```

### Step 5: Component main loop

Now we are ready to start the main loop of the component, which is responsible for reading from IN port(s), optionally processing the data and sending results to the OUT port(s). For a simple _Passthru_ component we just forward the data from input to output ports. In Cascades the data is sent as multipart ZeroMQ message therefore we need to read all function arguments passed to the handler.

```
receiver.on('message', function() {
  // Note that separate message parts come as function arguments.
  var args = Array.apply(null, arguments);
  process.stdout.write("received: " + args[1].toString());

  // Pass array of strings/buffers to send multipart messages.
  sender.send(args);
  process.stdout.write("done!")
});
```

## Complete example of a Passthru component

```
#!/usr/bin/env node

// Parse CLI arguments
var opts = require("nomnom")
   .option('debug', {
      flag: true,
      default: false,
      help: 'Enable debug mode'
   })
   .option('json', {
      flag: true,
      help: 'Print component registration data in JSON'
   })
   .option('port.in', {
      help: "Component input port endpoint"
   })
   .option('port.out', {
      help: "Component out port endpoint"
   })
   .parse();

// Print JSON info if requested
if (opts.json) {
  console.log(JSON.stringify({
    "description":"Forwards received IP to the output without any modifications",
    "elementary":true,
    "inports":[
      {"name":"IN","type":"all","description":"Input port for receiving IPs","required":true,"addressable":false}
    ],
    "outports":[
      {"name":"OUT","type":"all","description":"Output port for sending IPs","required":true,"addressable":false}
    ]
  }));
  process.exit(0);
}

// Validate port endpoints
if (opts['port.in'] == undefined || opts['port.out'] == undefined) {
  console.log("ERROR: no port endpoints provided!");
  process.exit(1);
}

// Create zmq sockets & run the main loop
var zmq  = require('zmq')
  , receiver = zmq.socket('pull')
  , sender = zmq.socket('push');

receiver.on('message', function() {
  // Note that separate message parts come as function arguments.
  var args = Array.apply(null, arguments);
  process.stdout.write("received: " + args[1].toString());

  // Pass array of strings/buffers to send multipart messages.
  sender.send(args);
  process.stdout.write("done!")
});

receiver.bindSync(opts['port.in']);
sender.connect(opts['port.out']);

process.on('SIGINT', function () {
  receiver.close();
  sender.close();
  process.exit(0);
});

```