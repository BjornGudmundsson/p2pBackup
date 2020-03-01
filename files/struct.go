package files

import (
	"os"
	"strconv"
	"time"
)

//FILE is just a constant
const FILE = "file"

//DIR is just a constant
const DIR = "dir"

//File keeps info about a found file
type File struct {
	//Local name of the file
	Name string
	//What kind of file it is
	Type string
	//Absolute location of the file w.r.t the starting dir
	Path string
	//Modified is the time in which the file was modified last
	Modified time.Time
	//Size is the size of the file
	Size int64
}

//NewFile takes in the name of a file
//and returns all the necessary info about that file.
func NewFile(f os.FileInfo, dir string) File {
	var t string
	if f.IsDir() {
		t = DIR
	} else {
		t = FILE
	}
	return File{
		Name:     f.Name(),
		Type:     t,
		Path:     dir,
		Modified: f.ModTime(),
		Size:     f.Size(),
	}
}

func (f File) String() string {
	n := "Name: " + f.Name + "\n"
	dir := "Type: " + f.Type + "\n"
	path := "Path: " + f.Path + "\n"
	size := "Size: " + strconv.Itoa(int(f.Size)) + "\n"
	mod := "Modified: " + f.Modified.String() + "\n"
	return n + dir + path + size + mod
}

//Equal takes in two file objects and returns
//Whether they are referring to the same file
func (f File) Equal(f2 File) bool {
	return f.Name == f2.Name && f.Path == f2.Path && f.Type == f2.Type
}
