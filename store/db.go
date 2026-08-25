package store

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

type DB struct{ SQL *sql.DB }

func Open(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	d := &DB{SQL: db}
	if err = d.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return d, nil
}
func (d *DB) Close() error {
	if d == nil || d.SQL == nil {
		return nil
	}
	return d.SQL.Close()
}
func (d *DB) initialize() error {
	stmts := []string{`CREATE TABLE IF NOT EXISTS stores(id TEXT PRIMARY KEY,name TEXT,region TEXT,manager TEXT,active INTEGER,created_at TEXT)`, `CREATE TABLE IF NOT EXISTS inspections(id TEXT PRIMARY KEY,store_id TEXT,inspector TEXT,status TEXT,summary TEXT,opened_at TEXT,closed_at TEXT)`, `CREATE TABLE IF NOT EXISTS findings(id TEXT PRIMARY KEY,inspection_id TEXT,title TEXT,severity TEXT,status TEXT,owner_id TEXT,description TEXT,due_at TEXT,resolution TEXT)`, `CREATE TABLE IF NOT EXISTS assignments(id TEXT PRIMARY KEY,finding_id TEXT,owner_id TEXT,note TEXT,assigned_at TEXT,accepted INTEGER)`, `CREATE TABLE IF NOT EXISTS remediations(id TEXT PRIMARY KEY,finding_id TEXT,author TEXT,text TEXT,attachment TEXT,submitted_at TEXT,approved INTEGER)`, `CREATE TABLE IF NOT EXISTS users(id TEXT PRIMARY KEY,name TEXT,role TEXT,email TEXT,enabled INTEGER)`}
	for _, s := range stmts {
		if _, err := d.SQL.Exec(s); err != nil {
			return fmt.Errorf("schema: %w", err)
		}
	}
	return nil
}
