build = go build
BINARY = a
port = 8084
udpPort = 3001
baseDir = "./test"
peerList = "./peers.txt"
rules = rules.toml
buildFlags = -o $(BINARY)
binaryFlags = -p $(port) -base $(baseDir) -peers $(peerList) -udp $(udpPort) -backuprules $(rules)
MAKE = make


build:
	$(build) $(buildFlags)

run: 
	$(MAKE); ./$(BINARY) $(binaryFlags)

clean:
	rm $(BINARY); fuser -k $(port)/tcp; fuser -k $(udpPort)/udp;