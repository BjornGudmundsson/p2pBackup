package crypto

import (
	"encoding/hex"
	"github.com/BjornGudmundsson/p2pBackup/kyber"
	"github.com/BjornGudmundsson/p2pBackup/kyber/sign/anon"
	"io/ioutil"
	"strings"
)

type NoMatchingPublicKey struct{
}

func (e NoMatchingPublicKey) Error() string {
	return "There was not matching public key in the anonymity set"
}

func NewNoMatchingPublicKey() NoMatchingPublicKey {
	return NoMatchingPublicKey{}
}


type Authenticator interface {
	Sign(scalar kyber.Scalar, msg []byte, l []byte) ([]byte, error)
	Verify(msg []byte, sig []byte, link []byte) ([]byte, error)
	GetAnonSet() anon.Set
	GetSuite() anon.Suite
}

type AnonAuthenticator struct {
	suite anon.Suite
	set anon.Set
}

func (a *AnonAuthenticator) GetSuite() anon.Suite {
	return a.suite
}

func (a *AnonAuthenticator) Sign(scalar kyber.Scalar, msg []byte, l []byte) ([]byte, error) {
	mine := -1
	suite := a.suite
	p := suite.Point().Base()
	p = p.Mul(scalar, p)
	for i, pk := range a.set {
		if pk.Equal(p) {
			mine = i
			break
		}
	}
	if mine == -1 {
		return nil, NewNoMatchingPublicKey()
	}
	return anon.Sign(a.suite, msg, a.set, l, mine, scalar), nil
}

func (a *AnonAuthenticator) Verify(msg, sig, link []byte) ([]byte, error) {
	return anon.Verify(a.suite, msg, a.set, link, sig)
}

func (a *AnonAuthenticator) GetAnonSet() anon.Set {
	return a.set
}

func NewAnonAuthenticator(suite anon.Suite, fn string) (Authenticator, error) {
	set, e := GetAnonymitySet(fn, suite)
	if e != nil {
		return nil, e
	}
	return &AnonAuthenticator{
		suite: suite,
		set:   set,
	}, nil
}

//GetAnonymitySet takes in a filename and a ciphersuite and returns
//a set of public keys(kyber points)
func GetAnonymitySet(fn string, suite anon.Suite) (anon.Set, error) {
	d, e := ioutil.ReadFile(fn)
	if e != nil {
		return nil, e
	}
	list := strings.Split(string(d), "\n")
	set := make(anon.Set, 0)
	for _, hx := range list {
		if len(hx) == 0 {
			break
		}
		b, e := hex.DecodeString(hx)
		if e != nil {
			return nil, e
		}
		p := suite.Point()
		e = p.UnmarshalBinary(b)
		if e != nil {
			return nil, e
		}
		set = append(set, p)
	}
	return set, nil
}
