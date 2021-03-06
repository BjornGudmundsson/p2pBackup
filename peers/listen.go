package peers

import (
	"errors"
	"fmt"
	"github.com/BjornGudmundsson/p2pBackup/files"
	"github.com/BjornGudmundsson/p2pBackup/kyber/util/key"
	"github.com/BjornGudmundsson/p2pBackup/purb"
	"github.com/BjornGudmundsson/p2pBackup/purb/purbs"
	"net"
	"strconv"
	"strings"
)


//Shorthand for a function handle a connection
type communicationHandler func(Communicator) error

const localhost = "127.0.0.1"

//ListenUDP sets up a UDP server
//on the given port.
func ListenUDP(port string, container Container, enc *EncryptionInfo) error {
	serverAddr, e := net.ResolveUDPAddr("udp", port)
	if e != nil {
		return e
	}
	conn, e := net.ListenUDP("udp", serverAddr)
	if e != nil {
		return e
	}
	defer conn.Close()
	buf := make([]byte, 4096)
	for {
		n, addr, e := conn.ReadFromUDP(buf)
		fmt.Println("Got packddet from: ", addr.String())
		fmt.Println("Message: ", string(buf[:n]))
		if e != nil {
			fmt.Println(e.Error())
		}
	}
	return nil
}

//ListenTCP starts a new tcp server on the given port
func ListenTCP(port string, encInfo *EncryptionInfo, backupHandler files.BackupHandler) error {
	l, e := net.Listen("tcp4", port)
	if e != nil {
		return e
	}
	defer l.Close()
	handler := createHandler(encInfo, backupHandler)
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			comm := NewTCPCommunicatorFromConn(c, encInfo)
			go handler(comm)
		}
	}
	return nil
}

func handleFailure(c Communicator, e error, m purbs.SuiteInfoMap) {
	if c != nil {
		c.SendMessage([]byte(e.Error()))
		c.CloseChannel()
	}
}

func isDownload(msg string) (bool, error) {
	if strings.Contains(msg, DOWNLOAD) {
		return true, nil
	}
	return false, nil
}

func getUploadHandler(suite purbs.Suite,encInfo *EncryptionInfo, backupHandler files.BackupHandler,b []byte) func(c Communicator) error {
	f := func(c Communicator) error {
		freshPair := key.NewKeyPair(suite)
		publicBytes, e := freshPair.Public.MarshalBinary()
		if e != nil {
			return nil
		}
		sig, e := signPublicKey(publicBytes, encInfo, b)
		if e != nil {
			return e
		}
		e = c.SendMessage(sig)
		if e != nil {
			return e
		}
		blob, e := c.GetNextMessage()
		if e != nil {
			return e
		}
		data, e := verifyPURB(freshPair.Private, suite, blob, encInfo)
		if e != nil {
			return e
		}
		ind := backupHandler.AddBackup(data)
		if ind == -1 {
			return new(ErrorCouldNotAppend)
		}
		msg := "Ok " + strconv.FormatInt(ind, 10)
		e = c.SendMessage([]byte(msg))
		return c.CloseChannel()
	}
	return f
}

func firstReply(c Communicator,encInfo *EncryptionInfo, backupHandler files.BackupHandler) (communicationHandler, error) {
	s, e := c.GetNextMessage()
	if e != nil {
		return nil, e
	}
	normal, e := GetNormalHandler(s, encInfo, backupHandler)
	if e == nil {
		return normal, nil
	}
	msg, e := encInfo.DecodePURBAnon(s)
	if e != nil {
		return nil, e
	}
	download, e := isDownload(string(msg))
	if e != nil {
		return nil, e
	}
	if download {
		start, size, e := getIndexes(string(msg))
		if e != nil {
			return nil, e
		}
		downloadHandler := getDownloadHandler(backupHandler, start, size)
		return downloadHandler, nil
	}
	suite, e := purb.GetSuite(string(msg))
	if e != nil {
		return nil, e
	}
	handler := getUploadHandler(suite, encInfo, backupHandler, s)
	return handler, nil
}

func createHandler(encInfo *EncryptionInfo, backupHandler files.BackupHandler) func(Communicator) {
	suiteMap := encInfo.RetrievalInfo.SuiteInfos
	f := func(c Communicator) {
		handler, e := firstReply(c, encInfo, backupHandler)
		if e != nil {
			handleFailure(c, e, suiteMap)
		} else {
			e = handler(c)
			if e != nil {
				handleFailure(c, e, suiteMap)
			}
		}
	}
	return f
}


