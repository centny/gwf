#!/bin/bash
go install github.com/Centny/gwf/netw/atests
$GOPATH/bin/atests 1>o.log 2>e.log