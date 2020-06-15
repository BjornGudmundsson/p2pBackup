n=20
for i in $(seq 1 $n);
do
    echo "starting";
    bash latency_test.sh | grep "Elapsed:" >> measurement.txt
    sleep 30;
done;   