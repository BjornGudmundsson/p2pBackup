package peers

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/BjornGudmundsson/p2pBackup/files"
	"net"
	"strconv"
)

func getDownloadHandler(bh files.BackupHandler, start, size int64) tcpHandler {
	f := tcpHandler(func(c net.Conn) error {
		data, e := bh.ReadFrom(start, size)
		if e != nil {
			return e
		}
		_, e = fmt.Fprintf(c, hex.EncodeToString(data))
		if e != nil {
			return e
		}
		return c.Close()
	})
	return f
}

func RetrieveBackup(p *TCPPeer,info BackupInfo, encryptionInfo *EncryptionInfo) ([]byte, error) {
	start, size := info.StartIndex, info.Size
	msg := strconv.FormatInt(start, 10) + ";" + strconv.FormatInt(size, 10)
	c, e := getTCPConn(p)
	if e != nil {
		return nil, e
	}
	_, e = fmt.Fprintf(c, msg)
	success, e := performDownloadChallenge(c, info, encryptionInfo)
	if e != nil {
		return nil, e
	}
	if !success {
		return nil, new(ErrorUnableToProveStorage)
	}
	ct, e := getBackup(c, info, encryptionInfo)
	if e != nil {
		return nil, e
	}
	blob, e := decryptDownload(ct, info)
	return encryptionInfo.DecodePURBBackup(blob)
}

func performDownloadChallenge(c net.Conn, info BackupInfo, encryptionInfo *EncryptionInfo) (bool, error) {
	return true, nil
}

func decryptDownload(d []byte, info BackupInfo) ([]byte, error) {
	//TODO: Make it actually decrypt the data once the encryption has been made
	return d, nil
}

func getBackup(c net.Conn, info BackupInfo, encryptionInfo *EncryptionInfo) ([]byte, error) {
	reader := bufio.NewReader(c)
	s, e := reader.ReadString('\n')
	if e != nil {
		return nil, e
	}
	return []byte(s), nil
}


