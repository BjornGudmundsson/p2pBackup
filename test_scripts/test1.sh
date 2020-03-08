

#Making the test directories
dummyContent="deadbeef lmao"
dir="testDir"
createDir="mkdir $dir"
createFiles="touch $dir/t1.txt $dir/t2.txt $dir/t3.txt"
$createDir
$createFiles
echo $dummyContent >> $dir/t1.txt
echo $dummyContent >> $dir/t2.txt
echo $dummyContent >> $dir/t3.txt
cat $dir/t1.txt
cat $dir/t2.txt
cat $dir/t3.txt
make build
./a -base=$dir &
p=$!
echo "Process running as: $p"
sleep 5s
kill -9 $p
cleanup="rm -rf testDir"
$cleanup