[![Stories in Ready](https://badge.waffle.io/cascades-fbp/cascades.png?label=ready&title=Ready)](https://waffle.io/cascades-fbp/cascades)

## Maya Flow-Based Toy Robot Engine

This language engine has been adapted from the language-agnostic flow-based language Cascades in order to create a toy sandbox to experiment with dynamic generation of DAGs (Directed Acyclic Graphs). It has a limited dictionary of components that can be wired together to do various complicated tasks. 

### Cascades - Language-Agnostic Programming Framework for Data-Driven Applications

Cascades is language-agnostic [dataflow](http://en.wikipedia.org/wiki/Dataflow_programming) programming framework that implements a number of the [Flow-Based Programming](http://en.wikipedia.org/wiki/Flow-based_programming) (FBP) concepts. Although we consider Cascades to be a general-purpose programming framework we have created it having a, as J. P. Morrison described: a *data factory mental image in mind, where the application is expressed as a series of transforms on data streams - which requires fundamental changes from the old von Neumann thinking in the way programmers build applications*.

### Features of Cascades

 * Cascades is cross-platform (can be used on OSX/Linux/Windows, tested also on Raspberry PI and Beagleboard Black) and completely written in [Go](http://golang.org/) programming language. 
 * Uses [ZeroMQ](http://zeromq.org) for connections between the components.
 * The core components are also written in Go, but you are free to choose any programming language of your choice as long as there are bindings for it (currently it supports: C, C++, C#, CL, Delphi, Erlang, F#, Felix, Haskell, Java, Objective-C, PHP, Python, Lua, Ruby, Ada, Basic, Clojure, Go, Haxe, Node.js, ooc, Perl, and Scala)
 * Supports flows defined using NoFlo's [JSON](http://noflojs.org/documentation/json/) format or FBP [DSL](http://noflojs.org/documentation/fbp/) language.

Visit our [Wiki](https://github.com/sibeshkar/maya-toy-engine/wiki) for details on how to get started.

## Usage

General usage is described below:

```
NAME:
   cascades - A Cascades FBP runtime/scheduler for the FBP applications.

USAGE:
   cascades [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   run      Runs a given graph defined in the .fbp or .json formats
   library  Manages a library of components
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --file, -f 'library.json'  components library file
   --debug, -d       enable extra output for debug purposes
   --help, -h        show help
   --version, -v     print the version
```

## Get Dependencies 
```
go get github.com/sibeshkar/maya-toy-engine/cmd/cascades

```

## Building Binary

Building the binary: 
```
go build ./cmd/cascades

```

## Adding new components

Building new components
```
go build -o components_built/maya/click ./components/maya/click
./cascades library add ./components_built/
```
Build in bulk : 
```
./build.sh
```

Add components to library, will be added to library.json : 
```
./cascades library add ./components_built/
OR 
./addToLibrary.sh
```

## Run FBP file

To test that everything is working as expected, run this test.

```
./cascades --debug run ./examples/maya_test.fbp
```

## License

The MIT License

Copyright &copy; 2020

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
