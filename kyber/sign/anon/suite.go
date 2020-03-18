package anon

import (
	kyber "github.com/BjornGudmundsson/p2pBackup/kyber"
)

// Suite represents the set of functionalities needed by the package anon.
type Suite interface {
	kyber.Group
	kyber.Encoding
	kyber.XOFFactory
	kyber.Random
}
