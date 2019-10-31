# TCP/UDP socket components for Cascades FBP

This repository contains the following components:

 * TCP server
 * TBD...


## TCP Server component

Usage of the server component:

```
$ ./components/socket/tcp-server
Usage of ./components/socket/tcp-server:
  -debug=false: Enable debug mode
  -json=false: Print component documentation in JSON
  -port.in="": Component's input port endpoint
  -port.options="": Component's options port endpoint
  -port.out="": Component's output port endpoint
```

Example (TCP echo server):

```
# Configure server and forward its output to packet cloning component
'localhost:9999' -> OPTIONS Server(sockets/tcp-server) OUT -> IN Copy(core/splitter)

# Show incoming data on the screen
Copy OUT[0] -> IN Log(core/console)

# Feedback to respond a user with what was received
Copy OUT[1] -> IN Server
```

Running example:

```
$ bin/cascades run examples/tcp-server.fbp
runtime | Starting processes...
runtime | Activating processes by sending IIPs...
runtime | Sending 'localhost:9999' to socket 'tcp://127.0.0.1:5000'
```

Connect to the server and type `Hello world`:

```
± |master ✗| → telnet localhost 9999
Trying ::1...
telnet: connect to address ::1: Connection refused
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
Hello world
```

You should see something like this on the screen (and receive `Hello world` string back in the telnet window):

```
Log     | [
Log     | IP: b729e449-a093-4e37-6d55-7550bdb86efa
Log     | IP: Hello world
Log     | ]
```
