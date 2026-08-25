package store

import (
	"storeinspection/domain"
	"time"
)

func (d *DB) SaveFinding(f domain.Finding) error {
	_, e := d.SQL.Exec(`INSERT INTO findings(id,inspection_id,title,severity,status,owner_id,description,due_at,resolution) VALUES(?,?,?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET status=excluded.status,owner_id=excluded.owner_id,resolution=excluded.resolution`, f.ID, f.InspectionID, f.Title, f.Severity, f.Status, f.OwnerID, f.Description, f.DueAt.Format(time.RFC3339), f.Resolution)
	return e
}
func (d *DB) ListFindings(inspectionID string) ([]domain.Finding, error) {
	rows, e := d.SQL.Query(`SELECT id,inspection_id,title,severity,status,owner_id,description,due_at,resolution FROM findings WHERE inspection_id=? ORDER BY id`, inspectionID)
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
func (d *DB) SaveAssignment(a domain.Assignment) error {
	_, e := d.SQL.Exec(`INSERT INTO assignments(id,finding_id,owner_id,note,assigned_at,accepted) VALUES(?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET owner_id=excluded.owner_id,accepted=excluded.accepted`, a.ID, a.FindingID, a.OwnerID, a.Note, a.AssignedAt.Format(time.RFC3339), a.Accepted)
	return e
}
