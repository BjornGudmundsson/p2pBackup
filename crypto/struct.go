package crypto

import (
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/util/key"
)

var cipherSuiteRegistry map[string]func() (KeyPair, error)

func setCipherSuiteRegistry() {
	cipherSuiteRegistry["Ed25519"] = newEd25519KeyPair
}

func init() {
	setCipherSuiteRegistry()
}

//SecretKey is the representation of a
//generic private key
type SecretKey []byte

//PublicKey is the representation of a
//generic public key
type PublicKey []byte

//KeyPair represents a generic public/private
//key pair for any desired ciphersuite
type KeyPair interface {
	Private() (SecretKey, error)
	PublicKey() (PublicKey, error)
}

func newEd25519KeyPair() (KeyPair, error) {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	kp := key.NewKeyPair(suite)
	ed := &Ed25519KeyPair{
		sk: kp.Private,
		pk: kp.Public,
	}
	return ed, nil
}

//Ed25519KeyPair is a wrappper around
//a pair of a scalar and a point on an
//elliptic curve with the edwards curve parameters.
type Ed25519KeyPair struct {
	sk kyber.Scalar
	pk kyber.Point
}

//Private returns the marshalled version of the private key and
//returns an error if, for some reason, it could not be marshalled.
func (ed *Ed25519KeyPair) Private() (SecretKey, error) {
	return ed.sk.MarshalBinary()
}

//PublicKey returns the marshalled version of the public key
//and returns an error if, for some reason, it could not be marshalled.
func (ed *Ed25519KeyPair) PublicKey() (PublicKey, error) {
	return ed.pk.MarshalBinary()
}

//CipherSuiteNotFound is an error thrown if
//the user asked for a ciphersuite that is not registered
type CipherSuiteNotFound struct {
	suite string
}

func (e *CipherSuiteNotFound) Error() string {
	return "Ciphersuite " + e.suite + " is not registered and could not be found"
}

//NewCipherSuiteNotFoundError returns a new ciphersuite not
//found error with the given suite name.
func NewCipherSuiteNotFoundError(suite string) error {
	return &CipherSuiteNotFound{
		suite: suite,
	}
}
