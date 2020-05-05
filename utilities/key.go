package utilities

import (
	"encoding/hex"
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