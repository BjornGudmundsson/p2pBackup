package crypto

//NewKeyPair takes in a protocol(i.e. different curve, latticebased etc)
//and returns a public/private key pair.
func NewKeyPair(suite string) (KeyPair, error) {
	if gen, ok := cipherSuiteRegistry[suite]; ok {
		return gen()
	}
	return nil, NewCipherSuiteNotFoundError(suite)
}

//NewSigner returns a new signer for a desired ciphersuite
//and returns an error if that suite could not be found.
func NewSigner(suite string) (Signer, error) {
	if sign, ok := signatureRegistry[suite]; ok {
		return sign, nil
	}
	return nil, NewCipherSuiteNotFoundError(suite)
}

//NewVerifier takes in a suite and returns the verifying function
//for that suite and returns an error if that suite is not registered.
func NewVerifier(suite string) (Verifier, error) {
	if ver, ok := verifierRegistry[suite]; ok {
		return ver, nil
	}
	return nil, NewCipherSuiteNotFoundError(suite)
}
