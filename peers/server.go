package peers

import (
	"fmt"
	"github.com/BjornGudmundsson/p2pBackup/files"
	"github.com/BjornGudmundsson/p2pBackup/utilities"
	"strconv"
	"time"
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
	if protocol == GIT {
		return NewGitServer(flags, container, enc), nil
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

func NewGitServer(flags utilities.Flags, container Container, enc *EncryptionInfo) func() error {
	f := func() error {
		go updatePeersGit(flags, container, enc)
		return FindPeersGit(flags, container, enc)
	}
	return f
}

func updatePeersGit(flags utilities.Flags, container Container, enc *EncryptionInfo) error {
	pingTimer := flags.GetString("ping")
	t, e := time.ParseDuration(pingTimer)
	if e != nil {
		return e
	}
	for {
		time.Sleep(t)
		pw, un, email := flags.GetString("gitpw"), flags.GetString("gituser"), flags.GetString("email")
		repo := flags.GetString("repo")
		port := flags.GetInt("fileport")
		ip := flags.GetString("ip")
		m := encodePeer(enc, ip+seperator+strconv.Itoa(port))
		e = PushMessageParallel(repo, un, pw, email, m)
		if e != nil {
			fmt.Println(e)
		}
	}
}

func encodePeer(enc *EncryptionInfo, p string) string {
	return p
}

func FindPeersGit(flags utilities.Flags, container Container, enc *EncryptionInfo) error {
	pingTimer := flags.GetString("ping")
	t, e := time.ParseDuration(pingTimer)
	if e != nil {
		return e
	}
	repo := flags.GetString("repo")
	pw, un := flags.GetString("gitpw"), flags.GetString("gituser")
	for {
		msgs, e := GetCommitMessages(repo, un, pw, t)
		if e != nil {
			fmt.Println(e)
			time.Sleep(t)
			continue
		}
		peers := make([]Peer, 0)
		for _, msg := range msgs {
			p, e := getPeerFromMsg(msg, container, enc)
			if e != nil {
				continue
			}
			peers = append(peers, p)
		}
		container.New(peers)
	}
	return nil
}

func getPeerFromMsg(msg string, container Container, enc *EncryptionInfo) (Peer, error) {
	return NewPeer(msg)
}