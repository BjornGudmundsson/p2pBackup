package files

import (
	"errors"
	"time"
)

//GetTimePassedSinceModified take in a file and
//returns how much time has passed since the file was
//modified in seconds.
func GetTimePassedSinceModified(f File) float64 {
	now := time.Now()
	diff := now.Sub(f.Modified)
	return diff.Seconds()
}

//HasPassedTimeLimit takes in a file and a duration and
//returns nil if the time limit has been exceeded or an
//error if not.
func HasPassedTimeLimit(f File, t time.Duration) error {
	elapsed := GetTimePassedSinceModified(f)
	if t.Seconds() <= elapsed {
		return nil
	}
	return errors.New("Not enough time has elapsed")
}
