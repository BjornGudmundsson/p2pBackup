for i in 1 2 3 4 5
do
    echo "starting";
    bash test_scripts/retrieve_backup.sh | grep "Elapsed" >> m.txt;
    sleep 10s;
done