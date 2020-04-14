peerFile = peers.txt
createPeerfile = touch $(peerFile);
goTest = go test -v -tags=vartime;

test:
	$(createPeerfile) $(goTest) rm $(peerFile)