#!/bin/bash
set -e
for((i=1;i<3;i++));
do
	sleep $1
	url="$2&process=0.$i"
	echo $url
	curl -s -o a.txt $url
	echo running $i
done
echo 11111 $1