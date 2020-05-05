package peers

import (
	"github.com/BjornGudmundsson/p2pBackup/files"
	"github.com/BjornGudmundsson/p2pBackup/utilities"
	"strconv"
)

type Server interface {
	Listen() error
	FindPeers() error
}

type ServerImplementation struct {
	bh files.BackupHandler
	enc *EncryptionInfo
	listen func() error
	find func () error
}

func NewServer(flags utilities.Flags, handler files.BackupHandler, enc *EncryptionInfo, container Container) (Server, error) {
	server := &ServerImplementation{}
	protocol := flags.GetString("protocol")
	listen, e := getFileProtocol(flags, protocol, handler, enc)
	if e != nil {
		return nil, e
	}
	server.listen = listen
	find := flags.GetString("find")
	findProtocol, e := getFindProtocol(flags, find, enc, container)
	if e != nil {
		return nil, e
	}
	server.find = findProtocol
	return server, nil
}

func getFileProtocol(flags utilities.Flags, protocol string, handler files.BackupHandler, enc *EncryptionInfo) (func() error, error) {
	if protocol == tcp {
		return NewTCPServer(flags, handler, enc), nil
	}
	return nil, new(ErrorProtocolNotFound)
}

func getFindProtocol(flags utilities.Flags, protocol string, enc *EncryptionInfo, container Container) (func() error, error) {
	if protocol == udp {
		return NewUDPServer(flags, container, enc), nil
	}
	return nil, new(ErrorProtocolNotFound)
}

func (s *ServerImplementation) Listen() error {
	return s.listen()
}

func (s *ServerImplementation) FindPeers() error {
	return s.find()
}

func NewTCPServer(flags utilities.Flags, handler files.BackupHandler, enc *EncryptionInfo) func() error {
	port := strconv.Itoa(flags.GetInt("fileport"))
	f := func() error {
		e := ListenTCP(":" + port, enc, handler)
		return e
	}
	return f
}

func NewUDPServer(flags utilities.Flags, container Container, enc *EncryptionInfo) func() error {
	port := flags.GetString("udp")
	f := func() error {
		return ListenUDP(":" + port, container, enc)
	}
	return f
}