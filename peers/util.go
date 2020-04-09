package peers

import (
	"bufio"
	"encoding/hex"
	"github.com/BjornGudmundsson/p2pBackup/kyber"
	"github.com/BjornGudmundsson/p2pBackup/kyber/util/random"
	"github.com/BjornGudmundsson/p2pBackup/purb"
	"github.com/BjornGudmundsson/p2pBackup/purb/purbs"
	"net"
	"strconv"
	"strings"
)

const delim = ";"

const errorIndicator = "Error: "


//signPublicKey take a marshalled public key and returns the same key along
// with a valid signature of that key.
func signPublicKey(k []byte, signer *EncryptionInfo) ([]byte, error) {
	sig, e := signer.Sign(k)
	if e != nil {
		return nil, e
	}
	hexSig := hex.EncodeToString(sig)
	hexKey := hex.EncodeToString(k)
	return []byte(hexSig + delim + hexKey), nil
}

//verifyPublicKey takes in a public key and a signature and validates them and returns a new EC-point
//for the corresponding key from the alleged suite.
func verifyPublicKey(d []byte, verifier *EncryptionInfo, suite purbs.Suite) (kyber.Point, error) {
	b, e := hex.DecodeString(string(d))
	if e != nil {
		return nil, e
	}
	sigKey := strings.Split(string(b), delim)
	if len(sigKey) != 2 {
		return nil, new(ErrorIncorrectFormat)
	}
	sig, e := hex.DecodeString(sigKey[0])
	if e != nil {
		return nil, e
	}
	key, e := hex.DecodeString(sigKey[1])
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
	spl := strings.Split(s, delim)
	if len(spl) < 2 {
		return nil, new(ErrorIncorrectFormat)
	}
	sig, e := hex.DecodeString(spl[0])
	if e != nil {
		return nil, e
	}
	data := []byte(strings.Join(spl[1:], delim))
	_, e = verifier.Verify(data, sig)
	return data, e
}

//signAndPURB takes in a piece of data, signs it and the purbifies the data given along with it signature.
func signAndPURB(signer *EncryptionInfo, recipients []purbs.Recipient, suite purbs.Suite, data []byte) ([]byte, error) {
	sig, e := signer.Sign(data)
	if e != nil {
		return nil,  e
	}
	hexSign := hex.EncodeToString(sig)
	signedBlob := []byte(hexSign + ";" + string(data))
	params := purbs.NewPublicFixedParameters(signer.RetrievalInfo.SuiteInfos, false)
	p, e := purbs.Encode(signedBlob, recipients, random.New(), params, false)
	if e != nil {
		return nil, e
	}
	return p.ToBytes(), nil
}


func getTCPConn(p *Peer) (net.Conn, error) {
	c, e := net.Dial("tcp", p.Addr.String()+":"+strconv.Itoa(p.Port))
	return c, e
}

func readNBytesFromConnection(c net.Conn, n int) ([]byte, error) {
	reader := bufio.NewReader(c)
	buffer := make([]byte, n)
	_, e := reader.Read(buffer)
	if e != nil {
		return nil, e
	}
	return buffer, nil
}

func checkIfFailed(d []byte) error {
	s := string(d)
	isError := strings.Contains(s, errorIndicator)
	if isError {
		return new(ErrorFailedProtocol)
	}
	return nil
}

