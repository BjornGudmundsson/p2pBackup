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
func Update(w time.Duration, dir string, rules files.BackupData, container Container, enc *EncryptionInfo, h files.LogWriter) {
	for {
		time.Sleep(w)
		BackupDate(dir, rules, container, enc, h)
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

func BackupDate(dir string, rules files.BackupData, container Container, enc *EncryptionInfo, h files.LogWriter) {
	peers := container.GetPeerList()
	backupFiles, e := files.FindAllFilesToBackup(rules, dir)
	if e != nil {
		fmt.Println(e)
	} else {
		data, e := files.ToBytes(backupFiles)
		if e != nil {
			fmt.Println("Could not read the files: ", e)
		} else {
			if e != nil {
				panic(e)
			}
			indexes := make([]uint64, 0)
			hasBeenBackedUp, e := h.CheckIfBackedUp(data)
			if !hasBeenBackedUp && len(data) != 0 && e == nil {
				ct, e := enc.PURBBackup(data)
				if e != nil {
					return
				}
				for _, peer := range peers {
					now := time.Now().Nanosecond()
					comm, e := NewCommunicatorFromPeer(peer, enc)
					if e != nil {
						fmt.Println(e)
						continue
					}
					index, e := UploadData(ct, comm, enc)
					if e != nil {
						fmt.Println(e.Error())
					} else {
						fmt.Println("Success")
						indexes = append(indexes, index)
					}
					duration := time.Now().Nanosecond() - now
					if duration > 0 {
						fmt.Println("Elapsed: ", time.Now().Nanosecond()-now, ",", len(data))
					}
				}
				log := h.NewLog(data, indexes, ct)
				e = h.Log(log)
			}
		}
	}
}
