export GOPATH=$PWD
mkdir -p bin
cd bin
go build zoneserver
go build webserver
