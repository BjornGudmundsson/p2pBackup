package main

import (
	"flag"
	"fmt"
	"net/http"

	_ "github.com/BjornGudmundsson/p2pBackup/files"
	"github.com/BjornGudmundsson/p2pBackup/pages"
)

func main() {
	port := flag.String("p", "8080", "Which port to run the server on")
	flag.Parse()
	pages.ServeScripts()
	http.HandleFunc("/", pages.IndexPage)
	http.HandleFunc("/backup", pages.BackupFile)
	fmt.Println("Running server on port: ", *port)
	http.ListenAndServe(":"+(*port), nil)
	fmt.Println("Bjorn")
}
