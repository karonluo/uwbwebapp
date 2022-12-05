#!/bin/bash
cd ../*go
./kill.sh
cd ../uwbwebapp
go build -o ../uwbwebapp_go/cmd/main main.go
cd ../*go
./run.sh
cd ../uwbwebapp
# cp conf/WebConfig.json ../uwbwebapp_go/conf

