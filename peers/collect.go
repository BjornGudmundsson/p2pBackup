package peers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/BjornGudmundsson/p2pBackup/files"
)

//Update is meant to be run as a seperate thread that periodically checks for data
//to send and backup to its peers. The wait parameter defines the amount of time to wait between
//searching for new backups and the basedir says where to find the files to backup. rules is
//used to assist in automatic filter of non-backupable files.
func Update(wait time.Duration, basedir string, rules files.BackupData, peerContainer Container, backupLog string, encInfo *EncryptionInfo, handler files.LogWriter) {
	for {
		time.Sleep(wait)
		peers := peerContainer.GetPeerList()
		backupFiles, e := files.FindAllFilesToBackup(rules, basedir)
		if e != nil {
			fmt.Println(e)
		} else {
			data, e := files.ToBytes(backupFiles)
			if e != nil {
				fmt.Println("Could not read the files")
				fmt.Println(e)
			} else {
				if e != nil {
					panic(e)
				}
				indexes := make([]uint64, 0)
				hasBeenBackedUp, e := handler.CheckIfBackedUp(data)
				if !hasBeenBackedUp && len(data) != 0 && e == nil {
					ct, e := encInfo.PURBBackup(data)
					if e != nil {
						fmt.Println("Could not purbify", e)
						continue
					}
					for _, peer := range peers {
						index, e := SendTCPData(data, peer, encInfo)
						if e != nil {
							fmt.Println("Could not send data over tcp")
							fmt.Println(e.Error())
						} else {
							indexes = append(indexes, index)
						}
					}
					log := handler.NewLog(data, indexes, ct)
					e = handler.Log(log)
				}
			}
		}
	}
}

func extractIndexFromMessage(msg string) (uint64, error) {
	fields := strings.Split(msg, " ")
	if len(fields) != 2 {
		return 0, new(ErrorIncorrectFormat)
	}
	ind, e := strconv.Atoi(fields[1])
	return uint64(ind), e
}
