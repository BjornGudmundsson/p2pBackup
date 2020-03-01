package files

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testFiles []File

func init() {
	f1 := File{
		Name: "f1",
		Type: "file",
		Path: "./test",
	}
	f2, f3 := f1, f1
	f2.Name = "f2"
	f3.Name = "f3"
	testFiles = []File{f1, f2, f3}
}

const fn = "files.go"

func CheckInArray(t *testing.T, a []File, f File) {
	check := false
	for _, fp := range a {
		check = check || fp.Equal(f)
	}
	assert.True(t, check)
}

func TestExists(t *testing.T) {
	assert.True(t, Exists(fn), "This file should 'exist'")
	assert.False(t, Exists("garbage"), "This file should not exist")
}

func TestReadingDir(t *testing.T) {
	files, e := GetFilesFromDir("./test")
	assert.Nil(t, e, "Should be able to read the test directory")
	for i, f := range files {
		assert.True(t, f.Equal(testFiles[i]), "Some file found did not match the expected value")
	}
}

func TestTraversingDir(t *testing.T) {
	t1 := File{
		Name: "t1",
		Type: FILE,
		Path: "./test2",
	}
	t2 := File{
		Name: "t2",
		Type: FILE,
		Path: "./test2/test2",
	}
	t3 := t2
	t3.Name = "t3"
	ts := []File{t1, t2, t3}
	files, e := TraverseDirForFiles("./test2")
	assert.Nil(t, e)
	for _, f := range ts {
		fmt.Println(f)
		CheckInArray(t, files, f)
	}
}
