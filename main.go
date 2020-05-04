package main

import (
	"encoding/hex"
	"fmt"
	"github.com/BjornGudmundsson/p2pBackup/crypto"
	"github.com/BjornGudmundsson/p2pBackup/kyber"
	"github.com/BjornGudmundsson/p2pBackup/purb/purbs"
	"github.com/BjornGudmundsson/p2pBackup/utilities"
	"net/http"
	"strconv"
	"time"

	"github.com/BjornGudmundsson/p2pBackup/purb"

	"github.com/BjornGudmundsson/p2pBackup/files"
	_ "github.com/BjornGudmundsson/p2pBackup/files"
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
	flags := utilities.NewFlags()
	if flags.GetBool("init") {
		keyPairs(flags.GetString("suite"))
		return
	}
	s, e := purb.GetSuite(flags.GetString("suite"))
	if e != nil {
		fmt.Println(e)
		return
	}
	authKey, e := hexToKey(flags.GetString("authkey"), s)
	if e != nil {
		fmt.Println("Not a valid authentication key: ", e)
		return
	}
	container, e := peers.NewContainerFromFile(flags.GetString("peers"))
	if e != nil {
		fmt.Println("Could not get the peer list")
		return
	}
	logfile := flags.GetString("logfile")
	logHandler, e := files.NewEncryptedLogWriter(logfile, authKey.String())
	if e != nil {
		fmt.Println("Can't create the log")
		return
	}
	sk, e := hex.DecodeString(flags.GetString("key"))
	if e != nil {
		fmt.Println("Your secret key was not valid, here's a new one: Pretend there is a secret key")
		return
	}

	info, e := purb.NewKeyInfo(sk, s, flags.GetString("Suites"))
	if e != nil {
		fmt.Println("Error: " , e.Error())
	}
	udp:= flags.GetString("udp")
	backupHandler := files.NewBackupBuffer("./"+flags.GetString("storage"))
	auth, e := crypto.NewAnonAuthenticator(s, flags.GetString("set"))
	encInfo := peers.NewEncryptionInfo(auth, authKey, nil, info, "", nil)
	if flags.GetBool("gui") {
		pages.ServeScripts()
		http.HandleFunc("/", pages.IndexPage)
		http.HandleFunc("/backup", pages.BackupFile)
		go http.ListenAndServe(":"+flags.GetString("p"), nil)
	}
	base := flags.GetString("base")
	if flags.GetBool("retrieve") {
		backup, e := peers.RetrieveFromLogs(logHandler, encInfo, container)
		if e != nil {
			fmt.Println(e)
		} else {
			files.ReconstructBackup(backup, base)
		}
	} else {
		//PrintInfo(baseDir, peersList, udpPort, rules, storageFile, backupLogs, filePort, false)
		timer, e := time.ParseDuration(flags.GetString("backuprate"))
		if e != nil {
			panic(e)
		}
		backupRules := files.CreateRules(flags.GetString("backuprules"))
		go peers.Update(timer, base, backupRules, container, logfile, encInfo, logHandler)
		go peers.ListenUDP(":" + udp)
		peers.ListenTCP(":"+strconv.Itoa(flags.GetInt("fileport")), encInfo, backupHandler)
	}
}
