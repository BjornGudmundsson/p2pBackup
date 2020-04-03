install = go get -u -tags=vartime -v ./...
build = go build
BINARY = a
port = 8084
udpPort = 3001
baseDir = "./test"
peerList = "./peers.txt"
rules = rules.toml
testKey = "3868d3a8fb5a8b8c234eeb20aac0d0de8377fb57ff68a7393468dfc5e338a7e7"
testAuthKey = "785535dbff6e8f2e08d9173e2fcd1cd77f5f09951b474ec77d14c8c924a50156"
setFile = set.txt
authFlag =  -authkey=$(testAuthKey) -set=$(setFile)
buildFlags = -o $(BINARY) -tags=vartime -v
binaryFlags = -p $(port) -base $(baseDir) -peers $(peerList) -udp $(udpPort) -backuprules $(rules) -key=$(testKey)
initFlag = -init=true 
MAKE = make
key2 = 187c3601cc7da912ed583244943b233d732182c1b88e7863ba941c114db8d97d

install:
	$(install)
build:
	$(build) $(buildFlags)

run: 
	$(MAKE) build; ./$(BINARY) $(binaryFlags) $(authFlag)
init:
	$(MAKE) build; ./$(BINARY) $(initFlag)

clean:
	rm $(BINARY); fuser -k $(port)/tcp; fuser -k $(udpPort)/udp;