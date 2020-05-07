package files

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

//AppendToFile take in a file and a slice of bytes
//and appends them to the end of the file.
func AppendToFile(f File, data []byte) error {
	var fn string
	if f.Path == "" {
		fn = "./" + f.Name
	} else {
		fn = f.Path + "/" + f.Name
	}
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

//ToBytes takes a list of files that are to be backed up
//and returns them as a single slice of bytes
func ToBytes(files []File) ([]byte, error) {
	data := make([]byte, 0)
	for _, f := range files {
		name := "Name:" + f.Name + " "
		path := "Path:" + f.Path + " "
		size := "Size:" + strconv.Itoa(int(f.Size)) + " "
		sum := name + path + size + "\n"
		d, e := ioutil.ReadFile(f.Path + "/" + f.Name)
		if e != nil {
			return nil, e
		}
		data = append(data, []byte(sum)...)
		data = append(data, d...)
	}
	compressed, e := compressData(data, compress)
	if e != nil {
		return nil, e
	}
	return compressed, nil
}

//GetFile takes in a filename and returns
//a file descriptor object that describes the given name.
//It is assumed that the filename is the absolute name of the file,
//at least in terms of relative the the running directory.
//GetFile returns an error if the file descriptor could not be made.
func GetFile(fn string) (*File, error) {
	f, e := os.Open(fn)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	inf, e := f.Stat()
	if e != nil {
		return nil, e
	}
	file := NewFile(inf, "")
	return &file, nil
}


func GetLastNBytes(fn string, n int64) ([]byte, error) {
	f, e := os.OpenFile(fn, os.O_RDONLY, os.ModeAppend)
	defer f.Close()
	if e != nil {
		return nil, e
	}
	info, e := f.Stat()
	if e != nil {
		return nil, e
	}
	size := info.Size()
	if size == 0 {
		return make([]byte, n), nil
	}
	if size < n {
		d := make([]byte, size)
		_, e = f.Read(d)
		return d, e
	}
	offset := size - n
	d := make([]byte, n)
	_, e  = f.ReadAt(d, offset)
	return d, e
}


func compressData(d []byte, do bool) ([]byte, error) {
	if !do {
		return d, nil
	}
	var buf bytes.Buffer
	wr, e := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if e != nil {
		return nil, e
	}
	wr.Write(d)
	e = wr.Close()
	if e != nil {
		return nil, e
	}
	wr.Reset(&buf)
	return buf.Bytes(), nil
}

func decompressData(p []byte, do bool) ([]byte, error) {
	if !do {
		return p, nil
	}
	buf := new(bytes.Buffer)
	_, e := buf.Write(p)
	if e != nil {
		fmt.Println("Can't write")
		return nil, e
	}
	zr, e := gzip.NewReader(buf)
	if e != nil {
		return nil, e
	}
	rbuf := new(bytes.Buffer)
	_, e = io.Copy(rbuf, zr)
	if e != nil {
		return nil, e
	}
	if e != nil {
		fmt.Println("Bjo")
		return nil, e
	}
	return rbuf.Bytes(), e
}

func ReplaceBytes(fn string, newData []byte) error {
	os.Remove(fn)
	f, e := os.Create(fn)
	_, e = f.Write(newData)
	return e
}