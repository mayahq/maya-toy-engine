Cascades support any programming language, which is supported by ZeroMQ. That includes C, C++, C#, CL, Delphi, Erlang, F#, Felix, Haskell, Java, Objective-C, PHP, Python, Lua, Ruby, Ada, Basic, Clojure, Go, Haxe, Node.js, ooc, Perl, and Scala. But not limited to them only. The complete list can be found at this page [ZeroMQ bindings](http://zeromq.org/bindings:_start). 

The only requirements to a component are:

 * Executable should support `--debug` flag (e.g. `--debug=true` should enable output of extra debugging information into _stdout_)
 * Executable should support `--json` flag (e.g. `--json` should output component registration information and exit with status code 0)
 * The ports endpoints should be read from command line argument named using pattern `port.<name>` (e.g. `--port.in`, `--port.options`, `--port.out`).