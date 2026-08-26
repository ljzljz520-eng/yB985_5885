package store

import (
	"storeinspection/domain"
	"time"
)

func (d *DB) SaveInspection(i domain.Inspection) error {
	_, e := d.SQL.Exec(`INSERT INTO inspections(id,store_id,inspector,status,summary,opened_at,closed_at) VALUES(?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET status=excluded.status,summary=excluded.summary,closed_at=excluded.closed_at`, i.ID, i.StoreID, i.Inspector, i.Status, i.Summary, i.OpenedAt.Format(time.RFC3339), i.ClosedAt.Format(time.RFC3339))
	return e
}
func (d *DB) GetInspection(id string) (domain.Inspection, error) {
	return scanInspection(d.SQL.QueryRow(`SELECT id,store_id,inspector,status,summary,opened_at,closed_at FROM inspections WHERE id=?`, id))
}
func (d *DB) ListInspections() ([]domain.Inspection, error) {
	rows, e := d.SQL.Query(`SELECT id,store_id,inspector,status,summary,opened_at,closed_at FROM inspections ORDER BY opened_at`)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []domain.Inspection
	for rows.Next() {
		i, e := scanInspection(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, i)
	}
	return out, rows.Err()
}
func (d *DB) DeleteInspection(id string) error {
	tx, e := d.SQL.Begin()
	if e != nil {
		return e
	}
	if _, e = tx.Exec(`DELETE FROM findings WHERE inspection_id=?`, id); e == nil {
		_, e = tx.Exec(`DELETE FROM inspections WHERE id=?`, id)
	}
	if e != nil {
		tx.Rollback()
		return e
	}
	return tx.Commit()
}
