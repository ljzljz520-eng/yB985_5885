package store

import (
	"database/sql"
	"fmt"
	"time"
)

type Health struct {
	Stores, Inspections, Findings int
	Database                      string
}

func (d *DB) Health() Health {
	h := Health{Database: "sqlite"}
	d.SQL.QueryRow(`SELECT COUNT(*) FROM stores`).Scan(&h.Stores)
	d.SQL.QueryRow(`SELECT COUNT(*) FROM inspections`).Scan(&h.Inspections)
	d.SQL.QueryRow(`SELECT COUNT(*) FROM findings`).Scan(&h.Findings)
	return h
}
func (d *DB) Vacuum() error { _, e := d.SQL.Exec(`VACUUM`); return e }
func (d *DB) PurgeEvents(before time.Time) (int64, error) {
	r, e := d.SQL.Exec(`DELETE FROM events WHERE at<?`, before.Format(time.RFC3339))
	if e != nil {
		return 0, e
	}
	return r.RowsAffected()
}
func (d *DB) Ping() error {
	if d == nil || d.SQL == nil {
		return fmt.Errorf("database unavailable")
	}
	return d.SQL.Ping()
}
func (d *DB) Transaction(fn func(*sql.Tx) error) error {
	tx, e := d.SQL.Begin()
	if e != nil {
		return e
	}
	if e = fn(tx); e != nil {
		tx.Rollback()
		return e
	}
	return tx.Commit()
}
func (d *DB) ResetFindings(inspection string) error {
	return d.Transaction(func(tx *sql.Tx) error {
		_, e := tx.Exec(`UPDATE findings SET status=?,resolution='' WHERE inspection_id=?`, "open", inspection)
		return e
	})
}
