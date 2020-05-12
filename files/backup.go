package files

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

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
}

//AddBackup takes in data and backs it up by writing to the file
//or by buffering it up
func (bb *BackupBuffer) AddBackup(d []byte) int64 {
	bb.mtx.Lock()
	defer bb.mtx.Unlock()
	fn := bb.fn
	fl, e := GetFile(fn)
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
	buffer := make([]byte, size)
	_, e = f.ReadAt(buffer, start)
	return buffer, e
}

//writeToFile runs while the system is running and periodically adds new backups
func (bb *BackupBuffer) writeToFile() {
	fn := bb.fn
	for {
		//Sleep to allow different backups to 'buffer up'
		time.Sleep(bb.wait)
		//Locking the file while adding the data
		bb.mtx.Lock()
		//Write the current buffer to the file
		buffer := bb.buffer
		if len(buffer) == 0 || buffer == nil {
			bb.mtx.Unlock()
			continue
		}
		fl, e := GetFile(fn)
		if e != nil {
			bb.mtx.Unlock()
			continue
		}
		e = AppendToFile(*fl, buffer)
		if e != nil {
			bb.mtx.Unlock()
			continue
		}
		bb.buffer = nil
		bb.mtx.Unlock()
	}
}

func NewBackupBuffer(fn string) BackupHandler {
	bb := &BackupBuffer{
		fn: fn,
		buffer: make([]byte, 0),
		wait: time.Second,
	}
	//Keep a loop running that periodically writes to file
	go bb.writeToFile()
	return bb
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