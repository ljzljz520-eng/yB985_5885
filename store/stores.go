package store

import (
	"storeinspection/domain"
	"time"
)

func (d *DB) SaveStore(s domain.Store) error {
	_, e := d.SQL.Exec(`INSERT INTO stores(id,name,region,manager,active,created_at) VALUES(?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET name=excluded.name,region=excluded.region,manager=excluded.manager,active=excluded.active`, s.ID, s.Name, s.Region, s.Manager, s.Active, s.CreatedAt.Format(time.RFC3339))
	return e
}
func (d *DB) GetStore(id string) (domain.Store, error) {
	var s domain.Store
	var active int
	var created string
	e := d.SQL.QueryRow(`SELECT id,name,region,manager,active,created_at FROM stores WHERE id=?`, id).Scan(&s.ID, &s.Name, &s.Region, &s.Manager, &active, &created)
	s.Active = active != 0
	s.CreatedAt, _ = time.Parse(time.RFC3339, created)
	return s, e
}
func (d *DB) ListStores() ([]domain.Store, error) {
	rows, e := d.SQL.Query(`SELECT id,name,region,manager,active,created_at FROM stores ORDER BY id`)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []domain.Store
	for rows.Next() {
		var s domain.Store
		var a int
		var c string
		if e = rows.Scan(&s.ID, &s.Name, &s.Region, &s.Manager, &a, &c); e != nil {
			return nil, e
		}
		s.Active = a != 0
		s.CreatedAt, _ = time.Parse(time.RFC3339, c)
		out = append(out, s)
	}
	return out, rows.Err()
}
func scanInspection(r interface{ Scan(...any) error }) (domain.Inspection, error) {
	var i domain.Inspection
	var c, o string
	e := r.Scan(&i.ID, &i.StoreID, &i.Inspector, &i.Status, &i.Summary, &o, &c)
	i.OpenedAt, _ = time.Parse(time.RFC3339, o)
	if c != "" {
		i.ClosedAt, _ = time.Parse(time.RFC3339, c)
	}
	return i, e
}
