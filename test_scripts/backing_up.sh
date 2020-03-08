

#Making the test directories
dummyContent="deadbeef lmao"
dir="testDir"
createDir="mkdir $dir"
createFiles="touch $dir/t1.txt $dir/t2.txt $dir/t3.txt"
$createDir
$createFiles
#Populating the files such that they can be meaningfully backed up
echo $dummyContent >> $dir/t1.txt
echo $dummyContent >> $dir/t2.txt
echo $dummyContent >> $dir/t3.txt
cat $dir/t1.txt
cat $dir/t2.txt
cat $dir/t3.txt
#Building the project
make build
#Run the p2p-backupservice
./a -base=$dir &
p=$!

echo "Process running as: $p"
sleep 5s
#Kill the service
kill -9 $p
hash="e206ef3632fbf007b8f151ea08e4a1d8c95a5c4acac5a57fdaace8b24127d559"
log="$(cat backuplog.txt)"
echo $log
#Cleanup the test environment and clean up the program
cleanup="rm -rf testDir"
make clean
$cleanup