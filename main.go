package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/BjornGudmundsson/p2pBackup/crypto"
	"github.com/BjornGudmundsson/p2pBackup/kyber"
	"github.com/BjornGudmundsson/p2pBackup/purb/purbs"
	"net/http"
	"strconv"
	"time"

	"github.com/BjornGudmundsson/p2pBackup/purb"

	"github.com/BjornGudmundsson/p2pBackup/files"
	_ "github.com/BjornGudmundsson/p2pBackup/files"
	"github.com/BjornGudmundsson/p2pBackup/kyber/group/curve25519"
	"github.com/BjornGudmundsson/p2pBackup/kyber/util/key"
	"github.com/BjornGudmundsson/p2pBackup/pages"
	"github.com/BjornGudmundsson/p2pBackup/peers"
)

func keyPairs(suite string) {
	s, e := purb.GetSuite(suite)
	if e != nil {
		fmt.Println("No suite available")
		return
	}
	keys := key.NewKeyPair(s)
	sk, e := keys.Private.MarshalBinary()
	if e != nil {
		fmt.Println("Could not generate a new key pair because: ", e.Error())
	} else {
		info, e := purb.NewKeyInfo(sk, s, "")
		if e != nil {
			fmt.Println(e.Error())
		} else {
			fmt.Println(info)
		}

	}
}

func PrintInfo(baseDir, peersList, udpPort, rules, storageFile,backupLogs *string, filePort *int, toPrint bool) {
	if toPrint {
		fmt.Println("Backing up files from: ", *baseDir)
		fmt.Println("Reading peers from: ", *peersList)
		fmt.Println("Listening for udp packets on: ", *udpPort)
		fmt.Println("Reading rules from: ", *rules)
		fmt.Println("Storing backups in: ", *storageFile)
		fmt.Println("Backup download port is at port: ", *filePort)
		fmt.Println("Storing backup logs at: ", *backupLogs)
	}
}


func hexToKey(hx string, suite purbs.Suite) (kyber.Scalar, error) {
	b, e := hex.DecodeString(hx)
	if e != nil {
		return nil, e
	}
	scalar := suite.Scalar()
	e = scalar.UnmarshalBinary(b)
	return scalar, e
}

func main() {
	port := flag.String("p", "8080", "Which port to run the server on")
	baseDir := flag.String("base", ".", "Base is the basedirectory in which all files will be backed up from. If not provided it will default to the running directory")
	peersList := flag.String("peers", "peers.txt", "Peers is the file in which the data about other peers is stored")
	udpPort := flag.String("udp", "5000", "UDP is the port that will be used for the udp socket")
	rules := flag.String("backuprules", "", "backuprules is the toml file in which the specifications for the backup are kept")
	gui := flag.Bool("gui", false, "Gui says whether or not a gui should be displayed or not. Defaults to false")
	storageFile := flag.String("storage", "backup.txt", "Storage is the location in which you prefer to store your peers backups")
	filePort := flag.Int("fileport", 3000, "The port in which a tcp connection can be made to send the backup")
	backupLogs := flag.String("logfile", "backuplog.txt", "This is where the users wishes to store all log of backups they have performed")
	updateTimer := flag.String("backuprate", "1s", "Backuprate tells how often the system should scan for whether it should update")
	initialize := flag.Bool("init", false, "Init is to tell wheter the user wants to get a new private/public key pair")
	key := flag.String("key", "", "What the secret key that the user will be using to encrypt their backups")
	authString := flag.String("authkey", "", "The key used to prove you are part of the p2p backup group")
	suite := flag.String("suite", curve25519.NewBlakeSHA256Curve25519(true).String(), "What ciphersuite the user decides to use")
	suiteFile := flag.String("Suites", "", "Where all the known suites are stored in a TOML file")
	setString := flag.String("set", "", "The file containing the public keys of the anonymity set")
	flag.Parse()
	if *gui {
		pages.ServeScripts()
		http.HandleFunc("/", pages.IndexPage)
		http.HandleFunc("/backup", pages.BackupFile)
		fmt.Println("Running server on port: ", *port)
		go http.ListenAndServe(":"+(*port), nil)
	}
	if *initialize {
		keyPairs(*suite)
	} else {
		PrintInfo(baseDir, peersList, udpPort, rules, storageFile, backupLogs, filePort, false)
		timer, e := time.ParseDuration(*updateTimer)
		if e != nil {
			panic(e)
		}
		backupRules := files.CreateRules(*rules)

		sk, e := hex.DecodeString(*key)
		if e != nil {
			fmt.Println("Your secret key was not valid, here's a new one: Pretend there is a secret key")
			return
		}
		s, e := purb.GetSuite(*suite)
		if e != nil {
			fmt.Println(e)
			return
		}
		info, e := purb.NewKeyInfo(sk, s, *suiteFile)
		if e != nil {
			fmt.Println("Error: " , e.Error())
		}
		authKey, e := hexToKey(*authString, s)
		if e != nil {
			fmt.Println("Not a valid authentication key: ", e)
			return
		}
		container, e := peers.NewContainerFromFile(*peersList)
		if e != nil {
			fmt.Println("Could not get the peer list")
			return
		}
		logHandler, e := files.NewEncryptedLogWriter(*backupLogs, authKey.String())
		if e != nil {
			fmt.Println("Can't create the log")
			return
		}
		backupHandler := files.NewBackupBuffer("./"+*storageFile)
		auth, e := crypto.NewAnonAuthenticator(s, *setString)
		encInfo := peers.NewEncryptionInfo(auth, authKey, nil, info, "", nil)
		go peers.Update(timer, *baseDir, backupRules, container, *backupLogs, encInfo, logHandler)
		go peers.ListenUDP(":" + *udpPort)
		peers.ListenTCP(":"+strconv.Itoa(*filePort), encInfo, backupHandler)
	}
}
