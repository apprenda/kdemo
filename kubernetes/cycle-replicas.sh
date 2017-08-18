for (( ; ; ))
do
   echo "Setting Replicas to 16 [ hit CTRL+C to stop]"
   kubectl scale deployment kdemo-dep --replicas=16
   sleep 10
   echo "Setting Replicas to 32 [ hit CTRL+C to stop]"
   kubectl scale deployment kdemo-dep --replicas=32
   sleep 10
   echo "Setting Replicas to 1 [ hit CTRL+C to stop]"
   kubectl scale deployment kdemo-dep --replicas=1
   sleep 10
done