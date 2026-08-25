package store

import (
	"database/sql"
	"storeinspection/domain"
)

type StatusCount struct {
	Status string
	Count  int
}

func (d *DB) CountStatuses() ([]StatusCount, error) {
	rows, e := d.SQL.Query(`SELECT status,COUNT(*) FROM inspections GROUP BY status ORDER BY status`)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []StatusCount
	for rows.Next() {
		var x StatusCount
		if e = rows.Scan(&x.Status, &x.Count); e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (d *DB) CountFindingsBySeverity() map[string]int {
	out := map[string]int{}
	rows, e := d.SQL.Query(`SELECT severity,COUNT(*) FROM findings GROUP BY severity`)
	if e != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var k string
		var n int
		if rows.Scan(&k, &n) == nil {
			out[k] = n
		}
	}
	return out
}
func (d *DB) FindingOwner(id string) (string, error) {
	var owner sql.NullString
	e := d.SQL.QueryRow(`SELECT owner_id FROM findings WHERE id=?`, id).Scan(&owner)
	return owner.String, e
}
func (d *DB) UpdateResolution(id, resolution string) error {
	_, e := d.SQL.Exec(`UPDATE findings SET resolution=? WHERE id=?`, resolution, id)
	return e
}
func (d *DB) MarkApproved(id string, approved bool) error {
	_, e := d.SQL.Exec(`UPDATE remediations SET approved=? WHERE finding_id=?`, approved, id)
	return e
}
func (d *DB) ActiveUsers() ([]domain.User, error) {
	rows, e := d.SQL.Query(`SELECT id,name,role,email,enabled FROM users WHERE enabled=1 ORDER BY name`)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []domain.User
	for rows.Next() {
		var u domain.User
		var en int
		if e = rows.Scan(&u.ID, &u.Name, &u.Role, &u.Email, &en); e != nil {
			return nil, e
		}
		u.Enabled = en != 0
		out = append(out, u)
	}
	return out, rows.Err()
}
