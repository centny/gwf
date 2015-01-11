#!/bin/bash
go install github.com/Centny/gwf/netw/atestc
$GOPATH/bin/atestc 192.168.1.130:7686 1>o.log 2>e.log