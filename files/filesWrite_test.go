package files

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const writeFile = "writeFile.txt"

func TestWrite(t *testing.T) {
	f, e := os.Open("writeFile.txt")
	defer f.Close()
	assert.Nil(t, e, "There should be no error when opening a file")
	inf, e2 := f.Stat()
	assert.Nil(t, e2, "There should be no error when getting the file info")
	file := NewFile(inf, ".")
	e = AppendToFile(file, []byte("deadbeef lmao"))
	assert.Nil(t, e, "Should be able to append to the file")
	buf := make([]byte, len([]byte("deadbeef lmao")))
	_, e = f.Read(buf)
	assert.Nil(t, e, "Should be able to read from the file")
	s := string(buf)
	assert.Equal(t, s, "deadbeef lmao")
}

func TestAppend(t *testing.T) {
	f, e := os.Open("appendFile.txt")
	defer f.Close()
	assert.Nil(t, e, "There should be no error when opening a file")
	inf, e2 := f.Stat()
	assert.Nil(t, e2, "There should be no error when getting the file info")
	file := NewFile(inf, ".")
	e = AppendToFile(file, []byte("deadbeef lmao \n"))
	updatedInfo, einf := f.Stat()
	assert.Nil(t, einf, "Getting info again failed")
	fileUpdated := NewFile(updatedInfo, ".")
	e = AppendToFile(fileUpdated, []byte("More beef lmao"))
	assert.Nil(t, e, "Should be able to append to the file after having written to it")
	buf := make([]byte, len([]byte("deadbeef lmao \nMore beef lmao")))
	_, e = f.Read(buf)
	assert.Nil(t, e, "Should be able to read from the file")
	s := string(buf)
	assert.Equal(t, s, "deadbeef lmao \nMore beef lmao")
}
