package suites

import (
	"github.com/BjornGudmundsson/p2pBackup/kyber/group/edwards25519"
)

func init() {
	register(edwards25519.NewBlakeSHA256Ed25519())
}
