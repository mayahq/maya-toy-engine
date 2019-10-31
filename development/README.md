Cascades FBP Development Environment Layout
===========

Scripts &amp; configuration required for setting up a Cascades development environment 

## Setup

```
$ ./setup.sh
This script configures current directory as Cascades FBP development environment.
Proceed? [y/n] y
...
DONE!
```

## Re-build cascade CLI and components

```
$ ./build
Building core components:
 -> core/console... ✔
 -> core/delay... ✔
 -> core/drop... ✔
 ...
```

## Release

```
$ release.sh v0.1.0
...
DONE!
```