func getIndexes(s string) (int64, int64, error) {
	fields := strings.Split(s, ";")
	if len(fields) != 3 {
		return -1, -1, new(ErrorFailedProtocol)
	}
	start, e := strconv.Atoi(fields[1])
	if e != nil {
		return -1, -1, e
	}
	size, e := strconv.Atoi(fields[2])
	if e != nil {
		return -1, -1, e
	}
	return int64(start), int64(size), nil
}
//UploadData takes in a slice of bytes
//and sends it the given peer.
func UploadData(d []byte, comm Communicator, encInfo *EncryptionInfo) (uint64, error) {
	info := encInfo.RetrievalInfo
	suite := info.Suite
	//fmt.Fprintf(conn, info.Suite.String() + "\n")
	suiteBlob, e := encInfo.PURBAnon([]byte(suite.String()))
	if e != nil {
		return 0, e
	}
	e = comm.SendMessage(suiteBlob)
	if e != nil {
		return 0, e
	}
	reply, e := comm.GetNextMessage()
	if e != nil {
		return 0, e
	}
	pk, e := verifyPublicKey(reply, encInfo, suite, suiteBlob)
	if e != nil {
		fmt.Println("Yolo")
		return 0, e
	}
	marshalledKey, e := pk.MarshalBinary()
	if e != nil {
		return 0, e
	}
	recipient, e := purb.NewRecipient(marshalledKey, suite)
	if e != nil {
		return 0, e
	}
	recipients := []purbs.Recipient{recipient}
	blob, e := signAndPURB(encInfo, recipients, suite, d)
	if e != nil {
		return 0, e
	}
	//blob := encInfo.Enc.EncodeToString(signedBlob)
	e = comm.SendMessage(blob)
	if e != nil {
		return 0, e
	}
	//fmt.Fprintf(conn, blob)
	message, e := comm.GetNextMessage()
	if e != nil {
		return 0, e
	}
	index, e := extractIndexFromMessage(string(message))
	if e != nil {
		return 0, e
	}
	e = comm.CloseChannel()
	return index, e
}

func SimpleUpload(d []byte, comm Communicator, encInfo *EncryptionInfo) (uint64, error) {
	sig, e := encInfo.NormalSign(d)
	if e != nil {
		return 0, e
	}
	encData := encInfo.Enc.EncodeToString(d)
	encSig := encInfo.Enc.EncodeToString(sig)
	s := "Norm" + seperator + encData + seperator + encSig
	comm.SendMessage([]byte(s))
	resp, e := comm.GetNextMessage()
	fmt.Println("Yolo")
	if e != nil {
		return 0, e
	}
	spl := strings.Split(string(resp), seperator)
	if len(spl) < 2 {
		return 0, errors.New("failed response")
	}
	ind, resSig := spl[0], spl[1]
	e = encInfo.NormalVerify([]byte(ind), []byte(resSig))
	if e != nil {
		return 0, e
	}
	u, e := strconv.Atoi(ind)
	if e != nil {
		return 0, e
	}
	return uint64(u), nil
}

func GetNormalHandler(msg []byte, info *EncryptionInfo, bh files.BackupHandler) (communicationHandler, error) {
	s := string(msg)
	if !strings.Contains(s, "Norm") {
		return nil, errors.New("this was not a normal request")
	}
	return func(c Communicator) error {
		spl := strings.Split(s, seperator)
		chunk, sig := spl[1], spl[2]
		data, e := info.Enc.DecodeFromString(chunk)
		if e != nil {
			return e
		}
		signature, e := info.Enc.DecodeFromString(sig)
		if e != nil {
			return e
		}
		v := info.NormalVerify(data, signature)
		if v != nil {
			return v
		}
		ind := bh.AddBackup(data)
		if ind == -1 {
			return errors.New("could not add to permanent storage")
		}
		u := strconv.FormatInt(ind, 10)
		su, e := info.Sign([]byte(u))
		if e != nil {
			return e
		}
		r := u + seperator + info.Enc.EncodeToString(su)
		return c.SendMessage([]byte(r))
	}, nil
}