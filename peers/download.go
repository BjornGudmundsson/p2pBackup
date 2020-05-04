package peers

import (
	aes2 "crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/BjornGudmundsson/p2pBackup/crypto"
	"github.com/BjornGudmundsson/p2pBackup/files"
	"github.com/BjornGudmundsson/p2pBackup/kyber"
	"github.com/BjornGudmundsson/p2pBackup/kyber/group/edwards25519"
	"github.com/BjornGudmundsson/p2pBackup/kyber/proof/dleq"
	"github.com/BjornGudmundsson/p2pBackup/kyber/util/random"
	"github.com/BjornGudmundsson/p2pBackup/purb/purbs"
	"strconv"
	"strings"
	"time"
)
const DOWNLOAD = "Download"
const wait = time.Second

func getDownloadHandler(bh files.BackupHandler, start, size int64) func(Communicator) error {
	f := func(c Communicator) error {
		data, e := bh.ReadFrom(start, size)
		if e != nil {
			return e
		}
		verified, e := verifyDownload(c, data)
		if e != nil  {
			return e
		}
		if !verified {
			return new(ErrorUnableToVerify)
		}
		encryptedData := encryptData(data)
		e = c.SendMessage(encryptedData)
		if e != nil {
			return e
		}
		return c.CloseChannel()
	}
	return f
}

func RetrieveFromLogs(logs files.LogWriter, enc *EncryptionInfo, container Container) ([]byte, error) {
	log, e := logs.GetLatestLog()
	if e != nil {
		return nil, e
	}
	return RetrieveBackup(log, container, enc)
}

func RetrieveBackup(log files.Log, container Container, enc *EncryptionInfo) ([]byte, error) {
	indexes := []uint64(log.Retrieve())
	size := log.Size()
	for i := 0; i < 5;i++ {
		time.Sleep(wait)//Sleep since it can take some time to get an up to date peer list
		peers := container.GetPeerList()
		for _, index := range indexes {
			//Iterate over all possible indexes since each peer may have a different
			msg := DOWNLOAD + delim + strconv.FormatUint(index, 10) + delim + strconv.FormatUint(size, 10)
			for _, peer := range peers {
				pt, e := getBackupFromPeer(peer, enc, log, msg)
				if e != nil {
					continue
				}
				return pt, nil
			}
		}
	}
	return nil, new(ErrorCouldNotRetrieveBackup)
}

func getBackupFromPeer(peer Peer, enc *EncryptionInfo, log files.Log, msg string) ([]byte, error) {
	c, e := NewCommunicatorFromPeer(peer, enc)
	if e != nil {
		return nil, e
	}
	msgBlob, e := enc.PURBAnon([]byte(msg))
	if e != nil {
		return nil, e
	}
	e = c.SendMessage(msgBlob)
	if e != nil {
		return nil, e
	}
	hasBackup, e := performDownloadChallenge(c, log)
	if e != nil  || !hasBackup{
		return nil, e
	}
	ct, e := c.GetNextMessage()
	if e != nil {
		return nil, e
	}
	blob, e := decryptAndVerifyData(ct, log)
	if e != nil {
		return nil, e
	}
	pt, e := enc.DecodePURBBackup(blob)
	if e != nil {
		return nil, e
	}
	return pt, nil
}

func dataToScalar(suite purbs.Suite, d []byte) (kyber.Scalar, error) {
	digest := sha256.Sum256(d)
	hexDigest := hex.EncodeToString(digest[:])
	x, e := crypto.PrivateKeyFromPassword(hexDigest, suite)
	if e != nil {
		return nil, e
	}
	return x, nil
}

func verifyDownload(c Communicator, d []byte) (bool, error) {
	rng := random.New()
	suite := edwards25519.NewBlakeSHA256Ed25519()
	G, H := suite.Point().Pick(rng), suite.Point().Pick(rng)
	x, e := dataToScalar(suite, d)
	if e != nil {
		return false, e
	}
	msg := []byte(G.String() + seperator + H.String())
	e = c.SendMessage(msg)
	if e != nil {
		return false, e
	}
	resp, e := c.GetNextMessage()
	if e != nil {
		return false, e
	}
	proof, e := unmarshallProof(resp, suite)
	if e != nil {
		return false, e
	}
	Gc, Hc := G.Clone(), H.Clone()
	xG, xH := Gc.Mul(x, G), Hc.Mul(x, H)
	e = proof.Verify(suite, G, H, xG, xH)
	if e != nil {
		return false, nil
	}
	return true, nil
}

