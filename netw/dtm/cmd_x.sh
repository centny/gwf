#!/bin/bash
set -e
sleep 1
echo $1 $2 $3
if [ "$1" == "exit" ];then
	exit 100
fi