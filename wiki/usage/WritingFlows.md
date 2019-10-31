The flows for Cascades can be written using one of the coordination language formats: 

 * [FBP DSL](http://noflojs.org/documentation/fbp/)
 * NoFlo's [JSON](http://noflojs.org/documentation/json/) schema

## Writing FBP DSL

### Basics

The syntax of the FBP DSL is very simple. Consider the following example:

```
'some data' -> PORT ProcessName(ComponentName)
```

Here _some data_ or in terms of FBP _Initial Information Packet_ (IIP) is sent to a process named _ProcessName_ based on component _ComponentName_ using its input port called _PORT_. We `->` to define a connection between the building blocks of the flow.

```
ProcessA(ComponentA) OUT -> IN ProcessB(ComponentB) OUT -> IN ProcessC(ComponentC)`
```

In the example above we have defined two connections in a single line. Once _ProcessA_, _ProcessB_ and _ProcessC_ are defined we can use them in other parts of the flow:

```
ProcessA ERR -> IN Log(Console)
ProcessB ERR -> Log
ProcessC ERR -> Log
```

As you can see from above we use only process name to refer to the same process defined above.

### Array ports

Array ports are named ports that have indices. Each index is a separate port and can be connected independently. This feature allows to have a dynamic number of slots in a component without chancing its source code and re-compilation. An array port can be used as an input and as an output port of a component.
```
Ticker1(Timer) OUT -> IN[0] Mux(Merge) OUT -> IN Log(Console)
Ticker2(Timer) OUT -> IN[1] Mux
Ticker3(Timer) OUT -> IN[2] Mux
```
In the example above we connected three timers into an array port of the _Mux_ process. 

### Exporting ports

One of the FBP features besides the power of configurable modularity is a step-wise decomposition using so called _subnets_ or _subgraphs_. We can design a flow and make it re-usable in the future, become a part of a large flow. For this purpose we need a way to export input and output ports so the flow can be used as any other _black box_ component.

Consider creating a file _subticker.fbp_ with the following definition:

```
INPORT=Ticker.INTERVAL:OPTIONS
OUTPORT=Forward.OUT:TICKS

Ticker(Ticker) OUT -> IN Forward(Passthru)
```

Here we are exporting the input port _INTERVAL_ of the process _Ticker_ under a new name _OPTIONS_. The output of the _Forward_ process via its _OUT_ component is exported under _TICKS_ name. Then you can use this flow (after registering it in the _Components Library_ as described in [Manage components library](Components-library)) in other flows:

```
'15s' -> OPTIONS NewTicker(subticker)
NewTicker TICKS -> IN Log(Console)
```

Without looking into a source code of each component we cannot see if it is an elementary component or a composite one (so called _subnet_).

## Writing NoFlo's JSON

The JSON representation of the flow is documented in the [NoFlo's JSON Documentation](http://noflojs.org/documentation/json/). Writing JSON by hand is not the best programming experience. We have added JSON support to Cascades only to be able to use third-party visual tools such as [NoFlo UI](http://noflojs.org/noflo-ui/).

### Example

A typical ticker example flow written using NoFlo's JSON:

```
{
  "properties": {
    "name": "Send tick via passthrough to the console"
  },
  "processes": {
    "Time Ticker": {
      "component": "core/ticker"
    },
    "IP Forwarder": {
      "component": "core/passthru"
    },
    "IP Logger": {
      "component": "core/console"
    }
  },
  "connections": [
    {
      "data": "5s",
      "tgt": {
        "process": "Time Ticker",
        "port": "INTERVAL"
      }
    },
    {
      "src": {
        "process": "Time Ticker",
        "port": "OUT"
      },
      "tgt": {
        "process": "IP Forwarder",
        "port": "IN"
      }
    },
    {
      "src": {
        "process": "IP Forwarder",
        "port": "OUT"
      },
      "tgt": {
        "process": "IP Logger",
        "port": "IN"
      }
    }
  ]
}

```