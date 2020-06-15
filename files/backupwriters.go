package files

import (
	aes2 "crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/BjornGudmundsson/p2pBackup/crypto"
	"os"
)

type AppendEmptyBuffer struct {
	fn string
}

func NewAppendEmptyBuffer(fn string) BackupHandler {
	return &AppendEmptyBuffer{
		fn: fn,
	}
}

func (bb *AppendEmptyBuffer) AddBackup(d []byte) int64 {
	fn := bb.fn
	f, e := GetFile(fn)
	if e != nil {
		fmt.Println(e)
		return 0
	}
	e = AppendToFile(*f, d)
	if e != nil {
		return 0
	}
	size := f.Size
	return size
}

func (bb *AppendEmptyBuffer) ReadFrom(start, size int64) ([]byte, error) {
	fn := bb.fn
	f, e := os.Open(fn)
	defer f.Close()
	if e != nil {
		return nil, e
	}
	d := make([]byte, size)
	_, e = f.ReadAt(d, start)
	if e != nil {
		return nil, e
	}
	return d, nil
}

type AppendBufferFull struct {
	key []byte
	fn string
}

func NewAppendBufferFull(fn string, key []byte) BackupHandler {
	return &AppendBufferFull{
		key: key,
		fn: fn,
	}
}

func (bb *AppendBufferFull) AddBackup(d []byte) int64 {
	fn := bb.fn
	key := bb.key
	f, e := os.OpenFile(fn, os.O_RDWR, os.ModeAppend)
	if e != nil {
		return 0
	}
	written := make([]byte, KEYLEN)
	_, e = f.ReadAt(written, 0)
	block, e := aes2.NewCipher(key)
	if e != nil {
		return 0
	}
	block.Decrypt(written, written)
	part1, part2 := written[0:KEYLEN / 2], written[KEYLEN / 2: KEYLEN]
	size := int64(KEYLEN)
	if string(part1) == string(part2) {
		size = int64(binary.LittleEndian.Uint64(part1))
	}
	f.WriteAt(d, size)
	update := uint64(size) + uint64(len(d))
	updatedWritten := make([]byte, KEYLEN)
	binary.LittleEndian.PutUint64(updatedWritten, update)
	binary.LittleEndian.PutUint64(updatedWritten[KEYLEN / 2:], update)
	block.Encrypt(updatedWritten, updatedWritten)
	_, e = f.WriteAt(updatedWritten, 0)
	if e != nil {
		return 0
	}
	return size
}

func (bb *AppendBufferFull) ReadFrom(start, size int64) ([]byte, error) {fn := bb.fn
	f, e := os.Open(fn)
	if e != nil {
		return nil, e
	}
	d := make([]byte, size)
	_, e = f.ReadAt(d, start)
	if e != nil {
		return nil, e
	}
	return d, nil
}

type SegmentedBuffer struct {
	segmentSize int
	key []byte
	fn string
}

func NewSegmentBuffer(fn string, key []byte, segmentSize int) BackupHandler {
	return &SegmentedBuffer{
		segmentSize: segmentSize,
		key:         key,
		fn:          fn,
	}
}

func (bb *SegmentedBuffer) getWrittenBytes() int64 {
	fn, key := bb.fn, bb.key
	f, e := os.OpenFile(fn, os.O_RDWR, os.ModeAppend)
	defer f.Close()
	if e != nil {
		return 0
	}
	nonce := make([]byte, KEYLEN)
	_, e = f.ReadAt(nonce, 0)
	if e != nil {
		return 0
	}
	written := make([]byte, KEYLEN)
	_, e = f.ReadAt(written, KEYLEN)
	if e != nil {
		return 0
	}
	block, e := aes2.NewCipher(key)
	if e != nil {
		return 0
	}
	size := int64(2 * KEYLEN)
	ctr := cipher.NewCTR(block, nonce)
	ctr.XORKeyStream(written, written)
	part1, part2  := written[0: KEYLEN / 2], written[KEYLEN / 2: KEYLEN]
	if string(part1) == string(part2) {
		size = int64(binary.LittleEndian.Uint64(part1))
	}
	return size
}

