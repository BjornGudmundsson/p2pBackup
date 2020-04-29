package files

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)


type fileDataPair struct {
	f File
	d []byte
}

func ReconstructBackup(p []byte, dir string) error {
	d, e := decompressData(p, compress)
	if e != nil {
		return e
	}
	eol := "\n"
	data := string(d)
	files := make([]fileDataPair, 0)
	for len(data) != 0 {
		headerIndex := strings.Index(data, eol)
		if headerIndex == -1 {
			break
		}
		header := data[:headerIndex]
		f, e := fileFromHeader(header)
		if e != nil {
			return e
 		}
 		dataIndex := headerIndex + 1
 		size := f.Size
 		if len(data) < headerIndex + int(size) + 1 {
 			return errors.New("error: could not reconstruct the backup")
		}
 		content := data[dataIndex: headerIndex + int(size) + 1]
 		p := fileDataPair{
			f: f,
			d: []byte(content),
		}
		files = append(files, p)
		data = data[headerIndex + int(size) + 1:]
	}
	for _, p := range files {
		e := ReconstructFile(p.f, p.d, dir)
		if e != nil {
			return e
		}
	}
	return nil
}

func ReconstructFile(f File, content []byte, dir string) error {
	fn, path, fileType := f.Name, f.Path, f.Type
	if fileType == DIR {
		reconstructedDir := dir + "/" + path + "/" + fn
		e := os.MkdirAll(reconstructedDir, os.ModePerm)
		return e
	}
	var reconstructedDir string
	if path != "" {
		reconstructedDir = dir + strings.Replace(path, ".", "", 1)
	} else {
		reconstructedDir = dir
	}
	e := os.MkdirAll(reconstructedDir, os.ModePerm)
	if e != nil {
		return e
	}
	return ioutil.WriteFile(reconstructedDir + "/" + fn, content, 0644)
}

func fileFromHeader(header string) (File, error) {
	file := File{}
	fields := strings.Fields(header)
	if len(fields) != 3 {
		return file, new(ErrorInvalidAmountOfFields)
	}
	nameFields := strings.Split(fields[0], ":")
	if len(nameFields) != 2 {
		return file, new(ErrorInvalidFileFormat)
	}
	file.Name = nameFields[1]
	pathField := strings.Split(fields[1], ":")
	if len(pathField) != 2 {
		return file, new(ErrorInvalidFileFormat)
	}
	file.Path = pathField[1]
	sizeField := strings.Split(fields[2], ":")
	if len(sizeField) != 2 {
		return file, new(ErrorInvalidFileFormat)
	}
	size, e := strconv.Atoi(sizeField[1])
	if e != nil {
		return file, e
	}
	file.Size = int64(size)
	return file, nil
}