package crypto

import (
	aes2 "crypto/aes"
	"crypto/cipher"
)

func GetKeyStream(nonce, key []byte, size int64) ([]byte, error) {
	aes, e := aes2.NewCipher(key)
	if e != nil {
		return nil, e
	}
	ctr := cipher.NewCTR(aes, nonce)
	stream := make([]byte, size)
	ctr.XORKeyStream(stream, stream)
	return stream, nil
}

func MergeWithStream(nonce, key, oldStream []byte) error {
	aes, e := aes2.NewCipher(key)
	if e != nil {
		return e
	}
	ctr := cipher.NewCTR(aes, nonce)
	ctr.XORKeyStream(oldStream, oldStream)
	return nil
}

func EncryptCTR(nonce, key, data []byte) error {
	block, e := aes2.NewCipher(key)
	if e != nil {
		return e
	}
	ctr := cipher.NewCTR(block, nonce)
	ctr.XORKeyStream(data, data)
	return nil
}
