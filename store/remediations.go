package store

import (
	"storeinspection/domain"
	"time"
)

func (d *DB) SaveRemediation(r domain.Remediation) error {
	_, e := d.SQL.Exec(`INSERT INTO remediations(id,finding_id,author,text,attachment,submitted_at,approved) VALUES(?,?,?,?,?,?,?)`, r.ID, r.FindingID, r.Author, r.Text, r.Attachment, r.SubmittedAt.Format(time.RFC3339), r.Approved)
	return e
}
func (d *DB) SaveUser(u domain.User) error {
	_, e := d.SQL.Exec(`INSERT INTO users(id,name,role,email,enabled) VALUES(?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET enabled=excluded.enabled`, u.ID, u.Name, u.Role, u.Email, u.Enabled)
	return e
}
