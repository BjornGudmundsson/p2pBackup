package peers

import (
	"fmt"
	"github.com/BjornGudmundsson/p2pBackup/files"
	"strconv"
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
				c, e := NewCommunicatorFromPeer(peer, enc)
				if e != nil {
					continue
				}
				e = c.SendMessage([]byte(msg))
				if e != nil {
					return nil, e
				}
				hasBackup, e := performDownloadChallenge(c, log)
				if e != nil  || !hasBackup{
					continue
				}
				ct, e := c.GetNextMessage()
				if e != nil {
					continue
				}
				blob, e := decryptAndVerifyData([]byte(ct), log)
				if e != nil {
					fmt.Println(e)
					continue
				}
				pt, e := enc.DecodePURBBackup(blob)
				if e != nil {
					continue
				}
				return pt, nil
			}
		}
	}
	return nil, new(ErrorCouldNotRetrieveBackup)
}

func verifyDownload(c Communicator, d []byte) (bool, error) {
	return true, nil//TODO: Actually perform the ZKP challenge thingy.
}

func performDownloadChallenge(c Communicator, log files.Log) (bool, error) {
	return true, nil
}

func encryptData(d []byte) []byte {
	return d//Todo: Encrypt the data using
}

func decryptAndVerifyData(d []byte, log files.Log) ([]byte, error) {
	return d, nil
}


