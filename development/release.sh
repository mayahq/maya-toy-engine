#!/bin/bash -e

CASCADES_VERSION=$1
#OSARCHS=( "darwin/amd64" "linux/amd64" "linux/arm" "windows/amd64" );
OSARCHS=( "linux/amd64" );

source func.sh

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

if [ "$CASCADES_VERSION" = "" ]; then
    CASCADES_VERSION="DEV"
fi

if [ -d "$DIR/dist" ]; then
    rm -rf $DIR/dist/*
else
    mkdir -p $DIR/dist
fi


for v in "${OSARCHS[@]}"
do
    suffix=`echo "${v}" | tr / _`
    p="cascades-${CASCADES_VERSION}_${suffix}"
    d="${DIR}/dist/${p}"
    mkdir -p $d/components

    # Compile
    pushd $d
    gox -verbose -output="cascades" -osarch="${v}" github.com/cascades-fbp/cascades/cmd/cascades
    build_core
    build_sockets
    build_websockets
    build_bonjour
    build_influxdb
    build_http
    build_mqtt
    popd

    mkdir $d/examples
    cp $DIR/examples/tick.fbp $d/examples/
    cp $DIR/examples/tick.json $d/examples/
    cp $DIR/examples/regexp2json.fbp $d/examples/

    cd $DIR/dist && zip -9 -r "${p}.zip" "$p"
    #cd $DIR/dist && tar -zcvf "${p}.tar.gz" $p
    rm -rf $d
done

echo "DONE!"
exit 0