package store

import (
	"storeinspection/domain"
	"time"
)

func (d *DB) InspectionsByStore(storeID string) ([]domain.Inspection, error) {
	rows, e := d.SQL.Query(`SELECT id,store_id,inspector,status,summary,opened_at,closed_at FROM inspections WHERE store_id=? ORDER BY opened_at DESC`, storeID)
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
func (d *DB) InspectionsByStatus(status string) ([]domain.Inspection, error) {
	rows, e := d.SQL.Query(`SELECT id,store_id,inspector,status,summary,opened_at,closed_at FROM inspections WHERE status=? ORDER BY opened_at DESC`, status)
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
func (d *DB) FindingsDue(before time.Time) ([]domain.Finding, error) {
	rows, e := d.SQL.Query(`SELECT id,inspection_id,title,severity,status,owner_id,description,due_at,resolution FROM findings WHERE due_at<>'' AND due_at<=? AND status<>? ORDER BY due_at`, before.Format(time.RFC3339), domain.StatusClosed)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []domain.Finding
	for rows.Next() {
		var f domain.Finding
		var due string
		if e = rows.Scan(&f.ID, &f.InspectionID, &f.Title, &f.Severity, &f.Status, &f.OwnerID, &f.Description, &due, &f.Resolution); e != nil {
			return nil, e
		}
		f.DueAt, _ = time.Parse(time.RFC3339, due)
		out = append(out, f)
	}
	return out, rows.Err()
}
func (d *DB) DeleteStore(id string) error {
	_, e := d.SQL.Exec(`DELETE FROM stores WHERE id=?`, id)
	return e
}
