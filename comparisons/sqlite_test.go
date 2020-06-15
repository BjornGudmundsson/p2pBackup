package comparisons

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
	"time"
)

const iterations = 100
const dbString = "yo.db"

type tuple struct {
	size int64
	duration int64
	read int64
}

func (t tuple) ToSTring() []string {
	return []string{strconv.FormatInt(t.size, 10), strconv.FormatInt(t.duration, 10), strconv.FormatInt(t.read, 10)}
}

func TestAddToDB(t *testing.T) {
	fmt.Println("Test")
	db, e := sql.Open("sqlite3", dbString)
	if e != nil {
		fmt.Println(e)
		t.FailNow()
	}
	tuples := make([][]string, 0)
	for i := 0; i < iterations;i++ {
		d := make([]byte, 13)
		_, e = rand.Read(d)
		assert.Nil(t, e)
		key := sha256.Sum256([]byte(d))
		now := time.Now().Nanosecond()
		e := AddToDB(hex.EncodeToString(key[:]), hex.EncodeToString(d), db)
		end := time.Now().Nanosecond()
		if e != nil {
			fmt.Println(e)
			continue
		}
		nowRead := time.Now().Nanosecond()
		e = QueryDB(hex.EncodeToString(key[:]), db)
		endRead := time.Now().Nanosecond()
		if e != nil {
			fmt.Println(e)
		}
		tup := tuple{
			size: int64(13 * i),
			duration: int64(end - now),
			read: int64(endRead - nowRead),
		}
		if tup.duration > 0 {
			tuples = append(tuples, tup.ToSTring())
		}
		assert.Nil(t, e)
	}
	f, e := os.Create("results.csv")
	if e != nil {
		fmt.Println(e)
		t.FailNow()
		return
	}
	writer := csv.NewWriter(f)
	defer writer.Flush()
	for _, s := range tuples {
		writer.Write(s)
	}
}
