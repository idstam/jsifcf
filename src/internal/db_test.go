// Implement tests for db.go

package internal

import (
	"database/sql"
	"log"
	"testing"
)

func TestSqliteDB(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func TestAddSession(t *testing.T) {
	db := SqliteDB{}
	db.Init(":memory:", Logger{})
	id, err := db.AddSession("foobar")
	if err != nil {
		t.Error(err)
	}
	if id == 0 {
		t.Error("Session id is 0")
	}
}
func TestAddHash(t *testing.T) {
	db := SqliteDB{}
	db.Init(":memory:", Logger{})
	session, _ := db.AddSession("foobar")
	id, isDuplicate, err := db.AddHash(session, MD5, "hash")
	if err != nil {
		t.Error(err)
	}
	if id == 0 {
		t.Error("Hash id is 0")
	}
	if isDuplicate {
		t.Error("First insert should not be duplicate")
	}
	id2, isDuplicate, _ := db.AddHash(session, MD5, "hash")
	if id2 != id {
		t.Error("Hash id is not the same", id, id2)
	}
	if !isDuplicate {
		t.Error("Second insert should be duplicate")
	}

}

func TestAddFile(t *testing.T) {
	db := SqliteDB{}
	db.Init(":memory:", Logger{})
	sessionId, err := db.AddSession("foobar")
	if err != nil {
		t.Error(err)
	}
	hashId, _, err := db.AddHash(sessionId, 1, "hash")
	if err != nil {
		t.Error(err)
	}
	id, err := db.AddFile(sessionId, "path", "filename", hashId)
	if err != nil {
		t.Error(err)
	}
	if id == 0 {
		t.Error("File id is 0")
	}
}
func TestInit(t *testing.T) {
	db := SqliteDB{}
	db.Init(":memory:", Logger{})

	tables := []string{}
	rows, err := db.DB.Query("SELECT name FROM sqlite_schema WHERE type='table'")
	for rows.Next() {
		var name string
		rows.Scan(&name)
		tables = append(tables, name)
	}

	if err != nil {
		t.Error(err)
	}
	if !contains(tables, "files") {
		t.Errorf("Files table does not exist")
	}
	if !contains(tables, "sessions") {
		t.Errorf("Sessions table does not exist")
	}
	if !contains(tables, "hashes") {
		t.Errorf("Hashes table does not exist")
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
