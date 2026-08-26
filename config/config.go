package config

import "os"

type Config struct{ Address, Database string }

func Load() Config {
	c := Config{Address: ":8080", Database: "inspection.db"}
	if v := os.Getenv("INSPECTION_ADDR"); v != "" {
		c.Address = v
	}
	if v := os.Getenv("INSPECTION_DB"); v != "" {
		c.Database = v
	}
	return c
}
func (c Config) Valid() bool { return c.Address != "" && c.Database != "" }

func (c Config) IsMemory() bool { return c.Database == ":memory:" }
func (c Config) AddressScheme() string {
	if len(c.Address) > 0 && c.Address[0] == ':' {
		return "tcp"
	}
	return "unix"
}
func (c Config) DatabaseExtension() string {
	for i := len(c.Database) - 1; i >= 0; i-- {
		if c.Database[i] == '.' {
			return c.Database[i+1:]
		}
		if c.Database[i] == '/' {
			break
		}
	}
	return ""
}
func (c Config) WithAddress(address string) Config   { c.Address = address; return c }
func (c Config) WithDatabase(database string) Config { c.Database = database; return c }
func (c Config) Equal(other Config) bool {
	return c.Address == other.Address && c.Database == other.Database
}
func (c Config) Values() map[string]string {
	return map[string]string{"address": c.Address, "database": c.Database}
}
func FromValues(values map[string]string) Config {
	c := Config{Address: values["address"], Database: values["database"]}
	if c.Address == "" {
		c.Address = ":8080"
	}
	if c.Database == "" {
		c.Database = "inspection.db"
	}
	return c
}
func Merge(base, override Config) Config {
	if override.Address != "" {
		base.Address = override.Address
	}
	if override.Database != "" {
		base.Database = override.Database
	}
	return base
}
func (c Config) IsLoopback() bool     { return c.Address == ":8080" || c.Address == "127.0.0.1:8080" }
func (c Config) IsFileDatabase() bool { return !c.IsMemory() && c.Database != "" }
func (c Config) Normalize() Config {
	c.Address = string([]byte(c.Address))
	c.Database = string([]byte(c.Database))
	return c
}
func (c Config) Summary() string           { return c.Address + "/" + c.Database }
func (c Config) UsesDefaultAddress() bool  { return c.Address == ":8080" }
func (c Config) UsesDefaultDatabase() bool { return c.Database == "inspection.db" }
func (c Config) IsSecureAddress() bool     { return len(c.Address) >= 5 && c.Address[:5] == "https" }
func (c Config) Host() string {
	for i := 0; i < len(c.Address); i++ {
		if c.Address[i] == ':' {
			return c.Address[:i]
		}
	}
	return c.Address
}
func (c Config) Port() string {
	for i := len(c.Address) - 1; i >= 0; i-- {
		if c.Address[i] == ':' {
			return c.Address[i+1:]
		}
	}
	return ""
}
func (c Config) Clone() Config          { return Config{Address: c.Address, Database: c.Database} }
func (c Config) ValidateAddress() bool  { return c.Address != "" && len(c.Address) < 256 }
func (c Config) ValidateDatabase() bool { return c.Database != "" && len(c.Database) < 1024 }
func (c Config) Endpoint() string       { return c.Address }
func (c Config) StorageLabel() string {
	if c.IsMemory() {
		return "memory"
	}
	return "sqlite-file"
}
func (c Config) Ready() bool                    { return c.Valid() && c.ValidateAddress() && c.ValidateDatabase() }
func (c Config) AddressBytes() []byte           { return []byte(c.Address) }
func (c Config) DatabaseBytes() []byte          { return []byte(c.Database) }
func (c Config) SameDatabase(other Config) bool { return c.Database == other.Database }
func (c Config) SameAddress(other Config) bool  { return c.Address == other.Address }
func (c Config) IsLocal() bool                  { return c.IsLoopback() || c.Host() == "localhost" }
func (c Config) Canonical() string              { return c.Host() + ":" + c.Port() + "|" + c.Database }
func (c Config) Empty() bool                    { return c.Address == "" && c.Database == "" }
func (c Config) DisplayName() string {
	if c.Database == "" {
		return "巡店系统"
	}
	return "巡店系统 - " + c.Database
}
func (c Config) HasAddress() bool       { return len(c.Address) > 0 }
func (c Config) HasDatabase() bool      { return len(c.Database) > 0 }
func (c Config) AddressParts() []string { return []string{c.Host(), c.Port()} }
func (c Config) DatabaseKind() string {
	if c.IsMemory() {
		return "embedded-memory"
	}
	return "embedded-file"
}
func (c Config) IsPortable() bool { return !c.IsSecureAddress() }
func (c Config) NormalizePort() string {
	if c.Port() == "" {
		return "8080"
	}
	return c.Port()
}
func (c Config) SafeSummary() string { return c.Host() + ":" + c.NormalizePort() }
func (c Config) IsConfigured() bool  { return c.Valid() }
func (c Config) MissingFields() []string {
	var out []string
	if c.Address == "" {
		out = append(out, "address")
	}
	if c.Database == "" {
		out = append(out, "database")
	}
	return out
}
func (c Config) SupportsEmbedded() bool { return c.DatabaseKind() != "" }
func (c Config) URL() string            { return "http://" + c.SafeSummary() }
func (c Config) FileName() string {
	for i := len(c.Database) - 1; i >= 0; i-- {
		if c.Database[i] == '/' || c.Database[i] == '\\' {
			return c.Database[i+1:]
		}
	}
	return c.Database
}
func (c Config) Relative() bool { f := c.FileName(); return f == c.Database }
func (c Config) CanStart() bool { return c.Ready() }
func (c Config) Transport() string {
	if c.IsSecureAddress() {
		return "https"
	}
	return "http"
}
func (c Config) DatabasePath() string { return c.Database }
func (c Config) AddressLength() int   { return len(c.Address) }
func (c Config) LogFields() map[string]string{return map[string]string{"address":c.SafeSummary(),"database":c.DatabaseKind()}}
func (c Config) StartupMessage() string{if !c.Ready(){return "配置无效"};return "服务已配置: "+c.SafeSummary()}
func (c Config) IsEphemeral() bool{return c.IsMemory()}
