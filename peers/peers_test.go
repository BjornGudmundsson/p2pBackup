package peers

import (
	"fmt"
	"testing"

	"github.com/BjornGudmundsson/p2pBackup/files"

	"github.com/stretchr/testify/assert"
)

func TestParsePeer(t *testing.T) {
	desc := "127.0.0.1 8080 " + "\n"
	p, e := NewTCPPeer(desc)
	fmt.Println("Peer: ", p )
	assert.Nil(t, e)
}

func TestFileParser(t *testing.T) {
	desc1 := "127.0.0.1 8080 " + tcp + "\n"
	l := len(desc1)
	desc2 := "127.0.0.1 8081 "+ tcp + "\n"
	f := files.File{
		Name: "peers.txt",
		Path: ".",
	}
	e := files.AppendToFile(f, []byte(desc1))
	assert.Nil(t, e)
	f.Size = int64(l)
	e = files.AppendToFile(f, []byte(desc2))
	assert.Nil(t, e)
	peers, epeers := GetPeerList("peers.txt")
	assert.Nil(t, epeers)
	assert.Len(t, peers, 2, "There should be exactly two peers")
	p1, p2 := peers[0], peers[1]
	fmt.Println(string(p1.Marshall()), string(p2.Marshall()))
	assert.Equal(t, p1.Address().String(), "127.0.0.1")
	assert.Equal(t, p2.Address().String(), "127.0.0.1")
	assert.Equal(t, p1.Port(), 8080, "Port should be 8080")
	assert.Equal(t, p2.Port(), 8081, "Port should be 8081")
	assert.Equal(t, tcp, p1.TransmissionProtocol(), "Protocol should be tcp")
	assert.Equal(t, tcp, p2.TransmissionProtocol(), "Protocol should be tcp")
}

