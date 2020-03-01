package files

import (
	"os"
)

//AppendToFile take in a file and a slice of bytes
//and appends them to the end of the file.
func AppendToFile(f File, data []byte) error {
	fn := f.Name
	fd, e := os.OpenFile(fn, os.O_WRONLY, os.ModeAppend)
	if e != nil {
		return e
	}
	start := f.Size
	_, e = fd.WriteAt(data, start)
	if e != nil {
		return e
	}
	return nil
}
