package files

import (
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
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

type duration struct {
	Duration time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

//BackupData keeps track of the rules a file has to
//follow in order to be backuped
type BackupData struct {
	//MaxSize indicates how large the file can be in bytes.
	MaxSize int64
	//MinSize is the lower limit for file size in bytes.
	MinSize int64
	//TypesToExclude is an expression detailing what kind of files should not be backed up.
	TypesToExclude string
	//BlackListedFiles are files that have been chosen in particular not to be included.
	BlackListedFiles []string
	//MinTimeSinceModified is how many seconds must have elapsed for the file to be considered for backing up
	MinTimeSinceModified duration
}

//GetMinTime returns the min time in time.Duration
func (bd BackupData) GetMinTime() time.Duration {
	return bd.MinTimeSinceModified.Duration
}

//Include takes in a file and returns true
//if the file satisfies the rules set in place for a backup.
func (bd BackupData) Include(f File) bool {
	if bd.GetMinTime() > GetTimePassedSinceModified(f) {
		return false
	}
	r, e := regexp.Compile(bd.TypesToExclude)
	if e != nil {
		return false
	}
	if r.MatchString(f.Name) {
		return false
	}
	name := f.Path + "/" + f.Name
	for _, n := range bd.BlackListedFiles {
		if n == name {
			return false
		}
	}
	if bd.MinSize >= f.Size || bd.MaxSize < f.Size {
		return false
	}
	return true
}

//CreateRules takes in a name of a file that has
//the specified rules while if the no file is specified or the
//configuration file can't be found a set of default rules will be used.
//The file given must be a TOML file
func CreateRules(rules string) BackupData {
	var d BackupData
	_, e := toml.DecodeFile(rules, &d)
	if e != nil {
		return DefaultRules()
	}
	return d
}

//DefaultRules returns the rules if nothing else has been specified
func DefaultRules() BackupData {
	d := BackupData{}
	d.MaxSize = 1000
	d.MinSize = 0
	d.TypesToExclude = "([a-z]*).csv"
	d.BlackListedFiles = []string{}
	d.MinTimeSinceModified = duration{
		Duration: time.Second,
	}
	return d
}
