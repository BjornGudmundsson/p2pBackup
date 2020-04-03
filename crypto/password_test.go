	package crypto

import (
	"github.com/BjornGudmundsson/p2pBackup/kyber/util/random"
	"github.com/BjornGudmundsson/p2pBackup/purb"
	"github.com/BjornGudmundsson/p2pBackup/purb/purbs"
	"github.com/stretchr/testify/assert"
	"github.com/BjornGudmundsson/p2pBackup/kyber/group/curve25519"
	"testing"
)

func TestPrivateKeyFromPassword(t *testing.T) {
	pw := "Bjorn"
	suite, e := purb.GetSuite(curve25519.NewBlakeSHA256Curve25519(true).String())
	assert.Nil(t, e, "Should be able to get Suite")
	sk, e := PrivateKeyFromPassword(pw, suite)
	assert.Nil(t, e, "Should be able to generate a scalar")
	assert.NotNil(t, sk, "Private key should be non-nil")
	p := suite.Point().Base()
	pk := p.Mul(sk, p)
	assert.NotNil(t, pk, "Public key should be non-nil")
}

func TestPasswordToPURB(t *testing.T) {
	toEncrypt := "I sell seashells by the seashore"
	pw := "Bjorn"
	suite, e := purb.GetSuite(curve25519.NewBlakeSHA256Curve25519(true).String())
	assert.Nil(t, e, "Should be able to generate a suite")
	assert.NotNil(t, suite, "Suite should be non-nil")
	sk, e := PrivateKeyFromPassword(pw, suite)
	assert.Nil(t, e, "Should be able to get private key from scalar")
	assert.NotNil(t, sk, "Private key should be non-nil")
	p := suite.Point().Base()
	pk := p.Mul(sk, p)
	m, e := pk.MarshalBinary()
	assert.Nil(t, e, "Should be able to marshal the public key")
	recipient, e := purb.NewRecipient(m, suite)
	assert.Nil(t, e, "should be able to generate a recipient")
	suiteInfoMap := purb.GetSuiteInfos("")
	params := purbs.NewPublicFixedParameters(suiteInfoMap, false)
	recipients := []purbs.Recipient{recipient}
	data, e := purbs.Encode([]byte(toEncrypt), recipients, random.New(), params, false)
	assert.Nil(t, e, "should be able to encode the PURB")
	msk, e := sk.MarshalBinary()
	assert.Nil(t, e, "Should be able to unmarshall binary key")
	privRecipient, e := purb.NewPrivateRecipient(msk, suite)
	assert.Nil(t, e, "should be able to get a private recipient")
	blob := data.ToBytes()
	success, decrypted, e := purbs.Decode(blob, &privRecipient,params, false)
	assert.Nil(t, e , "Should be able to decode the PURB")
	assert.True(t, success, "Should be able to decode the PURB")
	assert.Equal(t, string(decrypted), toEncrypt, "The decrypted data should be the same as the unencrypted data")
}
