#!/bin/bash -e

source func.sh

###########################################################################
# Main
#

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

mkdir -p components

# Build cascades CLI
echo -n "Building runtime ... "
go get github.com/cascades-fbp/cascades/cmd/cascades
echo_pass

# Build components
build_core
build_sockets
build_websockets
build_bonjour
build_influxdb
build_http
build_mqtt
build_caf

exit 0;