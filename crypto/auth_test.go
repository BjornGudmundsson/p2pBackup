package crypto

import (
	"encoding/hex"
	"github.com/BjornGudmundsson/p2pBackup/kyber/group/curve25519"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnonAuthenticator_Sign(t *testing.T) {
	fn := "set.txt"
	suite := curve25519.NewBlakeSHA256Curve25519(true)
	key, e := hex.DecodeString("30944f2b07b37b792b4484a3b191447936bfa598bb84770f0da1a5019df5eee4")
	assert.Nil(t, e, "Hex of key should be of right length")
	msg := []byte("Bjorn")
	scalar := suite.Scalar()
	e = scalar.UnmarshalBinary(key)
	assert.Nil(t, e, "Should be a valid key")
	auth, e := NewAnonAuthenticator(suite, fn)
	assert.Nil(t, e, "Should be able to get a new authenticator")
	sig, e := auth.Sign(scalar, msg, nil)
	assert.Nil(t, e, "Should be able to sign")
	assert.NotNil(t, sig, "Sig should be non-empty")
	v, e := auth.Verify(msg, sig, nil)
	assert.Nil(t, e, "Should be a valid signature")
	assert.Equal(t, len(v), 0, "Should be no linkage tag")
}
