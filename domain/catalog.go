package domain

import (
	"sort"
	"strings"
	"time"
)

type Catalog struct {
	Stores     []Store
	Users      []User
	Checklists []Checklist
}

func (c *Catalog) AddStore(s Store) bool {
	if ValidateStore(s) != nil {
		return false
	}
	for _, x := range c.Stores {
		if x.ID == s.ID {
			return false
		}
	}
	c.Stores = append(c.Stores, s)
	return true
}
func (c *Catalog) UpdateStore(s Store) bool {
	for i := range c.Stores {
		if c.Stores[i].ID == s.ID {
			c.Stores[i] = s
			return true
		}
	}
	return false
}
func (c *Catalog) RemoveStore(id string) bool {
	for i := range c.Stores {
		if c.Stores[i].ID == id {
			c.Stores = append(c.Stores[:i], c.Stores[i+1:]...)
			return true
		}
	}
	return false
}
func (c Catalog) FindStore(id string) (Store, bool) {
	for _, s := range c.Stores {
		if s.ID == id {
			return s, true
		}
	}
	return Store{}, false
}
func (c Catalog) SearchStores(term string) []Store {
	term = strings.ToLower(term)
	var out []Store
	for _, s := range c.Stores {
		if term == "" || strings.Contains(strings.ToLower(s.Name), term) || strings.Contains(strings.ToLower(s.Region), term) {
			out = append(out, s)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
func (c *Catalog) AddUser(u User) bool {
	if u.ID == "" || u.Name == "" {
		return false
	}
	for _, x := range c.Users {
		if x.ID == u.ID {
			return false
		}
	}
	c.Users = append(c.Users, u)
	return true
}
func (c *Catalog) EnableUser(id string) bool {
	for i := range c.Users {
		if c.Users[i].ID == id {
			c.Users[i].Enabled = true
			return true
		}
	}
	return false
}
func (c *Catalog) DisableUser(id string) bool {
	for i := range c.Users {
		if c.Users[i].ID == id {
			c.Users[i].Enabled = false
			return true
		}
	}
	return false
}
func (c Catalog) FindUser(id string) (User, bool) {
	for _, u := range c.Users {
		if u.ID == id {
			return u, true
		}
	}
	return User{}, false
}
func (c Catalog) UsersByRole(role string) []User {
	role = NormalizeRole(role)
	var out []User
	for _, u := range c.Users {
		if NormalizeRole(u.Role) == role && u.Enabled {
			out = append(out, u)
		}
	}
	return out
}
func (c *Catalog) AddChecklist(ch Checklist) bool {
	if ch.Validate() != nil {
		return false
	}
	for _, x := range c.Checklists {
		if x.ID == ch.ID && x.Version == ch.Version {
			return false
		}
	}
	c.Checklists = append(c.Checklists, ch)
	return true
}
func (c Catalog) LatestChecklist(id string) (Checklist, bool) {
	var best Checklist
	found := false
	for _, x := range c.Checklists {
		if x.ID == id {
			if !found || x.Version > best.Version {
				best = x
				found = true
			}
		}
	}
	return best, found
}
func (c Catalog) ActiveStores() []Store {
	var out []Store
	for _, s := range c.Stores {
		if s.Active {
			out = append(out, s)
		}
	}
	return out
}
func (c Catalog) InactiveStores() []Store {
	var out []Store
	for _, s := range c.Stores {
		if !s.Active {
			out = append(out, s)
		}
	}
	return out
}
func (c Catalog) StoreCount() int     { return len(c.Stores) }
func (c Catalog) UserCount() int      { return len(c.Users) }
func (c Catalog) ChecklistCount() int { return len(c.Checklists) }
func (c Catalog) CreatedAfter(at time.Time) []Store {
	var out []Store
	for _, s := range c.Stores {
		if s.CreatedAt.After(at) {
			out = append(out, s)
		}
	}
	return out
}
func (c Catalog) Regions() []string {
	m := map[string]bool{}
	for _, s := range c.Stores {
		m[s.Region] = true
	}
	var out []string
	for r := range m {
		out = append(out, r)
	}
	sort.Strings(out)
	return out
}
func (c Catalog) Managers() []string {
	m := map[string]bool{}
	for _, s := range c.Stores {
		if s.Manager != "" {
			m[s.Manager] = true
		}
	}
	var out []string
	for x := range m {
		out = append(out, x)
	}
	sort.Strings(out)
	return out
}
func (c Catalog) EnabledUsers() []User {
	var out []User
	for _, u := range c.Users {
		if u.Enabled {
			out = append(out, u)
		}
	}
	return out
}
func (c Catalog) HasStore(id string) bool { _, ok := c.FindStore(id); return ok }
func (c Catalog) HasUser(id string) bool  { _, ok := c.FindUser(id); return ok }
func (c Catalog) HasChecklist(id, version string) bool {
	for _, x := range c.Checklists {
		if x.ID == id && x.Version == version {
			return true
		}
	}
	return false
}
func (c *Catalog) RenameStore(id, name string) bool {
	for i := range c.Stores {
		if c.Stores[i].ID == id && strings.TrimSpace(name) != "" {
			c.Stores[i].Name = name
			return true
		}
	}
	return false
}
func (c *Catalog) MoveStore(id, region string) bool {
	for i := range c.Stores {
		if c.Stores[i].ID == id && region != "" {
			c.Stores[i].Region = region
			return true
		}
	}
	return false
}
func (c *Catalog) AssignManager(id, manager string) bool {
	for i := range c.Stores {
		if c.Stores[i].ID == id {
			c.Stores[i].Manager = manager
			return true
		}
	}
	return false
}
func (c Catalog) RoleNames() []string {
	m := map[string]bool{}
	for _, u := range c.Users {
		m[NormalizeRole(u.Role)] = true
	}
	var out []string
	for r := range m {
		out = append(out, r)
	}
	sort.Strings(out)
	return out
}
func (c Catalog) ChecklistVersions(id string) []string {
	var out []string
	for _, x := range c.Checklists {
		if x.ID == id {
			out = append(out, x.Version)
		}
	}
	sort.Strings(out)
	return out
}
func (c Catalog) FilterActiveUsers(active bool) []User {
	var out []User
	for _, u := range c.Users {
		if u.Enabled == active {
			out = append(out, u)
		}
	}
	return out
}
func (c Catalog) Clone() Catalog {
	out := Catalog{Stores: append([]Store(nil), c.Stores...), Users: append([]User(nil), c.Users...), Checklists: append([]Checklist(nil), c.Checklists...)}
	return out
}
