package store

import (
	"storeinspection/domain"
	"time"
)

func (d *DB) EnsureEvents() error {
	_, e := d.SQL.Exec(`CREATE TABLE IF NOT EXISTS events(id TEXT PRIMARY KEY,aggregate_id TEXT,type TEXT,actor TEXT,payload TEXT,at TEXT)`)
	return e
}
func (d *DB) AppendEvent(e domain.Event) error {
	if err := d.EnsureEvents(); err != nil {
		return err
	}
	_, err := d.SQL.Exec(`INSERT INTO events(id,aggregate_id,type,actor,payload,at) VALUES(?,?,?,?,?,?)`, e.ID, e.AggregateID, e.Type, e.Actor, e.Payload, e.At.Format(time.RFC3339))
	return err
}
func (d *DB) Events(aggregate string) ([]domain.Event, error) {
	if err := d.EnsureEvents(); err != nil {
		return nil, err
	}
	rows, e := d.SQL.Query(`SELECT id,aggregate_id,type,actor,payload,at FROM events WHERE aggregate_id=? ORDER BY at,id`, aggregate)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []domain.Event
	for rows.Next() {
		var x domain.Event
		var at string
		if e = rows.Scan(&x.ID, &x.AggregateID, &x.Type, &x.Actor, &x.Payload, &at); e != nil {
			return nil, e
		}
		x.At, _ = time.Parse(time.RFC3339, at)
		out = append(out, x)
	}
	return out, rows.Err()
}
