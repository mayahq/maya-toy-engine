# Websocket components for Cascades FBP

This repository contains the following components:

 * Websocket Server
 * Websocket Client

## Websocket Server Component

Usage of Websocket server component:

```
$ ./components/websocket/server
Usage of ./components/websocket/server:
  -debug=false: Enable debug mode
  -json=false: Print component documentation in JSON
  -port.in="": Component's input port endpoint
  -port.options="": Component's options port endpoint
  -port.out="": Component's output port endpoint
```

Example (Websocket echo server):

```
# Configure server and forward its output to packet cloning component
'127.0.0.1:8000' -> OPTIONS Server(websocket/server)
Server OUT -> IN Demux(core/splitter)

# Show incoming data on the screen
Demux OUT[0] -> IN Log(core/console)

# Return receive data back to server's client
Demux OUT[1] -> IN Server
```

Running example:

```
$ bin/cascades run examples/tcp-server.fbp
runtime | Starting processes...
runtime | Activating processes by sending IIPs...
runtime | Sending 'localhost:9999' to socket 'tcp://127.0.0.1:5000'
```

Open http://www.websocket.org/echo.html and type in the following location `ws://localhost:9999/ws`. Then enter the following message `{"msg":"Hi!"}` and send it. You should receive exactly the same response back.

## Websocket Client Component

Usage of Websocket client component:

```
$ ./components/websocket/client
Usage of ./components/websocket/client:
  -debug=false: Enable debug mode
  -json=false: Print component documentation in JSON
  -port.in="": Component's input port endpoint
  -port.options="": Component's options port endpoint
  -port.out="": Component's output port endpoint
```

Example (Connecting to public websocket echo server):

```
# Configure client to connect to a public echo server
'ws://echo.websocket.org:80' -> OPTIONS Client(websocket/client)

# Output all received data on the screen
Client OUT -> IN Log(core/console)

# Configure ticker (sends timestamp every 3 seconds)
'3s' -> INTERVAL Ticker(core/ticker)

# Send ticker's timestamp to the client
Ticker OUT -> IN Client
```

Running example:

```
...
Ticker  | 1412264587
Client  | Sending data from IN port to websocket...
Client  | Writer: will send to websocket [49 52 49 50 50 54 52 53 56 55]
Client  | Reader: received from websocket [49 52 49 50 50 54 52 53 56 55]
Client  | Sending data from websocket to OUT port...
Log     | IP: 1412264587
...
```

As you can see above a tick `1412264587` is sent to Client, which sends it to Websocket server. The echo server responds back with the same message, which is then forwarded to Log component and shown in the terminal output.
