./b.sh
for((i=0;i<200;i++));
do
 rcmd -c 192.168.1.130:7686 run
 # rcmd -c 192.168.1.130:7686 heap
 echo $i
 sleep 1
done
