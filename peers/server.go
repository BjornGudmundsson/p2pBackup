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
	
	flag.Parse()
	return nil, nil
}
