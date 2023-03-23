#!/bin/bash
cd ~/*go
cp -R ~/gitee/uwbwebapp/web/* ./web/
./kill.sh
cd ~/gitee/uwbwebapp
go build -o ~/uwbwebapp_go/cmd/main ./main.go
go build -o ~/uwbwebapp_go/cmd/logger_client ./pkg/tools/logger_client/loggerclient.go
go build -o ./cmd/logger_client ./pkg/tools/logger_client/loggerclient.go
cd ~/*go
./run.sh
cd ~/gitee/uwbwebapp
# cp conf/WebConfig.json ../uwbwebapp_go/conf

