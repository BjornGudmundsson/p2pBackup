install = go get -u -tags=vartime -v ./...
build = go build
BINARY = a
port = 8084
udpPort = 3001
baseDir = "./test"
peerList = "./peers.txt"
rules = rules.toml
testKey = "3868d3a8fb5a8b8c234eeb20aac0d0de8377fb57ff68a7393468dfc5e338a7e7"
buildFlags = -o $(BINARY) -tags=vartime -v
binaryFlags = -p $(port) -base $(baseDir) -peers $(peerList) -udp $(udpPort) -backuprules $(rules) -key=$(testKey)
initFlag = -init=true 
MAKE = make

install:
	$(install)
build:
	$(build) $(buildFlags)

run: 
	$(MAKE) build; ./$(BINARY) $(binaryFlags)
init:
	$(MAKE) build; ./$(BINARY) $(initFlag)

clean:
	rm $(BINARY); fuser -k $(port)/tcp; fuser -k $(udpPort)/udp;