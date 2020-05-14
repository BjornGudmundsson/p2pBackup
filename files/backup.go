package files

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/BjornGudmundsson/p2pBackup/crypto"
	"github.com/BjornGudmundsson/p2pBackup/utilities"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

func getFileSize(f *os.File) int64 {
	info, e := f.Stat()
	if e != nil {
		return -1
	}
	return info.Size()
}

//Backup keeps track of the relevant
type Backup struct {
	Date   time.Time
	Hash   string
	Size   int32
	CTSize int32
}

type BackupHandler interface {
	AddBackup(d []byte) int64
	ReadFrom(start, end int64) ([]byte, error)
}

type BackupBuffer struct {
	fn string
	mtx sync.Mutex
	buffer []byte
	wait time.Duration
	key []byte
}

//AddBackup takes in data and backs it up by writing to the file
//or by buffering it up
func (bb *BackupBuffer) AddBackup(d []byte) int64 {
	bb.mtx.Lock()
	defer bb.mtx.Unlock()
	fl, e := bb.getFile()
	if e != nil {
		return -1
	}
	unwrittenSize := len(bb.buffer)
	bb.buffer = append(bb.buffer, d...)
	return fl.Size + int64(unwrittenSize)
}


//ReadFrom reads the bytes of the buffer from the start index and reads
//size amount of bytes.
func (bb *BackupBuffer) ReadFrom(start, size int64) ([]byte, error) {
	bb.mtx.Lock()
	defer bb.mtx.Unlock()
	fn := bb.fn
	f, e := os.OpenFile(fn, os.O_RDONLY, os.ModeAppend)
	if e != nil {
		return nil, e
	}
	fileSize := getFileSize(f)
	keyStream, e := bb.getKeyStream(fileSize)
	if e != nil {
		return nil, e
	}
	buffer := make([]byte, size)
	_, e = f.ReadAt(buffer, start)
	if e != nil {
		return nil, e
	}
	e = encryptPartKeyStream(keyStream, buffer, start)
	if e != nil {
		return nil, e
	}
	return buffer, nil
}

//writeToFile runs while the system is running and periodically adds new backups
func (bb *BackupBuffer) writeToFile() {
	for {
		time.Sleep(bb.wait)
		bb.mtx.Lock()
		buffer := make([]byte, len(bb.buffer))
		copy(buffer, bb.buffer)
		if len(buffer) == 0 || buffer == nil {
			bb.mtx.Unlock()
			continue
		}
		fl, e := bb.getFile()
		if e != nil {
			bb.mtx.Unlock()
			continue
		}
		f, e := os.OpenFile(bb.fn, os.O_RDWR, os.ModeAppend)
		if e != nil {
			bb.mtx.Unlock()
			continue
		}
		size := getFileSize(f)
		start := fl.Size
		totalWritten := start + int64(len(buffer))
		bb.updateMetadata(totalWritten, f)
		var keyStream []byte
		if size > totalWritten + 2 * KEYLEN {
			keyStream, e = bb.getKeyStream(size)
		} else {
			keyStream, e = bb.getKeyStream(totalWritten)
		}
		if e != nil {
			bb.mtx.Unlock()
			continue
		}
		if e = encryptPartKeyStream(keyStream, buffer, start) ;e != nil {
			bb.mtx.Unlock()
			continue
		}
		if e = AppendToFile(*fl, buffer);e != nil {
			bb.mtx.Unlock()
			continue
		}
		bb.buffer = nil
		bb.mtx.Unlock()
	}
}

func encryptPartKeyStream(keyStream, data []byte, start int64) error {
	l, ld := int64(len(keyStream)), int64(len(data))
	if l < ld{
		return errors.New("key stream too short")
	}
	startStream := keyStream[start:]
	for i, b := range data {
		k := startStream[i]
		data[i] = b ^ k
	}
	return nil
}

func (bb *BackupBuffer) getKeyStream(size int64) ([]byte, error) {
	f, e := os.OpenFile(bb.fn, os.O_RDWR, os.ModeAppend)
	if e != nil {
		return nil, e
	}
	fmt.Println("Size: ", size)
	nonce := make([]byte, KEYLEN)
	_, e = f.ReadAt(nonce, 0)
	if e != nil {
		return nil, e
	}
	keyStream, e := crypto.GetKeyStream(nonce, bb.key, size-KEYLEN)
	if e != nil {
		return nil, e
	}
	return keyStream, nil
}

func (bb *BackupBuffer) reEncryptFile(keyStream []byte, f *os.File) error {
	info , e := f.Stat()
	if e != nil {
		return e
	}
	size := info.Size()
	key := bb.key
	newNonce := make([]byte, KEYLEN)
	rand.Read(newNonce)
	crypto.MergeWithStream(newNonce, key, keyStream)
	_, e = f.WriteAt(newNonce, 0)
	if e != nil {
		return nil
	}
	data := make([]byte, size - KEYLEN)
	_, e = f.ReadAt(data, KEYLEN)
	if e != nil {
		return e
	}
	e = utilities.XORBuffers(data, keyStream)
	_, e = f.WriteAt(data, KEYLEN)
	return e
}

