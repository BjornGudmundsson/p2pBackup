package files

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"sort"
	"strconv"
	"testing"
	"time"
)

func testPerformanceBuffer(t *testing.T, bb BackupHandler) {
	numTests := 100
	data := "deadbeef lmao"
	for i := 0; i < numTests;i++ {
		size := len(data)
		now := time.Now()
		bb.AddBackup([]byte(data))
		end := time.Now()
		fmt.Println("Elapsed: ", end.Nanosecond() - now.Nanosecond(), " Size: ", size)
		data += "deadbeef lmao"
	}
}
/*
func TestAppendTimer(t *testing.T) {
	fmt.Println("Append empty timer")
	fn := "appendOnly.txt"
	bb := NewAppendEmptyBuffer(fn)
	testPerformanceBuffer(t, bb)
}

func TestAppendFullTimer(t *testing.T) {
	fmt.Println("Append full timer")
	key := []byte("YELLOW SUBMARINE")
	fn := "appendFull.txt"
	bb := NewAppendBufferFull(fn, key)
	testPerformanceBuffer(t, bb)
}

func TestSegmentBufferTimer(t *testing.T) {
	fmt.Println("Segment Buffer")
	segmentSize := 5
	key := []byte("YELLOW SUBMARINE")
	fn := "segmentBuffer.txt"
	bb := NewSegmentBuffer(fn, key, segmentSize)
	testPerformanceBuffer(t, bb)
}*/

type tuple struct {
	size int64
	durations []int
}

func (t tuple) ToString() []string {
	arr := make([]string, 1 + len(t.durations))
	arr[0] = strconv.FormatInt(t.size, 10)
	for i, d := range t.durations {
		arr[i + 1] = strconv.Itoa(d)
	}
	return arr
}


func TestBuffersPerformance(t *testing.T) {
	f1 := "appendPerformance.txt"
	f2 := "appendFullPerformance.txt"
	f3 := "segmentPerformance.txt"
	segmentSize := 5
	fileAppend, e := os.Create(f1)
	assert.Nil(t, e)
	fileAppendFull, e := os.Create(f2)
	assert.Nil(t, e)
	fileSegment, e := os.Create(f3)
	assert.Nil(t, e)
	key := []byte("YELLOW SUBMARINE")
	resultsAppends := make(map[int][]int)
	resultsAppendFull := make(map[int][]int)
	resultsSegment := make(map[int][]int)
	resultsReadAppends := make(map[int][]int)
	resultsReadFull := make(map[int][]int)
	resultsReadSegment := make(map[int][]int)
	data := "deadbeef lmao"
	iterations := 100
	avgCount := 20
	for i := 0; i < avgCount;i++ {
		sum := data
		populateFile(0, fileAppend)
		populateFile(10000, fileAppendFull)
		populateFile(10000, fileSegment)
		bbAppend := NewAppendEmptyBuffer(f1)
		bbAppendFull := NewAppendBufferFull(f2, key)
		bbSegment := NewSegmentBuffer(f3, key, segmentSize)
		for j := 0; j < iterations;j++ {
			s1 := updateTimerMap(resultsAppends, []byte(sum), bbAppend, j)
			s2 := updateTimerMap(resultsAppendFull, []byte(sum), bbAppendFull, j)
			s3 := updateTimerMap(resultsSegment, []byte(sum), bbSegment, j)
			if s1 == -1 || s2 == -1 || s3 == -1 {
				break
			}
			updateReadTimerMap(resultsReadAppends, s1, int64(len(data)), bbAppend, j * len(data))
			updateReadTimerMap(resultsReadFull, s2, int64(len(data)), bbAppendFull, j * len(data))
			updateReadTimerMap(resultsReadSegment, s3, int64(len(data)), bbSegment, j * len(data))
		}

	}
	os.Remove(f1)
	os.Remove(f2)
	os.Remove(f3)
	tuples := make([]tuple, 0)
	tuplesRead := make([]tuple, 0)
	for k, v := range resultsAppends {
		v2, ok2 := resultsAppendFull[k]
		v3, ok3 := resultsSegment[k]
		if ok2 && ok3 {
			tup := tuple{
				size:      int64(k),
				durations: []int{avgSlice(v), avgSlice(v2), avgSlice(v3)},
			}
			tuples = append(tuples, tup)
		}
	}
	for k, v := range resultsReadAppends {
		v2, ok2 := resultsReadFull[k]
		v3, ok3 := resultsReadSegment[k]
		if ok2 && ok3 {
			tup := tuple{
				size:      int64(k),
				durations: []int{avgSlice(v), avgSlice(v2), avgSlice(v3)},
			}
			tuplesRead = append(tuplesRead, tup)
		}
	}
	sort.SliceStable(tuplesRead, func(i, j int) bool {
		return tuplesRead[i].size < tuplesRead[j].size
	})
	sort.SliceStable(tuples, func(i, j int) bool {
		return tuples[i].size < tuples[j].size
	})
	c := "results.csv"
	c2 := "resultsRead.csv"
	r, _ := os.Create(c)
	r2, _ := os.Create(c2)
	writer := csv.NewWriter(r)
	readerWriter := csv.NewWriter(r2)
	defer writer.Flush()
	defer readerWriter.Flush()
	for _, tup := range tuples {
		writer.Write(tup.ToString())
	}
	for _, tup := range tuplesRead {
		readerWriter.Write(tup.ToString())
	}
}

func avgSlice(a []int) int {
	sort.Ints(a)
	l := len(a)
	return a[l / 2]
}

func updateTimerMap(m map[int][]int, data []byte, bb BackupHandler, ind int) int64 {
	size := len(data) * ind
	now := time.Now().Nanosecond()
	start1 := bb.AddBackup(data)
	if start1 == -1 {
		return -1
	}
	end := time.Now().Nanosecond()
	elapsed := end - now
	if timers, ok := m[size]; ok {
		timers = append(timers, elapsed)
		m[size] = timers
	} else {
		m[size] = []int{elapsed}
	}
	return start1
}

func updateReadTimerMap(m map[int][]int, start, size int64, bb BackupHandler, ind int) {
	now := time.Now().Nanosecond()
	_, e := bb.ReadFrom(start, size)
	end := time.Now().Nanosecond()
	if e != nil {
		fmt.Println(e)
		return
	}
	duration := end - now
	m[ind] = append(m[ind], duration)
}

func populateFile(n int, f *os.File) error {
	buffer := make([]byte, n)
	_, e := rand.Read(buffer)
	if e != nil {
		return e
	}
	_, e = f.WriteAt(buffer, 0)
	return e
}