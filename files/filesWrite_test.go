package files

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

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

func TestModified(t *testing.T) {
	f, e := os.Open("appendFile.txt")
	defer f.Close()
	assert.Nil(t, e, "There should be no error when opening a file")
	inf, e2 := f.Stat()
	assert.Nil(t, e2, "There should be no error when getting the file info")
	file := NewFile(inf, ".")
	e = AppendToFile(file, []byte("deadbeef lmao \n"))
	assert.Nil(t, e, "Should be able to append to the file after having written to it")
	inf, e2 = f.Stat()
	assert.Nil(t, e2, "Getting info again failed")
	fupdated := NewFile(inf, ".")
	time.Sleep(time.Second)
	elapsed := GetTimePassedSinceModified(fupdated)
	assert.True(t, elapsed >= time.Second, "At least one second should be measured to have passed")
}

func TestTomlRead(t *testing.T) {
	bd := CreateRules("./test.toml")
	assert.Equal(t, 2*time.Second, bd.GetMinTime(), "The toml should have 1 second in the modify time")
	assert.Equal(t, int64(500), bd.MaxSize, "500 bytes should be the maximum size allowed")
	assert.Equal(t, int64(0), bd.MinSize, "0 should be the minimum size")
	assert.Equal(t, "([a-z]*).csv", bd.TypesToExclude, "It should exclude all .txt files")

	f, e := os.Open("appendFile.txt")
	defer f.Close()
	assert.Nil(t, e, "There should be no error when opening a file")
	inf, e2 := f.Stat()
	assert.Nil(t, e2, "There should be no error when getting the file info")
	file := NewFile(inf, ".")
	e = AppendToFile(file, []byte("deadbeef lmao \n"))
	assert.Nil(t, e, "Should be able to append to the file after having written to it")
	time.Sleep(time.Second)
	inf, e2 = f.Stat()
	assert.Nil(t, e2, "There should be no error when getting the file info 2")
	file = NewFile(inf, ".")
	assert.False(t, bd.Include(file))
	time.Sleep(time.Second)
	inf, e2 = f.Stat()
	assert.Nil(t, e2, "There should be no error when getting the file info 3")
	file = NewFile(inf, ".")
	assert.True(t, bd.Include(file), "Enough time should have elapsed in order to get the file")

	defaultRules := CreateRules("garbage")
	assert.Equal(t, int64(1000), defaultRules.MaxSize, "Default max is 1000")
	assert.Equal(t, int64(0), defaultRules.MinSize, "Default min is 0")
	assert.Equal(t, "([a-z]*).csv",defaultRules.TypesToExclude, "Default has no csv")
	assert.Equal(t, 0, len(defaultRules.BlackListedFiles), "Empty by default")
	assert.Equal(t, time.Second, defaultRules.GetMinTime(), "Default modify time is 1 second")
}

func TestReadLastBytes(t *testing.T) {
	f, e := os.Open("appendFile.txt")
	defer f.Close()
	assert.Nil(t, e, "There should be no error when opening a file")
	inf, e2 := f.Stat()
	assert.Nil(t, e2, "There should be no error when getting the file info")
	file := NewFile(inf, ".")
	e = AppendToFile(file, []byte("My name is Bjorn"))
	assert.Nil(t, e, "Should be able to append to the file after having written to it")
	d, e := GetLastNBytes("appendFile.txt", 5)
	assert.Nil(t, e, "Should be able to read the last n bytes of the file")
	assert.Equal(t, string(d), "Bjorn", "The last 5 bytes should be 'Bjorn'")
}

func TestAppendLog(t *testing.T) {
	logFile := "logFile.txt"
	pw := "deadbeef"
	trustee := NewTrustee("Not real", "not a suite")
	location := NewLocation("Bjorn", 1)
	trustees := []Trustee{trustee}
	locations := []Location{location}
	entry := NewLogEntry(time.Now(), "ABCDF", 10, 10, locations, trustees)
	entry2 := NewLogEntry(time.Now(), "ABCDF", 11, 11, locations, trustees)
	entry3 := NewLogEntry(time.Now(), "ABCDF", 12, 12, locations, trustees)
	e := AddBackupLog(entry, logFile, pw)
	assert.Nil(t, e, "Should be able to add a new log entry")
	e = AddBackupLog(entry2, logFile, pw)
	assert.Nil(t, e, "Should be able to add a new log entry")
	e = AddBackupLog(entry3, logFile, pw)
	assert.Nil(t, e, "Should be able to add a new log entry")
	ct, e := ioutil.ReadFile(logFile)
	assert.Nil(t, e, "Should be able to access the content of the file")
	notContain := strings.Contains(string(ct), entry.String())
	assert.False(t, notContain, "The ciphertext should not contain the entry")
	d, e := DecryptBackupLog(logFile, pw)
	assert.Nil(t, e, "Should be able to decrypt it ")
	contain := strings.Contains(string(d), entry.String())
	assert.True(t, contain, "The decrypted file should contain the log")
	contain = strings.Contains(string(d), entry2.String())
	assert.True(t, contain, "The decrypted file should contain the log")
	contain = strings.Contains(string(d), entry3.String())
	assert.True(t, contain, "The decrypted file should contain the log")
	fmt.Println(string(d))
}
