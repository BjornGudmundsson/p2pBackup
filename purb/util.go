package purb

import (
	"errors"
	"github.com/BjornGudmundsson/p2pBackup/kyber"
	"github.com/BjornGudmundsson/p2pBackup/purb/purbs"
	"github.com/BjornGudmundsson/p2pBackup/kyber/group/curve25519"
)

func GetSuite(suite string) (purbs.Suite, error) {
	if suite == curve25519.NewBlakeSHA256Curve25519(true).String() {
		return curve25519.NewBlakeSHA256Curve25519(true), nil 
	}
	return nil, errors.New("Given suite could not be found")
}

func GetSuiteInfos(fn string) (purbs.SuiteInfoMap) {
	if fn == "" {
		return getDummySuiteInfo()
	}
	//TODO:Implement reading a TOML file to get the suiteInfoMap
	return nil
}


func getDummySuiteInfo() purbs.SuiteInfoMap {
	info := make(purbs.SuiteInfoMap)
	cornerstoneLength := 32 // defined by Curve 25519
	entryPointLength := 16 + 4 + 4 + 16 // 16-byte symmetric key + 2 * 4-byte offset positions + 16-byte authentication tag
	info[curve25519.NewBlakeSHA256Curve25519(true).String()] = &purbs.SuiteInfo{
		AllowedPositions: []int{12 + 0*cornerstoneLength, 12 + 1*cornerstoneLength, 12 + 3*cornerstoneLength, 12 + 4*cornerstoneLength},
		CornerstoneLength: cornerstoneLength, EntryPointLength: entryPointLength}
	return info
}

type KeyInfo struct {
	Suite purbs.Suite
	SuiteInfos purbs.SuiteInfoMap
	PrivateKey kyber.Scalar
	PublicKey kyber.Point
}

func NewKeyInfo(sk []byte, suite purbs.Suite, suiteFile string) (*KeyInfo, error) {
	s := suite.Scalar()
	e := s.UnmarshalBinary(sk)
	if e != nil {
		return nil, e
	}
	suiteInfo := GetSuiteInfos(suiteFile)
	p := suite.Point().Base()
	p = p.Mul(s, p)
	return &KeyInfo{
		PrivateKey:s,
		PublicKey:p,
		SuiteInfos: suiteInfo,
		Suite:suite,
	}, nil
}

func (info *KeyInfo) String() string {
	sk := "SK: " + info.PrivateKey.String() + "\n"
	pk := "PK: " + info.PublicKey.String() + "\n"
	s := "Suite: " + info.Suite.String()
	return sk + pk + s
}


func NewRecipient(pk []byte, suite purbs.Suite) (purbs.Recipient, error) {
	p := suite.Point()
	e := p.UnmarshalBinary(pk)
	rec := purbs.Recipient{}
	if e != nil {
		return rec, e
	}
	rec.PublicKey = p
	rec.Suite = suite
	rec.SuiteName = suite.String()
	return rec, nil
}

func NewPrivateRecipient(sk []byte, suite purbs.Suite) (purbs.Recipient, error) {
	s := suite.Scalar()
	e := s.UnmarshalBinary(sk)
	rec := purbs.Recipient{}
	if e != nil {
		return rec, e
	}
	p := suite.Point().Base()
	p = p.Mul(s, p)
	rec.PrivateKey = s
	rec.PublicKey = p
	rec.Suite = suite
	rec.SuiteName = suite.String()
	return rec, nil
}