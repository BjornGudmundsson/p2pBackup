package crypto

import (
	"crypto/sha256"
	"encoding/binary"
	"github.com/BjornGudmundsson/p2pBackup/kyber"
	"github.com/BjornGudmundsson/p2pBackup/purb/purbs"
	"golang.org/x/crypto/pbkdf2"
)
const SALT = "SALT"
const KEYLEN = 16
const ITERATIONS = 5
//PrivateKeyFromPassword takes in a password and a given ciphersuite and returns
//a scalar(a privatekey) derived from that password.
func PrivateKeyFromPassword(pw string,suite purbs.Suite) (kyber.Scalar, error) {
	salt := sha256.New().Sum([]byte(pw))
	key := pbkdf2.Key([]byte(pw), salt, ITERATIONS, KEYLEN, sha256.New)
	v := binary.BigEndian.Uint64(key)
	sk := suite.Scalar()
	sk.SetInt64(int64(v))
	return sk, nil
}

func SymmetricKeyFromPassword(pw string) []byte {
	salt := sha256.Sum256([]byte(pw))
	key := pbkdf2.Key([]byte(pw), salt[:], ITERATIONS, KEYLEN, sha256.New)
	return key
}