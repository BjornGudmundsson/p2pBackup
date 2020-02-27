package files

import (
	"os"
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
		Name: f.Name(),
		Type: t,
		Path: dir,
	}
}

func (f File) String() string {
	n := "Name: " + f.Name + "\n"
	dir := "Type: " + f.Type + "\n"
	path := "Path: " + f.Path + "\n"
	return n + dir + path
}