func performDownloadChallenge(c Communicator, log files.Log) (bool, error) {
	challenge, e := c.GetNextMessage()
	if e != nil {
		return false, e
	}
	suite := edwards25519.NewBlakeSHA256Ed25519()
	G, H, e := getBasePoints(challenge, suite)
	if e != nil {
		return false, e
	}
	key := log.Key()
	x, e := crypto.PrivateKeyFromPassword(key, suite)
	proof, _, _, e := dleq.NewDLEQProof(suite, G, H, x)
	if e != nil {
		return false, e
	}
	//fmt.Println("Made proof")
	data, e := marshallProof(proof)
	if e != nil {
		return false, e
	}
	//fmt.Println("Marshalled proof")
	e = c.SendMessage(data)
	if e != nil {
		return false, nil
	}
	return true, nil
}

func marshallProof(proof *dleq.Proof) ([]byte, error) {
	c := proof.C
	r := proof.R
	vG := proof.VG
	vH := proof.VH
	data := []byte(c.String() + seperator + r.String() + seperator + vG.String() + seperator + vH.String())
	return data, nil
}

func unmarshallProof(d []byte, suite *edwards25519.SuiteEd25519) (*dleq.Proof, error) {
	fields := strings.Split(string(d), seperator)
	if len(fields) != 4 {
		fmt.Println("Bjo")
		return nil, new(ErrorFailedProtocol)
	}
	mC, e := hex.DecodeString(fields[0])
	if e != nil {
		return nil, e
	}
	mR, e := hex.DecodeString(fields[1])
	if e != nil {
		return nil, e
	}
	mVG, e := hex.DecodeString(fields[2])
	if e != nil {
		return nil, e
	}
	mVH, e := hex.DecodeString(fields[3])
	if e != nil {
		return nil, e
	}
	c, r, vG, vH := suite.Scalar(), suite.Scalar(), suite.Point(), suite.Point()
	e = c.UnmarshalBinary(mC)
	if e != nil {
		return nil, e
	}
	e = r.UnmarshalBinary(mR)
	if e != nil {
		return nil, e
	}
	e = vG.UnmarshalBinary(mVG)
	if e != nil {
		return nil, e
	}
	e = vH.UnmarshalBinary(mVH)
	if e != nil {
		return nil, e
	}
	return &dleq.Proof{
		C: c,
		R: r,
		VH: vH,
		VG: vG,
	}, nil
}

func getBasePoints(d []byte, suite *edwards25519.SuiteEd25519) (G, H kyber.Point, e error) {
	G, H = suite.Point(), suite.Point()
	spl := strings.Split(string(d), seperator)
	if len(spl) != 2 {
		e = new(ErrorIncorrectFormat)
		return
	}
	mG, e := hex.DecodeString(spl[0])
	if e != nil {
		return
	}
	mH, e := hex.DecodeString(spl[1])
	if e != nil {
		return
	}
	e = G.UnmarshalBinary(mG)
	if e != nil {
		return
	}
	e = H.UnmarshalBinary(mH)
	if e != nil {
		return
	}
	return
}

func encryptData(d []byte) []byte {
	digest := sha256.Sum256(d)
	iv := sha256.Sum256(digest[:])
	aes, e := aes2.NewCipher(digest[:][files.KEYLEN:])
	if e != nil {
		fmt.Println("Error", e)
		return nil
	}
	cbc := cipher.NewCBCEncrypter(aes, iv[:][files.KEYLEN:])
	padded := padData(aes.BlockSize(),d)
	dst := make([]byte, len(padded))
	cbc.CryptBlocks(dst, padded)
	return dst
}

func decryptAndVerifyData(d []byte, log files.Log) ([]byte, error) {
	key, size := log.Key(), log.Size()
	k, e := hex.DecodeString(key)
	if e != nil {
		return nil, e
	}
	iv := sha256.Sum256(k)
	aes, e := aes2.NewCipher(k[files.KEYLEN:])
	if e != nil {
		return nil, e
	}
	decrypter := cipher.NewCBCDecrypter(aes, iv[:][files.KEYLEN:])
	dst := make([]byte, len(d))
	decrypter.CryptBlocks(dst, d)
	if uint64(len(d)) < size {
		return nil, new(ErrorFailedProtocol)
	}
	blob := dst[:size]
	digest := sha256.Sum256(blob)
	hxDigest := hex.EncodeToString(digest[:])
	if hxDigest != key {
		return nil, new(ErrorFailedProtocol)
	}
	return blob, nil
}

func padData(blockSize int, d []byte) []byte {
	if len(d) % blockSize == 0 {
		return d
	}
	padded := make([]byte, blockSize)
	random.Bytes(padded, random.New())
	copy(padded, d)
	return padded
}

