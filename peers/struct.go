package peers

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/BjornGudmundsson/p2pBackup/crypto"
	"github.com/BjornGudmundsson/p2pBackup/files"
	"github.com/BjornGudmundsson/p2pBackup/kyber"
	"github.com/BjornGudmundsson/p2pBackup/kyber/sign/schnorr"
	"github.com/BjornGudmundsson/p2pBackup/kyber/util/random"
	"github.com/BjornGudmundsson/p2pBackup/purb"
	"github.com/BjornGudmundsson/p2pBackup/purb/purbs"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type file files.File

type Peer interface {
	Address() net.IP
	Port() int
	LastSeen() time.Time
	fmt.Stringer
	TransmissionProtocol() string
	Available() bool
	Marshall() []byte
	Unmarshall(d []byte) error
	ConnectorString() string
}

//TCPPeer is a container of how
//information about a peer is maintained.
type TCPPeer struct {
	Addr      net.IP
	port      int
	seen time.Time
	available bool
	protocol string
}

func (p *TCPPeer) Address() net.IP {
	return p.Addr
}

func (p *TCPPeer) Port() int {
	return p.port
}

func (p *TCPPeer) LastSeen() time.Time {
	return p.seen
}

func (p *TCPPeer) Available() bool {
	return p.available
}

func (p *TCPPeer) String() string {
	return "Addr: " + p.Address().String() + " " + "Port: " + strconv.Itoa(p.Port())
}

func (p *TCPPeer) TransmissionProtocol() string {
	return p.protocol
}

func (p *TCPPeer) Marshall() []byte {
	s := p.Addr.String() + seperator + strconv.Itoa(p.port) + seperator + p.protocol
	return []byte(s)
}

func (p *TCPPeer) Unmarshall(d []byte) error {
	p2, e := NewTCPPeer(string(d))
	if e != nil {
		return e
	}
	p.protocol = p2.TransmissionProtocol()
	p.port = p2.Port()
	p.Addr = p2.Address()
	p.seen = time.Now()
	p.available = true
	return nil
}

func (p *TCPPeer) ConnectorString() string {
	return p.Address().String()+":"+strconv.Itoa(p.Port())
}

//NewTCPPeer takes in a description string of the form
//[Name IP hex_of_public_key CipherSuite_being_used]
//and returns a pointer to a peer if it is a valid string
//else it returns an error.
func NewTCPPeer(desc string) (Peer, error) {
	fields := strings.Split(desc, " ")
	if len(fields) != 4  && len(fields) != 3 {
		return nil, errors.New("not the right amount of fields")
	}
	p := &TCPPeer{}
	ip := net.ParseIP(fields[0])
	if ip == nil {
		return nil, errors.New("could not parse IP")
	}
	port,e := strconv.Atoi(fields[1])
	if e != nil {
		return nil, e
	}
	p.Addr = ip
	p.port = port
	p.seen = time.Now()
	p.available = true
	p.protocol = fields[2]
	return p, nil
}

//GetPeerList takes in a file that has all the known peers
//and returns a slice with all of the peers in the file.
//If the file can't read or one of the peers is malformed it returns an error.
func GetPeerList(peerFile string) ([]Peer, error) {
	f, e := os.Open(peerFile)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	peers := make([]Peer, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		if txt == "" {
			continue
		}
		p, e := NewTCPPeer(txt)
		if e != nil {
			fmt.Println(txt)
			return nil, e
		}
		peers = append(peers, p)
	}
	return peers, nil
}

//EncryptionInfo keeps track of all of the
//cryptographic information needed to take part
//in the protocol.
type EncryptionInfo struct {
	Enc Encoder
	Auth crypto.Authenticator
	AuthKey kyber.Scalar
	Link []byte
	RetrievalInfo *purb.KeyInfo
	Password string//This parameter is entirely optional
	RecipientKeys []kyber.Point//This parameter is entirely optional
}

func NewEncryptionInfo(auth crypto.Authenticator, authKey kyber.Scalar, link []byte, info *purb.KeyInfo, pw string, recipients []kyber.Point) *EncryptionInfo {
	return &EncryptionInfo{
		Enc: NewB64Encoder(),
		Auth:          auth,
		AuthKey:       authKey,
		Link:          link,
		RetrievalInfo: info,
		Password:      pw,
		RecipientKeys:  recipients,
	}
}

func (enc *EncryptionInfo) Sign(msg []byte) ([]byte, error) {
	return enc.Auth.Sign(enc.AuthKey, msg, enc.Link)
}

func (enc *EncryptionInfo) NormalSign(msg []byte) ([]byte, error) {
	return schnorr.Sign(enc.Auth.GetSuite(), enc.AuthKey, msg)
}

func (enc *EncryptionInfo) NormalVerify(msg, sig []byte) error {
	set := enc.Auth.GetAnonSet()
	for _, p := range set {
		e := schnorr.Verify(enc.Auth.GetSuite(), p, msg, sig)
		if e != nil {
			return nil
		}
	}
	return errors.New("could not verify")
}

