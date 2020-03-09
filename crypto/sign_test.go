package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEd25519Success(t *testing.T) {
	kp, e := NewKeyPair(EDWARDS)
	assert.Nil(t, e, "There should exist a suite for the edwards curve")
	sk, e := kp.Private()
	assert.Nil(t, e, "should be able to generate a secret key")
	pk, e := kp.PublicKey()
	assert.Nil(t, e, "Should be able to generate a public key")
	assert.NotNil(t, sk, "The private key should be non-nil")
	assert.NotNil(t, pk, "The public key should be non-nil")
	signer, esign := NewSigner(EDWARDS)
	assert.Nil(t, esign, "There should exist a sign function for the edwards curve")
	msg := []byte("Bjorn is cool")
	sig, e := signer(sk, msg)
	assert.Nil(t, e, "Should be able to sign a given message")
	assert.NotNil(t, sig, "The signature should be non nil")
	verifier, e := NewVerifier(EDWARDS)
	assert.Nil(t, e, "There should exist a verifier for the given edwards curve")
	ver := verifier(pk, sig, msg)
	assert.Nil(t, ver, "There should not be an error when verifying the given message")
	garbage := verifier(pk, sig, []byte("Garbage message"))
	assert.NotNil(t, garbage, "There should be an error when verifying the garbage message")
}
