package peers

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/BjornGudmundsson/p2pBackup/kyber"
	"github.com/BjornGudmundsson/p2pBackup/kyber/util/key"
	"github.com/BjornGudmundsson/p2pBackup/kyber/util/random"
	"github.com/BjornGudmundsson/p2pBackup/purb"
	"github.com/BjornGudmundsson/p2pBackup/purb/purbs"
	"net"
	"strconv"
	"sync"

	"github.com/BjornGudmundsson/p2pBackup/files"
)

var fileLock *sync.Mutex

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
func ListenTCP(port string, backupFile string, info *purb.KeyInfo) {
	l, e := net.Listen("tcp4", port)
	if e != nil {
		panic(e)
	}
	defer l.Close()
	handler := createHandler(backupFile, info.SuiteInfos)
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

//BackupData takes in a file
//descriptor structure and
//data and appends it to the file
func BackupData(f file, data []byte) error {
	verified := verifyData(data)
	if !verified {
		return NotVerifiedError()
	}
	e := files.AppendToFile(files.File(f), data)
	if e != nil {
		return e
	}
	return nil
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

func getUploadHandler(suite purbs.Suite, fn string, suiteMap purbs.SuiteInfoMap) func(c net.Conn) error {
	fd, e := files.GetFile(fn)
	if e != nil {
		panic(e)
	}
	fl := file(*fd)
	f := func(c net.Conn) error {
		params := purbs.NewPublicFixedParameters(suiteMap, false)
		freshPair := key.NewKeyPair(suite)
		publicBytes, e := freshPair.Public.MarshalBinary()
		if e != nil {
			return nil
		}
		privateBytes, e := freshPair.Private.MarshalBinary()
		if e != nil {
			return e
		}
		self, e := purb.NewPrivateRecipient(privateBytes, suite)
		if e != nil {
			return e
		}
		_, e = fmt.Fprintf(c, hex.EncodeToString(publicBytes) + "\n")
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
		success, data, e := purbs.Decode(blob, &self, params, false)
		if e != nil {
			fmt.Println("Error: ", e.Error())
			return e
		}
		if !success {
			return errors.New("can't decode the initial purb")
		}
		e = BackupData(fl, data)
		if e != nil {
			return e
		}
		_, e = fmt.Fprintf(c, "Done writing\n")
		if e != nil {
			return e
		}
		return c.Close()
	}
	return f
}

func firstReply(conn net.Conn, fn string, m purbs.SuiteInfoMap) (tcpHandler, error) {
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
	handler := getUploadHandler(suite, fn, m)
	return handler, nil
}

func createHandler(fileName string, suiteMap purbs.SuiteInfoMap) func(net.Conn) {
	f := func(c net.Conn) {
		handler, e := firstReply(c, fileName, suiteMap)
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
func SendTCPData(d []byte, p *Peer, info *purb.KeyInfo) error {
	params := purbs.NewPublicFixedParameters(info.SuiteInfos, false)
	conn, e := net.Dial("tcp", p.Addr.String()+":"+strconv.Itoa(p.Port))
	if e != nil {
		return e
	}
	suite := info.Suite
	fmt.Fprintf(conn, info.Suite.String() + "\n")
	reply, e := bufio.NewReader(conn).ReadString('\n')
	reply = reply[:len(reply) - 1]
	pkBytes, e := hex.DecodeString(reply)
	if e != nil {
		return e
	}
	pk, e := getPublicKey(suite, pkBytes)
	if e != nil {
		return e
	}
	marshalledKey, e := pk.MarshalBinary()
	if e != nil {
		return e
	}
	recipient, e := purb.NewRecipient(marshalledKey, suite)
	if e != nil {
		return e
	}
	recipients := []purbs.Recipient{recipient}
	purb, e := purbs.Encode(d, recipients, random.New(), params, false)
	if e != nil {
		return e
	}
	blob := hex.EncodeToString(purb.ToBytes()) + "\n"
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

func getPublicKey(suite purbs.Suite, d []byte) (kyber.Point, error) {
	p := suite.Point()
	e := p.UnmarshalBinary(d)
	return p, e
}
