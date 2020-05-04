package peers

import (
	"flag"
	"github.com/BjornGudmundsson/p2pBackup/files"
)

type Server interface {
	Listen() error
	FindPeers() error
}

type TCPUDPServer struct {
	tcpPort string
	udpPort string
	bh files.BackupHandler
	enc *EncryptionInfo
	listen func(enc *EncryptionInfo, bh files.BackupHandler) error
	find func (enc *EncryptionInfo, container Container) error
}

func NewServer() (Server, error) {
	protocol := flag.String("protocol", tcp, "Which protocol to be used")
	flag.Parse()
	if *protocol == tcp {

	}
	return nil, nil
}

func NewTCPServer() func(*EncryptionInfo, files.BackupHandler) error {
	//port := flag.String("p", "8080", "Which port to run the server on")
	//flag.Parse()
	return nil
}