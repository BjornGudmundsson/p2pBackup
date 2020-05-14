package utilities

import (
	"encoding/hex"
	"errors"
	"github.com/BjornGudmundsson/p2pBackup/kyber"
	"github.com/BjornGudmundsson/p2pBackup/purb/purbs"
)

func HexToKey(hx string, suite purbs.Suite) (kyber.Scalar, error) {
	k, e := hex.DecodeString(hx)
	if e != nil {
		return nil, e
	}
	x := suite.Scalar()
	e = x.UnmarshalBinary(k)
	if e != nil {
		return nil, e
	}
	return x, nil
}

func XORBuffers(dst, src []byte) error {
	if len(dst) != len(src) {
		return errors.New("buffer lengths did not match")
	}
	for i, b := range src {
		b2 := dst[i]
		dst[i] = b ^ b2
	}
	return nil
}