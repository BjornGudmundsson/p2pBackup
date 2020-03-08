package peers

import (
	"bufio"
	"encoding/hex"
	"errors"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/BjornGudmundsson/p2pBackup/files"
)

type file files.File

//Peer is a container of how
//information about a peer is maintained.
type Peer struct {
	Name      string
	Addr      net.IP
	Port      int
	PublicKey []byte
	Suite     string
	TCP       int
}

//SuiteIsSupported take in a suite and says if that
//ciphersuite is supported.
func SuiteIsSupported(s string) bool {
	return true
}

//NewPeer takes in a description string of the form
//[Name IP hex_of_public_key CipherSuite_being_used]
//and returns a pointer to a peer if it is a valid string
//else it returns an error.
func NewPeer(desc string) (*Peer, error) {
	fields := strings.Split(desc, " ")
	if len(fields) != 6 {
		return nil, errors.New("Not enough fields")
	}
	p := &Peer{}
	p.Name = fields[0]
	ip := net.ParseIP(fields[1])
	if ip == nil {
		return nil, errors.New("Could not parse IP")
	}
	port, eport := strconv.Atoi(fields[2])
	if eport != nil {
		return nil, eport
	}
	tcpPort, etcp := strconv.Atoi(fields[3])
	if etcp != nil {
		return nil, etcp
	}
	p.TCP = tcpPort
	p.Port = int(port)
	p.Addr = ip
	d, ehex := hex.DecodeString(fields[4])
	if ehex != nil {
		return nil, ehex
	}
	p.PublicKey = d
	if !SuiteIsSupported(fields[5]) {
		return nil, errors.New("This suite is not supported")
	}
	p.Suite = fields[5]
	return p, nil
}

//GetPeerList takes in a file that has all the known peers
//and returns a slice with all of the peers in the file.
//If the file can't read or one of the peers is malformed it returns an error.
func GetPeerList(peerFile string) ([]*Peer, error) {
	f, e := os.Open(peerFile)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	peers := make([]*Peer, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		p, e := NewPeer(scanner.Text())
		if e != nil {
			return nil, e
		}
		peers = append(peers, p)
	}
	return peers, nil
}
