package files

import (
	"errors"
	"time"
)

//GetTimePassedSinceModified take in a file and
//returns how much time has passed since the file was
//modified in seconds.
func GetTimePassedSinceModified(f File) time.Duration {
	now := time.Now()
	diff := now.Sub(f.Modified)
	return diff
}

//HasPassedTimeLimit takes in a file and a duration and
//returns nil if the time limit has been exceeded or an
//error if not.
func HasPassedTimeLimit(f File, t time.Duration) error {
	elapsed := GetTimePassedSinceModified(f)
	if t-elapsed <= time.Duration(0) {
		return nil
	}
	return errors.New("Not enough time has elapsed")
}

//FindAllFilesToBackup takes in a set of rules and a base directory and
//returns all of the files that should be backed up according to the
//given rules and returns an error otherwise.
func FindAllFilesToBackup(rules BackupData, dir string) ([]File, error) {
	files, e := TraverseDirForFiles(dir)
	if e != nil {
		return nil, nil
	}
	fileToBackup := make([]File, 0)
	for _, file := range files {
		if rules.Include(file) {
			fileToBackup = append(fileToBackup, file)
		}
	}
	return fileToBackup, nil
}
