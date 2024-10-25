#!/bin/bash

if [ $# == 0 ]; then
    goPWD=$PWD/../../src
    cd $goPWD
    go env -w GOPROXY=https://goproxy.cn
    go get golang.org/x/net/websocket
fi

#go get go.mongodb.org/mongo-driver/bson
#go get go.mongodb.org/mongo-driver/mongo
#go get go.mongodb.org/mongo-driver/mongo/options