#!/bin/bash
##############################
#####Setting Environments#####
echo "Setting Environments"
set -e
export PATH=$PATH:$GOPATH/bin:$HOME/bin:$GOROOT/bin
##############################
######Install Dependence######
echo "Installing Dependence"
#go get github.com/go-sql-driver/mysql
#go get github.com/Centny/TDb
#go get code.google.com/p/go-uuid/uuid
##############################
#########Running Clear#########
if [ "$1" = "-u" ];then
 echo "Running Clear"
 rm -rf $GOPATH/src/github.com/Centny/Cny4go
 go get -u github.com/Centny/Cny4go
fi
#########Running Test#########
echo "Running Test"
pkgs="\
 github.com/Centny/Cny4go/smartio\
 github.com/Centny/Cny4go/log\
 github.com/Centny/Cny4go/util\
 github.com/Centny/Cny4go/dbutil\
 github.com/Centny/Cny4go/igtest\
 github.com/Centny/Cny4go/routing\
 github.com/Centny/Cny4go/routing/cookie\
 github.com/Centny/Cny4go/routing/filter\
 github.com/Centny/Cny4go/routing/httptest\
 github.com/Centny/Cny4go/jcr\
 github.com/Centny/Cny4go/pathc\
 github.com/Centny/Cny4go/mcobertura\
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

######
go install github.com/Centny/Cny4go/mrepo
go install github.com/Centny/Cny4go/mcobertura
go install github.com/Centny/Cny4go/jcr/jcr
go install github.com/Centny/Cny4go/pathc