func (bb *SegmentedBuffer) writeSegment(written int64, d []byte,f *os.File) (int64, error) {
	k := int64((bb.segmentSize + 1) * KEYLEN)
	div, mod := written / k, written % k
	//fmt.Println("Div mod: ", div, mod)
	ind := div * k
	segment := make([]byte, k)
	_, e := f.ReadAt(segment, ind)
	if e != nil {
		return 0, e
	}
	//fmt.Println("Segment: ", string(segment))
	nonce := segment[:KEYLEN]
	newNonce := make([]byte, KEYLEN)
	data := segment[KEYLEN:]
	e = crypto.EncryptCTR(nonce, bb.key, data)
	writtenInSegment := mod
	leftInSegment := k - writtenInSegment
	//fmt.Println("Left in segment: ", leftInSegment)
	if int64(len(d)) < leftInSegment {
		copy(segment[writtenInSegment:], d)
		e = crypto.EncryptCTR(newNonce, bb.key, data)
		copy(segment, newNonce)
		_, e = f.WriteAt(segment, ind)
		if e != nil {
			return -1, e
		}
		return int64(len(d)), nil
	}
	//fmt.Println("Len: ", len(d), len(data))
	d = d[:leftInSegment]
	copy(data[writtenInSegment - KEYLEN:], d)
	e = crypto.EncryptCTR(newNonce, bb.key, data)
	copy(segment, newNonce)
	_, e = f.WriteAt(segment, ind)
	//fmt.Println("Returning")
	buffer := make([]byte, k - KEYLEN)
	f.ReadAt(buffer, KEYLEN)
	crypto.EncryptCTR(newNonce, bb.key, buffer)
	return leftInSegment, e
}

func (bb *SegmentedBuffer) updateWrittenBytes(written int64, f *os.File) error {
	writtenSlice := make([]byte, KEYLEN)
	key := bb.key
	k := bb.segmentSize + 1
	binary.LittleEndian.PutUint64(writtenSlice, uint64(written))
	binary.LittleEndian.PutUint64(writtenSlice[KEYLEN / 2:], uint64(written))
	segment := make([]byte, KEYLEN * k)
	_, e := f.ReadAt(segment, 0)
	if e != nil {
		return e
	}
	nonce := segment[:KEYLEN]
	data := segment[KEYLEN:]
	block, e := aes2.NewCipher(key)
	if e != nil {
		return e
	}
	ctr := cipher.NewCTR(block, nonce)
	ctr.XORKeyStream(data, data)
	copy(data, writtenSlice)
	newNonce := make([]byte, KEYLEN)
	_, e = rand.Read(newNonce)
	if e != nil {
		return e
	}
	ctr2 := cipher.NewCTR(block, newNonce)
	ctr2.XORKeyStream(data, data)
	data = append(newNonce, data...)
	_, e = f.WriteAt(data, 0)
	return e
}

func (bb *SegmentedBuffer) AddBackup(d []byte) int64 {
	written := bb.getWrittenBytes()
	k := int64((bb.segmentSize + 1) * KEYLEN)
	f, e := os.OpenFile(bb.fn, os.O_RDWR, os.ModeAppend)
	if e != nil {
		return -1
	}
	defer f.Close()
	temp := written
	for d != nil && len(d) != 0 {
		n, e := bb.writeSegment(written, d, f)
		if e != nil {
			return -1
		}
		d = d[n:]
		written += n
		if written % k == 0 {
			written += KEYLEN
		}
	}
	e = bb.updateWrittenBytes(written, f)
	if e != nil {
		return -1
	}
	return temp
}

func (bb *SegmentedBuffer) ReadFrom(start, size int64) ([]byte, error) {
	k := int64(bb.segmentSize + 1) * KEYLEN
	div := start / k
	ind := div * k
	f, e := os.OpenFile(bb.fn, os.O_RDWR, os.ModeAppend)
	defer f.Close()
	if e != nil {
		return nil, e
	}
	tempSize := size
	readData := make([]byte, size)
	firstSegment := make([]byte, k)
	_, e = f.ReadAt(firstSegment, ind)
	if e != nil {
		return nil, e
	}
	mod := start % k
	nonce := firstSegment[:KEYLEN]
	data := firstSegment[KEYLEN:]
	e = crypto.EncryptCTR(nonce, bb.key, data)
	if e != nil {
		return nil, e
	}
	copy(readData, firstSegment[mod:])
	readBits := k - mod
	tempSize -= (k - mod)
	for tempSize > 0 {
		newStart := ind + k
		nextSegment := make([]byte, k)
		_, e := f.ReadAt(nextSegment, newStart)
		if e != nil {
			return nil, e
		}
		nonce := nextSegment[:KEYLEN]
		d := nextSegment[KEYLEN:]
		e = crypto.EncryptCTR(nonce, bb.key, d)
		if e != nil {
			return nil, e
		}
		if tempSize < int64(bb.segmentSize * KEYLEN) {
			copy(readData[readBits:], d[:tempSize])
			tempSize = 0
		} else {
			copy(readData[readBits:], d)
			tempSize -= int64(bb.segmentSize * KEYLEN)
			readBits += int64(bb.segmentSize * KEYLEN)
		}
	}
	return readData, nil
}