func (enc *EncryptionInfo) Verify(msg, sig []byte) ([]byte, error) {
	return enc.Auth.Verify(msg, sig, enc.Link)
}

//PURBAnon takes in a data that is supposed to be encrypted for the entire anonymity set. I.e. anyone in
//the peer-group can decode it using their secret key.
func (enc *EncryptionInfo) PURBAnon(d []byte) ([]byte, error) {
	params := purbs.NewPublicFixedParameters(enc.RetrievalInfo.SuiteInfos, false)
	suite := enc.RetrievalInfo.Suite
	recipients := make([]purbs.Recipient, len(enc.Auth.GetAnonSet()))
	for i, p := range enc.Auth.GetAnonSet() {
		r, e := pointToRecipient(p, suite)
		if e != nil {
			return nil,  e
		}
		recipients[i] = *r
	}
	pur, e := purbs.Encode(d, recipients, random.New(), params, false)
	if e != nil {
		return nil, e
	}
	return pur.ToBytes(), nil
}

func pointToRecipient(p kyber.Point, s purbs.Suite) (*purbs.Recipient, error) {
	m, e := p.MarshalBinary()
	if e != nil {
		return nil, e
	}
	r, e := purb.NewRecipient(m, s)
	if e != nil {
		return nil, e
	}
	return &r, e
}

func scalarToRecipient(x kyber.Scalar, s purbs.Suite) (*purbs.Recipient, error) {
	m, e := x.MarshalBinary()
	if e != nil {
		return nil, e
	}
	r, e := purb.NewPrivateRecipient(m, s)
	if e != nil {
		return nil, e
	}
	return &r, nil
}

func (enc *EncryptionInfo) PURBBackup(d []byte) ([]byte, error) {
	x := enc.RetrievalInfo.PrivateKey
	pw := enc.Password
	suite := enc.RetrievalInfo.Suite
	pwKey, e := crypto.PrivateKeyFromPassword(pw, suite)
	if e != nil {
		return nil, e
	}
	pwPoint := suite.Point().Base()
	pwPoint = pwPoint.Mul(pwKey, pwPoint)
	pwRecipient, e := pointToRecipient(pwPoint, suite)
	if e != nil {
		return nil, e
	}
	publicKey := suite.Point().Base()
	publicKey = publicKey.Mul(x, publicKey)
	selfRecipient, e := pointToRecipient(publicKey, suite)
	if e != nil {
		return nil, e
	}
	recipients := make([]purbs.Recipient, 0)
	recipients = append(recipients, *pwRecipient)
	recipients = append(recipients, *selfRecipient)
	for _, p  := range enc.RecipientKeys {
		r, e := pointToRecipient(p, suite)
		if e != nil {
			return nil, e
		}
		recipients = append(recipients, *r)
	}
	params := purbs.NewPublicFixedParameters(enc.RetrievalInfo.SuiteInfos, false)
	pur, e := purbs.Encode(d, recipients,random.New(), params, false)
	if e != nil {
		return nil, e
	}
	blob := pur.ToBytes()
	return blob, nil
}

func (enc *EncryptionInfo) DecodePURBBackup(blob []byte) ([]byte, error) {
	params := purbs.NewPublicFixedParameters(enc.RetrievalInfo.SuiteInfos, false)
	x := enc.RetrievalInfo.PrivateKey
	suite := enc.RetrievalInfo.Suite
	r, e := scalarToRecipient(x, suite)
	if e == nil {
		v, d, e := purbs.Decode(blob, r, params, false)
		if e == nil {
			if v {
				return d, nil
			}
		}
	}
	pw := enc.Password
	k, e := crypto.PrivateKeyFromPassword(pw, suite)
	if e == nil {
		kr, e := scalarToRecipient(k, suite)
		if e == nil {
			v, d, e := purbs.Decode(blob, kr, params, false)
			if e == nil && v {
				return d, nil
			}
		}
	}
	return nil, new(ErrorCouldNotDecode)
}

func (enc *EncryptionInfo) DecodePURBAnon(blob []byte) ([]byte, error) {
	suite := enc.RetrievalInfo.Suite
	x := enc.AuthKey
	m, e := x.MarshalBinary()
	if e != nil {
		return nil, e
	}
	recipient, e := purb.NewPrivateRecipient(m, suite)
	if e != nil {
		return nil, e
	}
	params := purbs.NewPublicFixedParameters(enc.RetrievalInfo.SuiteInfos, false)
	v, d, e := purbs.Decode(blob, &recipient, params, false)
	if e != nil {
		return nil, e
	}
	if !v {
		return nil, new(ErrorCouldNotDecode)
	}
	return d, nil
}

type BackupInfo struct {
	X kyber.Scalar//The secret to the backup.
	StartIndex int64
	Size int64
}