func (bb *BackupBuffer) updateMetadata(size int64, f *os.File) error {
	num := make([]byte, KEYLEN)
	u := uint64(size)
	binary.LittleEndian.PutUint64(num[:KEYLEN / 2], u)
	binary.LittleEndian.PutUint64(num[KEYLEN / 2:], u)
	_, e := f.WriteAt(num, KEYLEN)
	if e != nil {
		return e
	}
	return nil
}

func NewBackupBuffer(fn, pw string) BackupHandler {
	bb := &BackupBuffer{
		fn: fn,
		buffer: make([]byte, 0),
		wait: time.Second,
		key: crypto.SymmetricKeyFromPassword(pw),
	}
	//Keep a loop running that periodically writes to file
	go bb.writeToFile()
	return bb
}

func (bb *BackupBuffer) getFile() (*File, error) {
	fn := bb.fn
	f, e := os.OpenFile(fn, os.O_RDWR, os.ModeAppend)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil {
		return nil, e
	}
	size := info.Size()
	if size < 2 * KEYLEN {
		return nil, errors.New("not able to retrieve file metadata")
	}
	metadata := make([]byte, 2 * KEYLEN)
	_, e = f.Read(metadata)
	if e != nil {
		return nil, e
	}
	fl, e := GetFile(fn)
	if e != nil {
		return fl, e
	}
	_, written := metadata[:KEYLEN], metadata[KEYLEN:2 * KEYLEN]
	left, right := written[: KEYLEN / 2], written[KEYLEN / 2: KEYLEN]
	if string(left) != string(right) {
		fl.Size = 2 * KEYLEN
	} else {
		writtenBytes := int64(binary.LittleEndian.Uint64(left))
		fl.Size = 2 * KEYLEN + writtenBytes
	}
	return fl, nil
}


type LogWriter interface {
	CheckIfBackedUp(d []byte) (bool, error)
	Log(l Log) error
	GetLogs() ([]Log, error)
	GetLatestLog() (Log, error)
	NewLog(d []byte, loc Locations, key []byte) Log
}

type LogHandler struct {
	fn string//Name of the file where the logs are kep
	pw string//If the logfile is encrypted then this is not ""
	mtx sync.Mutex//Adding this in case this becomes concurrent later
}

func (lh *LogHandler) CheckIfBackedUp(d []byte) (bool, error) {
	digest := sha256.Sum256(d)
	hxDigest := hex.EncodeToString(digest[:])
	logs, e := lh.GetLogs()
	noLogs := new(ErrorNoLogs)
	if e != nil && !compareErrors(e, noLogs) {
		return false, e
	}
	for _, log := range logs {
		if hxDigest == log.Digest() {
			return true, nil
		}
	}
	return false, nil
}

func (lh *LogHandler) Log(l Log) error {
	return AddBackupLog(l, lh.fn, lh.pw)
}

func (lh *LogHandler) GetLogs() ([]Log, error) {
	ct, e := ioutil.ReadFile(lh.fn)
	if e != nil {
		return nil, e
	}
	k := pwToKey(lh.pw)
	key := k[:KEYLEN]
	iv := make([]byte, KEYLEN)
	block, e := aes.NewCipher(key)
	if e != nil {
		return nil, e
	}
	dec := cipher.NewCBCDecrypter(block, iv)
	pt := make([]byte, len(ct))
	dec.CryptBlocks(pt, ct)
	logs := NewEmptyLogEntry().FindLogs(pt)
	if len(logs) == 0 {
		return nil, new(ErrorNoLogs)
	}
	return logs, nil
}

func (lh *LogHandler) GetLatestLog() (Log, error) {
	logs, e := lh.GetLogs()
	if e != nil {
		return nil, e
	}
	l := len(logs)
	if l == 0 {
		return nil, new(ErrorNoLogs)
	}
	return logs[l - 1], nil
}

func (lh *LogHandler) NewLog(d []byte, loc Locations, key []byte) Log {
	digest := sha256.Sum256(d)
	digestKey := sha256.Sum256(key)
	log := LogEntry{
		indexes:loc,
		hash: hex.EncodeToString(digest[:]),
		sizeCT:uint64(len(key)),
		date:time.Now(),
		key: hex.EncodeToString(digestKey[:]),
	}
	return log
}

func NewEncryptedLogWriter(fn, pw string) (LogWriter, error) {
	return &LogHandler{
		fn: fn,
		pw: pw,
	}, nil
}

func NewEmptyLogEntry() Log {
	return LogEntry{}
}