package peers

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/BjornGudmundsson/p2pBackup/kyber"
	"github.com/BjornGudmundsson/p2pBackup/kyber/util/random"
	"github.com/BjornGudmundsson/p2pBackup/purb"
	"github.com/BjornGudmundsson/p2pBackup/purb/purbs"
	"net"
	"strconv"
	"strings"
)

const delim = ";"
const seperator = " "

const errorIndicator = "Error: "


//signPublicKey take a marshalled public key and returns the same key along
// with a valid signature of that key.
func signPublicKey(k []byte, signer *EncryptionInfo) ([]byte, error) {
	sig, e := signer.Sign(k)
	if e != nil {
		return nil, e
	}
	hexSig := signer.Enc.EncodeToString(sig)
	hexKey := signer.Enc.EncodeToString(k)
	return []byte(hexSig + seperator + hexKey), nil
}

//verifyPublicKey takes in a public key and a signature and validates them and returns a new EC-point
//for the corresponding key from the alleged suite.
func verifyPublicKey(d []byte, verifier *EncryptionInfo, suite purbs.Suite) (kyber.Point, error) {
	b, e := verifier.Enc.DecodeFromString(string(d))
	if e != nil {
		return nil, e
	}
	sigKey := strings.Split(string(b), seperator)
	if len(sigKey) != 2 {
		return nil, new(ErrorIncorrectFormat)
	}
	sig, e := verifier.Enc.DecodeFromString(sigKey[0])
	if e != nil {
		return nil, e
	}
	key, e := verifier.Enc.DecodeFromString(sigKey[1])
	if e != nil {
		return nil, e
	}
	_, e = verifier.Verify(key, sig)
	if e != nil {
		return nil, e
	}
	p := suite.Point()
	e = p.UnmarshalBinary(key)
	return p, e
}

//verifyPURB takes in PURBified data that is supposed to be shared with the key corresponding to the scalar
//given. The de-purbified data is supposed to be signed and then it can be verifed that it was encoded by
//someone from the allowed group.
func verifyPURB(x kyber.Scalar, suite purbs.Suite, blob []byte, verifier *EncryptionInfo) ([]byte, error) {
	params := purbs.NewPublicFixedParameters(verifier.RetrievalInfo.SuiteInfos, false)
	b, e := x.MarshalBinary()
	if e != nil {
		return nil, e
	}
	recipient, e := purb.NewPrivateRecipient(b, suite)
	if e != nil {
		return nil, e
	}
	v, d, e := purbs.Decode(blob, &recipient, params, false)
	if e != nil {
		return nil, e
	}
	if !v {
		return nil, new(ErrorCouldNotDecode)
	}
	s := string(d)
	spl := strings.Split(s, seperator)
	if len(spl) < 2 {
		return nil, new(ErrorIncorrectFormat)
	}
	sig := spl[0]
	sigDecoded, e := verifier.Enc.DecodeFromString(sig)
	if e != nil {
		return nil, e
	}
	data := []byte(strings.Join(spl[1:], seperator))
	_, e = verifier.Verify(data, sigDecoded)
	return data, e
}

//signAndPURB takes in a piece of data, signs it and the purbifies the data given along with it signature.
func signAndPURB(signer *EncryptionInfo, recipients []purbs.Recipient, suite purbs.Suite, data []byte) ([]byte, error) {
	sig, e := signer.Sign(data)
	if e != nil {
		return nil,  e
	}
	fmt.Println("Len sig: ", len(sig))
	signedBlob := []byte(signer.Enc.EncodeToString(sig) + seperator + string(data))
	params := purbs.NewPublicFixedParameters(signer.RetrievalInfo.SuiteInfos, false)
	p, e := purbs.Encode(signedBlob, recipients, random.New(), params, false)
	if e != nil {
		return nil, e
	}
	return p.ToBytes(), nil
}


func getTCPConn(p Peer) (net.Conn, error) {
	c, e := net.Dial("tcp", p.Address().String()+":"+strconv.Itoa(p.Port()))
	return c, e
}

type Encoder interface {
	EncodeToString(d []byte) string
	DecodeFromString(s string) ([]byte, error)
}

type hexEncoder struct{}

func (he hexEncoder) EncodeToString(d []byte) string {
	return hex.EncodeToString(d)
}

func (he hexEncoder) DecodeFromString(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

type base64Encoder struct {}

func (b64 base64Encoder) EncodeToString(d []byte) string {
	return base64.StdEncoding.EncodeToString(d)
}

func (b64 base64Encoder) DecodeFromString(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func NewHexEncoder() Encoder {
	return hexEncoder{}
}

func NewB64Encoder() Encoder {
	return base64Encoder{}
}
