package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/BjornGudmundsson/p2pBackup/files"
	_ "github.com/BjornGudmundsson/p2pBackup/files"
	"github.com/BjornGudmundsson/p2pBackup/pages"
	"github.com/BjornGudmundsson/p2pBackup/peers"
)

func KeyPairs(protocol string) ([]byte, []byte, string) {
	return nil, nil, ""
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
	flag.Parse()
	if *gui {
		pages.ServeScripts()
		http.HandleFunc("/", pages.IndexPage)
		http.HandleFunc("/backup", pages.BackupFile)
		fmt.Println("Running server on port: ", *port)
		go http.ListenAndServe(":"+(*port), nil)
	}
	timer, e := time.ParseDuration(*updateTimer)
	if e != nil {
		panic(e)
	}
	backupRules := files.CreateRules(*rules)
	fmt.Println("Backing up files from: ", *baseDir)
	fmt.Println("Reading peers from: ", *peersList)
	fmt.Println("Listening for udp packets on: ", *udpPort)
	fmt.Println("Reading rules from: ", *rules)
	fmt.Println("Storing backups in: ", *storageFile)
	fmt.Println("Backup download port is at port: ", *filePort)
	fmt.Println("Storing backup logs at: ", *backupLogs)
	go peers.Update(timer, *baseDir, backupRules, *peersList, *backupLogs)
	go peers.ListenUDP(":" + *udpPort)
	peers.ListenTCP(":"+strconv.Itoa(*filePort), "./"+*storageFile)
}
