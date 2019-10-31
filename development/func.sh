#!/bin/bash

###########################################################################
# Functions
#

# display a message in red with a cross by it
# example
# echo echo_fail "No"
function echo_fail {
  # echo first argument in red
  printf "\e[31m✘ ${1}"
  # reset colours back to normal
  printf "\033[0m\n"
}

# display a message in green with a tick by it
# example
# echo echo_fail "Yes"
function echo_pass {
  # echo first argument in green
  printf "\e[32m✔ ${1}"
  # reset colours back to normal
  printf "\033[0m\n"
}

buildComponents() {
    components=$1
    for c in "${components[@]}"
    do
        echo -n " -> ${c}... "
        go get -d $2
        if [ "$3" = "" ]; then
            go build -o components/$c $2/$c
        else
            go build -o components/$3/$c $2/$c
        fi
        echo_pass
    done
}

build_core() {
    echo "Building core components:"
    components=('core/console' 'core/delay' 'core/drop' 'core/exec' \
        'core/distinct' \
        'core/joiner' 'core/passthru' 'core/readfile' 'core/splitter' \
        'core/submatch' 'core/switch' 'core/template' 'core/ticker' \
        'debug/crasher' 'fs/walk' 'fs/watchdog')
    buildComponents $components "github.com/cascades-fbp/cascades/components"
}

build_sockets() {
    echo "Building sockets component:"
    components=('tcp-server')
    buildComponents $components "github.com/cascades-fbp/cascades-sockets" "sockets"
}

build_websockets() {
    echo "Building websocket components:"
    components=('client' 'server')
    buildComponents $components "github.com/cascades-fbp/cascades-websocket" "websocket"
}

build_bonjour() {
    echo "Building bonjour components:"
    components=('discover' 'register')
    buildComponents $components "github.com/cascades-fbp/cascades-bonjour" "bonjour"
}

build_influxdb() {
    echo "Building InfluxDB components:"
    components=('write')
    buildComponents $components "github.com/cascades-fbp/cascades-influxdb" "influxdb"
}

build_http() {
    echo "Building http components:"
    #components=('client' 'router' 'server')
    components=('client')
    buildComponents $components "github.com/cascades-fbp/cascades-http" "http"
}

build_mqtt() {
    echo "Building MQTT components:"
    components=('pub' 'sub')
    buildComponents $components "github.com/cascades-fbp/cascades-mqtt" "mqtt"
}

build_caf() {
    echo "Building CAF components:"
    components=('http-property' 'context')
    buildComponents $components "github.com/cascades-fbp/cascades-caf" "caf"
}
