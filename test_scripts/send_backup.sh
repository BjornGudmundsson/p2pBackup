RED='\033[0;31m'
NC='\033[0m'
GREEN='\033[0;32m'
#Making the test directories
dummyContent="deadbeef lmao"
dir="testDir"
createDir="mkdir $dir"
createFiles="touch $dir/t1.txt $dir/t2.txt $dir/t3.txt"
key1="3868d3a8fb5a8b8c234eeb20aac0d0de8377fb57ff68a7393468dfc5e338a7e7"
key2="3861d11d29b10a0d6677bef2d675d73ca55a530e8093fe3bb956ac43da13ab48"
authKey1="60bbcd3c9e79b514d10a12e8242b18fbdf96e51fd3bf7798b89af31515f956db"
publicAuthKey1="b9a616d0453af330d462a8afba820f4f22c77bac2d09aa60c17913ad9416c946"
authKey2="7029b076b7528d0fd076f9e7a31f1a0e7b5e83012b0fea2891dcf55f1c0b19f5"
publicAuthKey2="6eea06f2136b9c5bdebfcb9980b5a56fbdf315dc03304a78bc6451436b208a82"
setFlag="-set=set.txt"
$createDir
$createFiles
#Populating the files such that they can be meaningfully backed up
echo $dummyContent >> $dir/t1.txt;
echo $dummyContent >> $dir/t2.txt;
echo $dummyContent >> $dir/t3.txt;

touch log1.txt log2.txt backupfile1.txt backupfile2.txt set.txt;
echo $publicAuthKey1 >> set.txt;
echo $publicAuthKey2 >> set.txt;
touch peers1.txt peers2.txt;
#echo "Bjorn er cool" >> backupfile2.txt;
#echo "Bjorn er cool" >> backupfile1.txt;
echo "127.0.0.1 8081 " >> peers1.txt;
echo "127.0.0.1 8082 " >> peers2.txt;

make build;
#By default the name of the binary is a

#Running the first peer
./a -peers=peers1.txt -udp=3000 -fileport=8081 -logfile=log1.txt -base=$dir -storage=backupfile1.txt -key="$key1"  $setFlag -authkey="$authKey1" &
p1=$!;

#Running the second peer
./a -peers=peers2.txt -udp=3001 -fileport=8082 -logfile=log2.txt -base=$dir -storage=backupfile2.txt -key="$key2" $setFlag -authkey="$authKey2" &
p2=$!;

sleep 5s
b1="$(cat backupfile1.txt)";
b2="$(cat backupfile2.txt)";
#echo "$b1";
#echo "$b2";

#Make sure that both strings are non-empty and equal
if [ -n "$b1" -a -n "$b2" ]
then
    echo  -e "${GREEN}Passed${NC}"
else
    echo -e "${RED}Failed${NC}"
fi

rm set.txt
#Cleaning up after the test
cleanup="rm -rf testDir";
kill -9 $p1;
kill -9 $p2;
rm log1.txt log2.txt peers1.txt peers2.txt backupfile1.txt backupfile2.txt;
fuser -k 8081/tcp;
fuser -k 8082/tcp;
fuser -k 3000/udp;
fuser -k 3001/udp;
#make clean
$cleanup;
