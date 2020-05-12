package peers

import (
	"bufio"
	"fmt"
	"golang.org/x/net/proxy"
	"net"
)

type Communicator interface {
	SendClearMessage(msg []byte) error
	SendMessage(msg []byte) error
	GetNextMessage() ([]byte, error)
	CloseChannel() error
}

type TCPCommunicator struct {
	c net.Conn
	encoder Encoder
}

func (com TCPCommunicator) SendMessage(msg []byte) error {
	encoder := com.encoder
	encoded := encoder.EncodeToString(msg) + eof
	_, e := fmt.Fprintf(com.c, encoded)
	return e
}

func (com TCPCommunicator) SendClearMessage(msg []byte) error {
	_, e := fmt.Fprintf(com.c, string(msg) + eof)
	return e
}
func (com TCPCommunicator) GetNextMessage() ([]byte, error) {
	encoder := com.encoder
	reader := bufio.NewReader(com.c)
	s, e := reader.ReadString('\n')
	if e != nil {
		return nil, e
	}
	l := len(s)
	if l <= 1 {
		return nil, new(ErrorEmptyData)
	}
	s = s[: l - 1]
	decoded, e := encoder.DecodeFromString(s)
	if e != nil {
		return nil, e
	}
	return decoded, nil
}

func (com TCPCommunicator) CloseChannel() error {
	return com.c.Close()
}

func NewTCPCommunicatorFromPeer(p Peer, enc *EncryptionInfo) (Communicator, error) {
	conn, e := getTCPConn(p)
	if e != nil {
		return nil, e
	}
	return TCPCommunicator{
		c: conn,
		encoder: enc.Enc,
	}, nil
}

func NewTorCommunicatorFromPeer(p Peer, enc *EncryptionInfo) (Communicator, error) {
	dialer, e := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, nil)
	if e != nil {
		return nil, e
	}
	c, e := dialer.Dial("tcp", p.ConnectorString())
	if e != nil {
		return nil, e
	}
	com := NewTCPCommunicatorFromConn(c, enc)
	return com, nil
}

func NewTCPCommunicatorFromConn(c net.Conn, enc *EncryptionInfo) Communicator {
	return TCPCommunicator{
		c: c,
		encoder: enc.Enc,
	}
}

func NewCommunicatorFromPeer(p Peer, enc *EncryptionInfo) (Communicator, error) {
	protocol := p.TransmissionProtocol()
	if protocol == tcp {
		return NewTCPCommunicatorFromPeer(p, enc)
	}
	if protocol == tor {
		return NewTorCommunicatorFromPeer(p, enc)
	}
	return nil, nil
}