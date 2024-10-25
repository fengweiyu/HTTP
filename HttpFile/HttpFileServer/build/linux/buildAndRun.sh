#!/bin/bash
cd ../../src/
go build HttpFileServer.go
cp ./HttpFileServer ../build/linux/
./HttpFileServer