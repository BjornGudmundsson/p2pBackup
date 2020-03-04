package main

import (
	"flag"
	"fmt"
	"net/http"

	_ "github.com/BjornGudmundsson/p2pBackup/files"
	"github.com/BjornGudmundsson/p2pBackup/pages"
	"github.com/BjornGudmundsson/p2pBackup/peers"
)

func main() {
	port := flag.String("p", "8080", "Which port to run the server on")
	baseDir := flag.String("base", ".", "Base is the basedirectory in which all files will be backed up from. If not provided it will default to the running directory")
	peersList := flag.String("peers", "empty", "Peers is the file in which the data about other peers is stored")
	udpPort := flag.String("udp", "3000", "UDP is the port that will be used for the udp socket")
	rules := flag.String("backuprules", "", "backuprules is the toml file in which the specifications for the backup are kept")
	gui := flag.Bool("gui", false, "Gui says whether or not a gui should be displayed or not. Defaults to false")
	flag.Parse()
	if *gui {
		pages.ServeScripts()
		http.HandleFunc("/", pages.IndexPage)
		http.HandleFunc("/backup", pages.BackupFile)
		fmt.Println("Running server on port: ", *port)
	}
	fmt.Println("Backing up files from: ", *baseDir)
	fmt.Println("Reading peers from: ", *peersList)
	fmt.Println("Listening for udp packets on: ", *udpPort)
	fmt.Println("Reading rules from: ", *rules)
	go peers.ListenUDP(":" + *udpPort)
	e := http.ListenAndServe(":"+(*port), nil)
	fmt.Println(e.Error())
}
