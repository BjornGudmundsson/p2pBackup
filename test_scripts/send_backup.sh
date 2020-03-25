
#Making the test directories
dummyContent="deadbeef lmao"
dir="testDir"
createDir="mkdir $dir"
createFiles="touch $dir/t1.txt $dir/t2.txt $dir/t3.txt"
key1="3868d3a8fb5a8b8c234eeb20aac0d0de8377fb57ff68a7393468dfc5e338a7e7"
key2="3861d11d29b10a0d6677bef2d675d73ca55a530e8093fe3bb956ac43da13ab48"
$createDir
$createFiles
#Populating the files such that they can be meaningfully backed up
echo $dummyContent >> $dir/t1.txt
echo $dummyContent >> $dir/t2.txt
echo $dummyContent >> $dir/t3.txt

touch log1.txt log2.txt backupfile1.txt backupfile2.txt

touch peers1.txt peers2.txt

#Setting up different peer files
echo "Ulf 127.0.0.1 8081 3000 ABCDEF ECDSA" >> peers1.txt
echo "Bjo 127.0.0.1 8082 3001 ABCDEF ECDSA" >> peers2.txt

make build

#By default the name of the binary is a

#Running the first peers
./a -peers=peers1.txt -udp=3000 -fileport=8081 -logfile=log1.txt -base=$dir -storage=backupfile1.txt -key="3868d3a8fb5a8b8c234eeb20aac0d0de8377fb57ff68a7393468dfc5e338a7e7"&
p1=$!

./a -peers=peers2.txt -udp=3001 -fileport=8082 -logfile=log2.txt -base=$dir -storage=backupfile2.txt -key="3861d11d29b10a0d6677bef2d675d73ca55a530e8093fe3bb956ac43da13ab48"&
p2=$!

sleep 5s
b1="$(cat backupfile1.txt)"
b2="$(cat backupfile2.txt)"
echo "$b1"
echo "$b2"

if [ "$b1" = "$b2" ]
then
    echo  -e "${GREEN}Passed"
else
    echo -e "${RED}Failed"
fi

#Cleaning up after the test
cleanup="rm -rf testDir"
kill -9 $p1
kill -9 $p2
rm log1.txt log2.txt peers1.txt peers2.txt backupfile1.txt backupfile2.txt
fuser -k 8081/tcp
fuser -k 8082/tcp
fuser -k 3000/udp
fuser -k 3001/udp
#make clean
$cleanup
