package peers

import (
	"fmt"
	"time"

	"github.com/BjornGudmundsson/p2pBackup/files"
)

//Here are the functions and objects that get the data that is supposed to be collected
//and send it if meant to

func checkIfHasBeenBackedup(data []byte) bool {
	return true
}

//Update is meant to be run as a seperate thread that periodically checks for data
//to send and backup to its peers. The wait parameter defines the amount of time to wait between
//searching for new backups and the basedir says where to find the files to backup. rules is
//used to assist in automatic filter of non-backupable files.
func Update(wait time.Duration, basedir string, rules files.BackupData, peerFile, backupLog string) {
	peers, e := GetPeerList(peerFile)
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
	for {
		backupFiles, e := files.FindAllFilesToBackup(rules, basedir)
		if e != nil {
			fmt.Println(e)
		} else {
			data, e := files.ToBytes(backupFiles)
			if e != nil {
				fmt.Println("Could not read the files")
				fmt.Println(e)
			} else {
				fmt.Println(string(data))
				hasBeenBackedup := checkIfHasBeenBackedup(data)
				if !hasBeenBackedup {
					for _, peer := range peers {
						e = SendTCPData(data, peer)
						if e != nil {
							fmt.Println("Could not send data over tcp")
							fmt.Println(e.Error())
						}
					}
					files.AddBackup(data, backupLog)
				}
			}
		}
		time.Sleep(wait)
	}
}
