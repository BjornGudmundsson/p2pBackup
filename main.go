package main

import (
	"fmt"
	"github.com/BjornGudmundsson/p2pBackup/utilities"
	"net/http"
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

func main() {
	flags := utilities.NewFlags()
	if flags.GetBool("init") {
		keyPairs(flags.GetString("suite"))
		return
	}
	container, e := peers.NewContainerFromFile(flags.GetString("peers"))
	if e != nil {
		fmt.Println("Could not get the peer list: ", e)
		return
	}
	logfile := flags.GetString("logfile")
	logHandler, e := files.NewEncryptedLogWriter(logfile, flags.GetString("pw"))
	if e != nil {
		fmt.Println("Can't create the log: ", e)
		return
	}
	backupHandler := files.NewBackupBuffer("./"+flags.GetString("storage"))
	encInfo, e := peers.GetEncryptionInfoFromFlags(flags)
	if e != nil {
		fmt.Println(e)
		return
	}
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
		if flags.GetBool("update") {
			timer, e := time.ParseDuration(flags.GetString("backuprate"))
			if e != nil {
				fmt.Println(e)
				return
			}
			backupRules := files.CreateRules(flags.GetString("backuprules"))
			go peers.Update(timer, base, backupRules, container, logfile, encInfo, logHandler)
		}
		server, e := peers.NewServer(flags, backupHandler, encInfo, container)
		if e != nil {
			fmt.Println(e)
			return
		}
		go server.FindPeers()
		server.Listen()
	}
}
