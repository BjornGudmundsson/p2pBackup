package files

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

//Backup keeps track of the relevant
type Backup struct {
	Date   time.Time
	Hash   string
	Size   int32
	CTSize int32
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
	str := strconv.Itoa(size) + " " + hx + " " + now.String() + "\n"
	return AppendToFile(*f, []byte(str))
}
