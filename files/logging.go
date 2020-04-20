package files

import (
	aes2 "crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"golang.org/x/crypto/pbkdf2"
)

const (
	ITERATIONS = 5
	KEYLEN = 16
)

func padData(d []byte) []byte {
	l := len(d)
	m := l % KEYLEN
	if m == 0 {
		return d
	}
	padding := KEYLEN - m
	padded := make([]byte, l + padding)
	copy(padded, d)
	return padded
}

func pwToKey(pw string) []byte {
	salt := sha256.New().Sum([]byte(pw))
	key := pbkdf2.Key([]byte(pw), salt, ITERATIONS, KEYLEN, sha256.New)
	return key
}
func AddBackupLog(log Log, logFile string, pw string) error {
	salt := sha256.New().Sum([]byte(pw))
	key := pbkdf2.Key([]byte(pw), salt, ITERATIONS, KEYLEN, sha256.New)
	entry := log.MarshallToString()
	padded := padData([]byte(entry))
	iv, e := GetLastNBytes(logFile, KEYLEN)
	if e != nil {
		return e
	}
	if len(iv) != KEYLEN {
		return errors.New("could not construct an IV of sufficient length")
	}
	aes, e := aes2.NewCipher(key)
	if e != nil {
		return e
	}
	encryptedLog := make([]byte, len(padded))
	enc := cipher.NewCBCEncrypter(aes, iv)
	enc.CryptBlocks(encryptedLog, padded)
	f, e := GetFile(logFile)
	if e != nil {
		return e
	}
	return AppendToFile(*f, encryptedLog)
}
