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
 rm -rf $GOPATH/src/github.com/Centny/gwf
 go get -u github.com/Centny/gwf
fi
#########Running Test#########
echo "Running Test"
# pkgs="\
#  github.com/Centny/gwf/smartio\
#  github.com/Centny/gwf/log\
#  github.com/Centny/gwf/util\
#  github.com/Centny/gwf/dbutil\
#  github.com/Centny/gwf/igtest\
#  github.com/Centny/gwf/routing\
#  github.com/Centny/gwf/routing/cookie\
#  github.com/Centny/gwf/routing/filter\
#  github.com/Centny/gwf/routing/httptest\
#  github.com/Centny/gwf/routing/doc\
#  github.com/Centny/gwf/jcr\
#  github.com/Centny/gwf/pathc\
#  github.com/Centny/gwf/mcobertura\
#  github.com/Centny/gwf/hooks\
#  github.com/Centny/gwf/ini\
#  github.com/Centny/gwf/tutil\
#  github.com/Centny/gwf/pool\
#  github.com/Centny/gwf/netw\
#  github.com/Centny/gwf/netw/impl\
#  github.com/Centny/gwf/im\
# "
pkgs="\
  github.com/Centny/gwf/im/imc\
"
echo "mode: set" > a.out
for p in $pkgs;
do
 go test -v --coverprofile=c.out $p
 cat c.out | grep -v "mode" >>a.out
 go install $p
done
gocov convert a.out > coverage.json

##############################
#####Create Coverage Report###
echo "Create Coverage Report"
cat coverage.json | gocov-xml -b $GOPATH/src > coverage.xml
cat coverage.json | gocov-html coverage.json > coverage.html

######
go install github.com/Centny/gwf
go install github.com/Centny/gwf/im/imc
