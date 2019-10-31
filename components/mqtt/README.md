# MQTT components for Cascades FBP

## Usage

Usage of MQTT's publication component:

```
$ ./components/mqtt/pub
Usage of ./components/mqtt/pub:
  -debug=false: Enable debug mode
  -json=false: Print component documentation in JSON
  -port.err="": Component's output port endpoint
  -port.in="": Component's input port endpoint
  -port.options="": Component's input port endpoint
```

Usage of MQTT's subscription component:

```
$ ./components/mqtt/sub
Usage of ./components/mqtt/sub:
  -debug=false: Enable debug mode
  -json=false: Print component documentation in JSON
  -port.err="": Component's output port endpoint
  -port.options="": Component's input port endpoint
  -port.out="": Component's output port endpoint
```

## Examples

Publishing data from network to MQTT:

```
# Configure publisher with broker on localhost and default topic 'demo/1'
'tcp://127.0.0.1:1883/demo/1' -> OPTIONS MosquittoPublisher(mqtt/pub)

# Output all errors to console
MosquittoPublisher ERR -> IN Errors(core/console)

# Start a timer to send current timestamp every 2s to publisher
'2s' -> INTERVAL Ticker(core/ticker)
Ticker OUT -> IN MosquittoPublisher
```

Subscribing to data from MQTT:

```
# Configure subscriber. Pay attention to '%23' which is a '#'
'tcp://127.0.0.1:1883/%23' -> OPTIONS MosquittoSubscriber(mqtt/sub)

# Output all errors to console
MosquittoSubscriber ERR -> IN Errors(core/console)

# Output incoming substreams form MQTT into console as well
MosquittoSubscriber OUT -> IN Log(core/console)
```

Example out the output from runnning subscription example:

```
$ bin/cascades --debug run examples/mqtt-sub.fbp
MosquittoSubscriber | Started...
Log                 | [
Log                 | IP: /demo/1
Log                 | IP: 1412254327
Log                 | ]
Log                 | [
Log                 | IP: /demo/1
Log                 | IP: 1412254329
Log                 | ]
Log                 | [
Log                 | IP: /demo/1
Log                 | IP: 1412254331
Log                 | ]
```