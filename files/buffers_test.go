package files

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testBuffer(t *testing.T, bb BackupHandler) {
	data := []byte("deadbeef lmao")
	size := int64(len(data))
	start1 := bb.AddBackup(data)
	start2 := bb.AddBackup(data)
	d1, e := bb.ReadFrom(start1, size)
	if e != nil {
		t.FailNow()
		return
	}
	assert.Equal(t, string(data), string(d1))
	d2, e := bb.ReadFrom(start2, size)
	if e != nil {
		fmt.Println(e)
	}
	assert.Equal(t, string(data), string(d2))
	fmt.Println("Starting")
	start3 := bb.AddBackup(data)
	fmt.Println("Ending")
	d3, e := bb.ReadFrom(start3, size)
	assert.Equal(t, string(data), string(d3))
}
/*
func TestAppendEmptyBuffer(t *testing.T) {
	fn := "appendOnly.txt"
	bb := NewAppendEmptyBuffer(fn)
	testBuffer(t, bb)
}

func TestAppendFullBuffer(t *testing.T) {
	key := []byte("YELLOW SUBMARINE")
	fn := "appendFull.txt"
	bb := NewAppendBufferFull(fn, key)
	testBuffer(t, bb)
}

func TestSegmentBuffer(t *testing.T) {
	segmentSize := 5
	key := []byte("YELLOW SUBMARINE")
	fn := "segmentBuffer.txt"
	bb := NewSegmentBuffer(fn, key, segmentSize)
	testBuffer(t, bb)
	f, _ := os.Open(fn)
	info, _ := f.Stat()
	buffer := make([]byte, info.Size())
	f.ReadAt(buffer, 0)
	fmt.Println("Buffer: ", string(buffer))
	time.Sleep(30*time.Second)
}
*/