#!/bin/bash
##############################
#####Setting Environments#####
echo "Setting Environments"
set -e
export cpwd=`pwd`
export LD_LIBRARY_PATH=/usr/local/lib:/usr/lib
export PATH=$PATH:$GOPATH/bin:$HOME/bin:$GOROOT/bin
o_dir=build
if [ "$2" != "" ];then
	o_dir=$2/dtmc
fi
rm -rf $o_dir
mkdir $o_dir

#### Package ####
n_srv=dtmc
v_srv=0.0.1
##
d_srv="$n_srv"d
o_srv=$o_dir/$n_srv
mkdir $o_srv
mkdir $o_srv/conf
mkdir $o_srv/www
go build -o $o_srv/$n_srv github.com/Centny/gwf/cmd/dtmc
cp $n_srv.properties $o_srv/conf

###
if [ "$1" != "" ];then
	curl -o $o_srv/srvd $1/srvd
	curl -o $o_srv/srvd_i $1/srvd_i
	chmod +x $o_srv/srvd
	chmod +x $o_srv/srvd_i
	echo "./srvd_i \$1 dtmc \$2 \$3" >$o_srv/install.sh
	chmod +x $o_srv/install.sh
fi 
cd $o_dir
zip -r -q $n_srv.zip $n_srv
cd ../
echo "Package $n_srv..."
