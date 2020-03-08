package peers

import (
	"testing"

	"github.com/BjornGudmundsson/p2pBackup/files"

	"github.com/stretchr/testify/assert"
)

func TestParsePeer(t *testing.T) {
	desc := "Bjo 127.0.0.1 8080 3001 ABCDEF " + ECDSA
	_, e := NewPeer(desc)
	assert.Nil(t, e)
}

func TestFileParser(t *testing.T) {
	desc1 := "Bjo 127.0.0.1 8080 3001 ABCDEF " + ECDSA + "\n"
	l := len(desc1)
	desc2 := "Ulf 127.0.0.1 8081 3000 ABCDEF " + ECDSA + "\n"
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
	assert.Equal(t, p1.Name, "Bjo", "Names should equal")
	assert.Equal(t, p2.Name, "Ulf", "Names should equal")
	assert.Equal(t, p1.Addr.String(), "127.0.0.1")
	assert.Equal(t, p2.Addr.String(), "127.0.0.1")
	assert.Equal(t, p1.Port, 8080, "Port should be 8080")
	assert.Equal(t, p2.Port, 8081, "Port should be 8081")

}
