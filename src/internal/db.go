package internal

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteDB struct {
	DB     *sql.DB
	logger Logger
}

func (me *SqliteDB) Init(database string, logger Logger) {
	var err error
	me.logger = logger

	me.DB, err = sql.Open("sqlite3", database)
	if err != nil {
		logger.Fatal("SqliteDb.Init", err.Error())
	}

	if !me.execRaw(
		`CREATE TABLE IF NOT EXISTS sessions (
			id integer primary key autoincrement,
			computer text,
			started_at text, 
			done_at text,
			file_count integer,
			new_file_count integer
		);`) {
		logger.Fatal("SqliteDb.Init", "Could not init db")
	}

	if !me.execRaw(
		`CREATE TABLE IF NOT EXISTS hashes (
			id integer primary key autoincrement,
			session_id integer, 
			hash_algo integer, 
			hash_value text, 
			found_at text,
			FOREIGN KEY (session_id) REFERENCES sessions(id)
		);`) {
		logger.Fatal("SqliteDb.Init", "Could not init db")
	}
	if !me.execRaw(
		`CREATE TABLE IF NOT EXISTS files (
			id integer primary key autoincrement,
			session_id integer, 
			path text,
			file_name text, 
			hash_id integer,
			FOREIGN KEY (session_id) REFERENCES sessions(id),
			FOREIGN KEY (hash_id) REFERENCES hashes(id)
		);`) {
		logger.Fatal("SqliteDb.Init", "Could not init db")
	}
	if !me.execRaw("CREATE UNIQUE INDEX idx_hash ON hashes(hash_value)") {
		logger.Fatal("SqliteDb.Init", "Could not init db")
	}

}

func (me *SqliteDB) execRaw(sql string) bool {
	_, err := me.DB.Exec(sql)

	if err != nil {
		me.logger.Error("SqliteDb.execRaw", err.Error())
		return false
	}
	return true
}
func (me *SqliteDB) AddSession(computer string) (int64, error) {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	result, err := me.DB.Exec("INSERT INTO sessions (computer, started_at) VALUES (?, ?)", computer, now)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, nil
}
func (me *SqliteDB) AddHash(session int64, algo HashAlgo, hash string) (int64, bool, error) {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	isDuplicate := false
	result, err := me.DB.Exec("INSERT OR IGNORE INTO hashes (session_id, hash_algo, hash_value, found_at) VALUES (?, ?, ?, ?)",
		session, algo, hash, now)

	if err != nil {
		return 0, false, err
	}

	id, err := result.LastInsertId()

	if id == 0 {
		isDuplicate = true
		id, err = me.GetHash(hash)
	}

	if err != nil {
		return 0, false, err
	}

	return id, isDuplicate, nil
}

func (me *SqliteDB) GetHash(hash string) (int64, error) {

	row := me.DB.QueryRow("SELECT id FROM hashes WHERE hash_value = ?", hash)

	var id int64
	err := row.Scan(&id)

	if err != nil {
		me.logger.Error("SqliteDb.GetHash", err.Error())
		return 0, err
	}

	return id, nil
}
func (me *SqliteDB) AddFile(session int64, path string, filename string, hash int64) (int64, error) {
	result, err := me.DB.Exec("INSERT INTO files (session_id, path, file_name, hash_id) VALUES (?, ?, ?, ?)",
		session, path, filename, hash)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (me *SqliteDB) Get(connectionId string, key string) (string, bool, error) {

	row := me.DB.QueryRow("SELECT cache_value FROM cache_table WHERE connection_id = ? and cache_key = ?", connectionId, key)

	val := ""
	err := row.Scan(&val)
	if err != nil {
		return "", false, err
	}

	return val, true, nil

}
func (me *SqliteDB) Set(connectionId string, key string, value string) {

	stm, err := me.DB.Prepare("INSERT INTO cache_table (connection_id, cache_key, cache_value) VALUES(?, ?, ?) ON CONFLICT(connection_id, cache_key) DO UPDATE SET cache_value=excluded.cache_value;")

	if err != nil {
		log.Fatal(err)
	}

	defer stm.Close()

	_, err = stm.Exec(connectionId, key, value)
	if err != nil {
		log.Fatal(err)
	}
}
