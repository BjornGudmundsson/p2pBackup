package crypto

import (
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/sign/schnorr"
)

//Signature is an alias for a slice of bytes to represent a generic cryptographic signature
type Signature []byte

//Message is an alias for a slice of bytes to represent a message
type Message []byte

//Verifier is a generic way to verify a signature given any public key, signature and a message.
type Verifier func(PublicKey, Signature, Message) error

//Signer is an alias for a generic way to sign a message given any desired ciphersuite,
type Signer func(SecretKey, Message) (Signature, error)

//Ed25519Verification verifies a given message and signature pair for a given public key
//and returns an error if it could not be verified else it returns nil
func Ed25519Verification(pk PublicKey, sig Signature, msg Message) error {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	p := suite.Point()
	e := p.UnmarshalBinary(pk)
	if e != nil {
		return e
	}
	e = schnorr.Verify(suite, p, msg, sig)
	return e
}

//Ed25519Sign takes in a secret key and a message and signs it.
//It returns an error if it was unable to sign the message
func Ed25519Sign(sk SecretKey, msg Message) (Signature, error) {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	k := suite.Scalar()
	e := k.UnmarshalBinary(sk)
	if e != nil {
		return nil, e
	}
	return schnorr.Sign(suite, k, msg)
}
