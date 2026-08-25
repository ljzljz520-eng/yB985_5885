package config

type Policy struct {
	RequireAttachment, RequireApproval bool
	MaxOpenFindings                    int
}

func DefaultPolicy() Policy {
	return Policy{RequireAttachment: true, RequireApproval: true, MaxOpenFindings: 50}
}
func (p Policy) AllowsOpen(count int) bool {
	return p.MaxOpenFindings <= 0 || count < p.MaxOpenFindings
}
func (p Policy) Validate() bool { return p.MaxOpenFindings >= 0 }
