package files

import (
	"io/ioutil"
	"os"

	"github.com/BjornGudmundsson/p2pBackup/utilities"
)

//In this file, is everything required to find and traverse the file directories.

//Exists returns whether a file with a given file exists.
func Exists(fn string) bool {
	if _, err := os.Stat(fn); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		return false

	}
}

//GetFilesFromDir take in a directory name and reads
//all of the files from that directory.
func GetFilesFromDir(dir string) ([]File, error) {
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		return nil, e
	}
	ret := make([]File, len(files))
	for i, file := range files {
		ret[i] = NewFile(file, dir)
	}
	return ret, nil
}

//TraverseDirForFiles takes in a directory and returns all
//of the files in that directory
func TraverseDirForFiles(dir string) ([]File, error) {
	files, e := GetFilesFromDir(dir)
	if e != nil {
		return nil, e
	}
	ret := make([]File, 0)
	fileStack := utilities.NewStack()
	for _, f := range files {
		if f.Type == DIR {
			fileStack.Push(f)
		} else {
			ret = append(ret, f)
		}
	}
	for !fileStack.IsEmpty() {
		f := fileStack.Pop().(File)
		d, err := GetFilesFromDir(f.Path + "/" + f.Name)
		if err != nil {
			return nil, err
		}
		for _, f2 := range d {
			if f2.Type == DIR {
				fileStack.Push(f2)
			} else {
				ret = append(ret, f2)
			}
		}
	}
	return ret, nil
}
