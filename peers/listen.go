package peers

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/BjornGudmundsson/p2pBackup/files"
	"github.com/BjornGudmundsson/p2pBackup/kyber/util/key"
	"github.com/BjornGudmundsson/p2pBackup/purb"
	"github.com/BjornGudmundsson/p2pBackup/purb/purbs"
	"net"
)


//Shorthand for a function handle a connection
type tcpHandler func(net.Conn) error

const localhost = "127.0.0.1"

//ListenUDP sets up a UDP server
//on the given port.
func ListenUDP(port string) {
	serverAddr, e := net.ResolveUDPAddr("udp", port)
	if e != nil {
		panic(e)
	}
	conn, e := net.ListenUDP("udp", serverAddr)
	if e != nil {
		fmt.Println(conn)
		fmt.Println(e.Error())
	}
	defer conn.Close()
	buf := make([]byte, 4096)
	for {
		n, addr, e := conn.ReadFromUDP(buf)
		fmt.Println("Got packet from: ", addr.String())
		fmt.Println("Message: ", string(buf[:n]))
		if e != nil {
			fmt.Println(e.Error())
		}
	}
}

//ListenTCP starts a new tcp server on the given port
func ListenTCP(port string, encInfo *EncryptionInfo, backupHandler files.BackupHandler) {
	l, e := net.Listen("tcp4", port)
	if e != nil {
		panic(e)
	}
	defer l.Close()
	handler := createHandler(encInfo, backupHandler)
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		} else {

			go handler(c)
		}
	}
}

//verifyData takes in bytes and verifies that
//this data was sent by a known peer
func verifyData(data []byte) bool {
	//TODO: Write the function in such a way that it compares the data against all public keys.
	return true
}



func handleFailure(c net.Conn, e error, m purbs.SuiteInfoMap) {
	if c != nil {
		fmt.Fprintf(c, e.Error())
		c.Close()
	}
}

func isDownload(msg string) (bool, error) {
	return false, nil//TODO: /make it such that returns true if download, false else and an error if it means nothing
}

func getUploadHandler(suite purbs.Suite,encInfo *EncryptionInfo, backupHandler files.BackupHandler) func(c net.Conn) error {
	f := func(c net.Conn) error {
		freshPair := key.NewKeyPair(suite)
		publicBytes, e := freshPair.Public.MarshalBinary()
		if e != nil {
			return nil
		}
		sig, e := signPublicKey(publicBytes, encInfo)
		if e != nil {
			return e
		}
		_, e = fmt.Fprintf(c, hex.EncodeToString(sig) + "\n")
		if e != nil {
			return e
		}
		reader := bufio.NewReader(c)
		hxPurb, e := reader.ReadString('\n')
		if e != nil {
			return e
		}
		hxPurb = hxPurb[:len(hxPurb) - 1]
		blob, e := hex.DecodeString(hxPurb)
		if e != nil {
			return e
		}
		data, e := verifyPURB(freshPair.Private, suite, blob, encInfo)
		ind := backupHandler.AddBackup(data)
		if ind == -1 {
			return new(ErrorCouldNotAppend)
		}
		_, e = fmt.Fprintf(c, "Done writing\n")
		if e != nil {
			return e
		}
		return c.Close()
	}
	return f
}

func firstReply(conn net.Conn,encInfo *EncryptionInfo, backupHandler files.BackupHandler) (tcpHandler, error) {
	reader := bufio.NewReader(conn)
	s, e := reader.ReadString('\n')
	if e != nil {
		return nil, e
	}
	s = s[:len(s) - 1]
	download, e := isDownload(s)
	if e != nil {
		return nil, e
	}
	if download {
		//Handle download
		return nil, nil
	}
	suite, e := purb.GetSuite(s)
	if e != nil {
		return nil, e
	}
	handler := getUploadHandler(suite, encInfo, backupHandler)
	return handler, nil
}

func createHandler(encInfo *EncryptionInfo, backupHandler files.BackupHandler) func(net.Conn) {
	suiteMap := encInfo.RetrievalInfo.SuiteInfos
	f := func(c net.Conn) {
		handler, e := firstReply(c, encInfo, backupHandler)
		if e != nil {
			handleFailure(c, e, suiteMap)
		}
		e = handler(c)
		if e != nil {
			fmt.Println("Did not succeed")
			handleFailure(c, e, suiteMap)
		}
	}
	return f
}

//SendTCPData takes in a slice of bytes
//and sends it the given peer.
func SendTCPData(d []byte, p Peer, encInfo *EncryptionInfo) error {
	info := encInfo.RetrievalInfo
	conn, e := getTCPConn(p)
	if e != nil {
		return e
	}
	suite := info.Suite
	fmt.Fprintf(conn, info.Suite.String() + "\n")
	reply, e := bufio.NewReader(conn).ReadString('\n')
	reply = reply[:len(reply) - 1]
	pk, e := verifyPublicKey([]byte(reply), encInfo, suite)
	if e != nil {
		fmt.Println("Could not verify the signature")
		return e
	}
	marshalledKey, e := pk.MarshalBinary()
	if e != nil {
		fmt.Println("Could not marshalled to binary")
		return e
	}
	recipient, e := purb.NewRecipient(marshalledKey, suite)
	if e != nil {
		return e
	}
	recipients := []purbs.Recipient{recipient}
	signedBlob, e := signAndPURB(encInfo, recipients, suite, d)
	if e != nil {
		return e
	}
	blob := hex.EncodeToString(signedBlob) + "\n"
	fmt.Fprintf(conn, blob)
	message, e := bufio.NewReader(conn).ReadString('\n')
	message = message[:len(message) - 1]
	if e != nil {
		fmt.Println("Error: ", e.Error())
		return e
	}
	e = conn.Close()
	return e
}

