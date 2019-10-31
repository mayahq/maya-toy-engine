## Overview

Similar to a classical FBP system Cascades consists of two main subsystems:

 * _Runtime_ or _Scheduler_
 * _Components Library_

A _Runtime_ is responsible for executing a provided program. The implementation is based partially  on the ideas and source code of the Go port of the [Foreman](https://github.com/ddollar/foreman) (Procfile-based process manager) but it goes beyond just process management and also includes program definition parsing, validation, command-line arguments generation, library management, etc.

The _Components Library_ is a registry of _Components_, which, as described in the [Flowbased Wiki](https://github.com/flowbased/flowbased.org/wiki/Concepts), are the building blocks of the FBP program. They are usually written as classes, functions or small programs in conventional programming languages. 

Cascades provides some core components for typical data transformation, such as: _Ticker_, _Splitter_, _Joiner_, _ReadFile_, _Submatch_, _Switch_, _Template_, _WalkDir_, etc. Additional components can be either downloaded and registered in the _Library_ or implemented for specific project needs.

## Trade-offs of being language agnostic

In classical FBP a process is a _black box_ component that has its control thread and environment. Depending on the implementation a FBP process can be realized as a [green thread](http://en.wikipedia.org/wiki/Green_threads), a native thread or a system level process. If we want to support multiple implementation technologies and avoid build/compilation step for every network before execution the only option left for us is running FBP process as a _system level process_. This makes the execution of a network to be much heavier in comparison to JavaFBP's native threads based execution and NoFlo's single threaded execution. From this perspective executing a relatively large network (300 notes or more) has a different CPU/memory footprint in JavaFBP, NoFlo and Cascades. Therefore a granularity of the components should be selected carefully.

While the programs implemented with the same technology might have certain ways of interprocess communication, how to make two programs executed as separate processes, written in different programming languages and executed on either same of different computers exchange the data? You can use files (shared file system) or, which is the most convenient, networking? Using raw sockets, a database or a message broker (RabbitMQ or Mosquitto) is a good option for inter-process communication. In the classical FBP systems connections between the nodes in the network are implemented as bounded buffers (FIFO queues) with a defined capacity. From this point of view we consider using either broker-based or broker-less solution for connections between nodes. The usage of a broker brings additional prerequisites to the installation and execution, while the broker-less approach simplifies the end-user experience.

Based on these considerations we have decided to use [ZeroMQ](http://zeromq.org) for implementing connections between the processes:

 * Cross-platform
 * Supports message exchange via inproc, IPC, TCP, TPIC, multicast
 * Implements various communication patterns (although for FBP implementation PULL/PUSH sockets pair is sufficient)
 * Broker-less approach (therefore _zero_ in the name) simplifies the installation and developer experience
 * Bindings to variety of programming languages (C, C++, C#, CL, Delphi, Erlang, F#, Felix, Haskell, Java, Objective-C, PHP, Python, Lua, Ruby, Ada, Basic, Clojure, Go, Haxe, Node.js, ooc, Perl, and Scala) allows writing components in the language of your choice

## Coordination language

The coordination language is a language for describing the flows or networks, which are programs in FBP and in Cascades in particular. Although the best possible flows composition can be done using [Visual Programming](http://en.wikipedia.org/wiki/Visual_programming_language) by direct manipulation of data and its transformation we still consider writing text as an efficient and the most popular way of programming (on the other hand the development of proper visual programming tools might change this trend).

In language specific FBP systems there is always a way to write a flow using a conventional language the runtime is implemented with. For instance, in JavaFBP you can defined network with Java. In NoFlo.js a flow can be written either in CoffeeScript or directly in JavaScript. In a language agnostic runtime of Cascades we were faced by a choice of picking a one language to be a _coordination_ one, supporting any conventional language, developing our own Domain Specific Language (DSL) for dataflow or re-using any existing DLS. As a result we have discovered a [FBP DSL](http://noflojs.org/documentation/fbp/), which seemed to be simple and powerful enough to describe flows by a human in any text editor. Additionally we have added support for [JSON](http://noflojs.org/documentation/json/) format defined by NoFlo, which would allow Cascades to execute flows created using the soon-to-be-developed visual programming tools.

Example of FBP DSL:

```
# This is a simple flow to send ticks every 5 seconds
'5s' -> INTERVAL Ticker(core/ticker) OUT -> IN Forward(core/passthru)
Forward OUT -> IN Log(core/console)
```

Example of the same flow written in JSON DSL:

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

As you probably noticed JSON DSL still can be written/edited by humans but is more verbose than FBP DLS.

## Component design

A component in Cascades is a standalone executable program for a particular platform or cross-platform. The _black box_ concept is implemented in a component in the following way. A component has zero knowledge about its executing context (runtime, connections to other components) and communicates with an outside world only through the ports. The ports in a component are implemented using ZeroMQ sockets. 

But not every executable using ZeroMQ is a Cascades component. There is a minimal set of conventions to be implemented by a component to make it Cascades-compliant. The compliance is achieved by supported a set of command line arguments as described below:

 * Flag `help` should print a detailed help on a component in a human readable format
 * Flag `debug` is passed by a Runtime running in a debug mode. You can use this flag to increase the log verbosity level (writing into _stdout_ or _stderr_ streams) or performing other extra logic that would help to trace the component execution.
 * Flag `json` is used by a Library to extra meta information from a component and create a proper registration entry in the registry. Only registered components are resolved by a Runtime during execution. You can read more on registration format on the [Components library](Components-library) page.
 * The endpoints for configuring ZeroMQ sockets are provided by a runtime using the flags, which are named according to the scheme `port.*` (e.g. `port.in`, `port.options`, `port.err`)

An example of the `./my-component -help` output:
```
Usage of my-component:
  -debug=false: Enable debug mode
  -json=false: Print component documentation in JSON
  -port.interval="": Component's input port endpoint for configuring interval
  -port.out="": Component's output port endpoint to send timestamps to
```

## Component Library

The _Components Library_ is a registry of _Components_, which helps the _Registry_ / _Scheduler_ to parse the flow definition and initiate all required processes for execution. The Library might be also used by visual programming tool to assist a network designer during the flow composition. 

There can be several component libraries, each created for a specific problem domain. Developers are able to create own libraries from the available components and split/extend/share existing ones.

Another problem a Library should help solving is an organization of large components sets and reducing misuse of components by a network designer.

## Runtime / Scheduler

A _Runtime_ is responsible for executing a provided program. The implementation is based partially  on the ideas and source code of the Go port of the [Foreman](https://github.com/ddollar/foreman) (Procfile-based process manager) but it goes beyond just process management and also includes program definition parsing, validation, command-line arguments generation, library management, etc.
