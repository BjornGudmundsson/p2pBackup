package comparisons

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)
var insert = "INSERT INTO backups(BackupId, Backup) VALUES(?, ?)"
var query = "select Backup FROM backups WHERE BackupId = ?"
func AddToDB(key, data string, db *sql.DB) error {
	_, e := db.Exec(insert, key, data)
	return e
}

func QueryDB(key string, db *sql.DB) error {
	_, e := db.Exec(query, key)
	return e
}
