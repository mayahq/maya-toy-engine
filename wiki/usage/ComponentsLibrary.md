The _Components Library_ is a registry of _Components_ that helps the _Registry_ / _Scheduler_ to parse the flow definition and initiate all required processes for execution. The Library might be also used by visual programming tool to assist a network designer during the flow composition. 

There can be several component libraries, each created for a specific problem domain. Developers are able to create own libraries from the available components and split/extend/share existing ones.

Another problem a Library should help solving is an organization of large components sets and reducing misuse of components by a network designer.

## Location of a library

By default `cascades` command is looking for a library file named `library.json` in the current working directory. You can specify a library file explicitly using `--file` (or `-f`) flag. Examples:

```
$ cascades --file=/path/to/library library add demo
```
or
```
$ cascades --file=/path/to/library --debug run /path/to/flow.json
```

## Adding a folder with components

Adding a folder with components recursively is simple:

```
$ cascades --file=/path/to/library library add /path/to/components/folder/
```

## Adding a single component

Adding a single component requires specifying a registration name, which will be used in the flows. For instance, after registering a component under `myprefix/componentname` name in this way:

```
$ cascades --file=/path/to/library library add --name="myprefix/componentname" /path/to/components/folder/
```

you can use it in the FBP DSL program as shown below:

```
'some data' -> PORT MyProcess(myprefix/componentname)
```

## Adding a composite component (subnet)

The `library add` command can register subnets in the same way as you would register a single component or a folder of components: simply provide a path to a .json or .fbp file.
