package peers

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/BjornGudmundsson/p2pBackup/files"
	"net"
	"strconv"
	"time"
)
const DOWNLOAD = "Download"
const wait = time.Second

func getDownloadHandler(bh files.BackupHandler, start, size int64) tcpHandler {
	f := tcpHandler(func(c net.Conn) error {
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
		_, e = fmt.Fprintf(c, hex.EncodeToString(encryptedData) + "\n")
		if e != nil {
			return e
		}
		return c.Close()
	})
	return f
}

func RetrieveFromLogs(logs files.LogWriter, enc *EncryptionInfo, container Container) ([]byte, error) {
	log, e := logs.GetLatestLog()
	fmt.Println(log)
	if e != nil {
		return nil, e
	}
	return RetrieveBackup(log, container, enc)
}

func RetrieveBackup(log files.Log, container Container, enc *EncryptionInfo) ([]byte, error) {
	indexes := []uint64(log.Retrieve())
	size := log.Size()
	for i := 0; i < 5;i++ {
		fmt.Println("Trying to retrieve from all peers")
		time.Sleep(wait)//Sleep since it can take some time to get an up to date peer list
		peers := container.GetPeerList()
		for _, index := range indexes {
			//Iterate over all possible indexes since each peer may have a different
			msg := DOWNLOAD + delim + strconv.FormatUint(index, 10) + delim + strconv.FormatUint(size, 10) + "\n"
			for _, peer := range peers {
				c, e := getTCPConn(peer)
				if e != nil {
					continue
				}
				reader := bufio.NewReader(c)
				fmt.Fprintf(c, msg)
				hasBackup, e := performDownloadChallenge(c, log)
				if e != nil  || !hasBackup{
					continue
				}
				ct, e := reader.ReadString('\n')
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

func verifyDownload(c net.Conn, d []byte) (bool, error) {
	return true, nil//TODO: Actually perform the ZKP challenge thingy.
}

func performDownloadChallenge(c net.Conn, log files.Log) (bool, error) {
	return true, nil
}

func encryptData(d []byte) []byte {
	return d//Todo: Encrypt the data using
}

func decryptAndVerifyData(d []byte, log files.Log) ([]byte, error) {
	l := len(d)
	if l <= 1 {
		return nil, errors.New("Empty string")
	}
	ct := d[:l - 1]
	return hex.DecodeString(string(ct))
}


