package files

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"strconv"
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

//AddBackup takes in the data that is being backed up
//and the name of the file where the backup log will be kept.
func AddBackup(d []byte, fn string) error {
	f, e := GetFile(fn)
	if e != nil {
		return e
	}
	h := sha256.Sum256(d)
	hx := hex.EncodeToString(h[:])
	now := time.Now()
	size := len(d)
	str := strconv.Itoa(size) + ";" + hx + ";" + now.String() + "\n"
	return AppendToFile(*f, []byte(str))
}

type BackupBuffer struct {
	fn string
	//TODO: Add some sane data structure to buffer the backups
}

func (bb *BackupBuffer) AddBackup(d []byte) int64 {
	fn := bb.fn
	fl, e := GetFile(fn)
	if e != nil {
		return -1
	}
	e = AppendToFile(*fl, d)
	if e != nil {
		return -1
	}
	return fl.Size
}

func (bb *BackupBuffer) ReadFrom(start, size int64) ([]byte, error) {
	fn := bb.fn
	f, e := os.OpenFile(fn, os.O_RDONLY, os.ModeAppend)
	if e != nil {
		return nil, e
	}
	buffer := make([]byte, size)
	_, e = f.ReadAt(buffer, start)
	return buffer, e
}

func NewBackupBuffer(fn string) BackupHandler {
	return &BackupBuffer{
		fn: fn,
	}
}


type LogWriter interface {
	CheckIfBackedUp(d []byte) (bool, error)
	Log(l LogEntry) error
	GetLogs() ([]LogEntry, error)
	GetLatestLog() (LogEntry, error)
	NewLog(d []byte, indexes []int) LogEntry
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
	if e != nil {
		return false, e
	}
	for _, log := range logs {
		if hxDigest == log.Hash {
			return true, nil
		}
	}
	return false, nil
}

func (lh *LogHandler) Log(l LogEntry) error {
	return AddBackupLog(l, lh.fn, lh.pw)
}

func (lh *LogHandler) GetLogs() ([]LogEntry, error) {
	ct, e := ioutil.ReadFile(lh.fn)
	if e != nil {
		return nil, e
	}
	k := sha256.Sum256([]byte(lh.pw))

	return nil, nil
}

func (lh *LogHandler) GetLatestLog() (LogEntry, error) {
	return LogEntry{}, nil
}

func (lh *LogHandler) NewLog(d []byte, indexes []int) LogEntry {
	return LogEntry{}
}

func NewEncryptedLogWriter(fn, pw string) (LogWriter, error) {
	return nil, nil
}