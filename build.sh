#!/bin/bash
##############################
#####Setting Environments#####
echo "Setting Environments"
set -e
export PATH=$PATH:$GOPATH/bin:$HOME/bin:$GOROOT/bin
##############################
######Install Dependence######
echo "Installing Dependence"
if [ "$1" = "-u" ];then
  echo "updating github.com/Centny/Cny4go code..."
  go get -u github.com/Centny/Cny4go
fi
##############################
#########Running Test#########
echo "Running Test"
pkgs="\
 github.com/Centny/Cny4go/smartio\
 github.com/Centny/Cny4go/log\
 github.com/Centny/Cny4go/util\
 github.com/Centny/Cny4go/dbutil\
"
echo "mode: set" > a.out
for p in $pkgs;
do
 go test -v --coverprofile=c.out $p
 cat c.out | grep -v "mode" >>a.out
done
gocov convert a.out > coverage.json

##############################
#####Create Coverage Report###
echo "Create Coverage Report"
cat coverage.json | gocov-xml -b $GOPATH/src > coverage.xml
cat coverage.json | gocov-html coverage.json > coverage.html
