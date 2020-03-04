package crypto

type SecretKey interface {
	String() string
	Suite() string
	PublicKey() []byte
}
