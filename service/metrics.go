package service

import (
	"storeinspection/domain"
	"storeinspection/store"
)

type Summary struct{ Total, Open, Assigned, Submitted, Closed int }

func (s *Service) Summary(inspectionID string) (Summary, error) {
	fs, e := s.DB.ListFindings(inspectionID)
	if e != nil {
		return Summary{}, e
	}
	out := Summary{Total: len(fs)}
	for _, f := range fs {
		switch f.Status {
		case domain.StatusOpen:
			out.Open++
		case domain.StatusAssigned:
			out.Assigned++
		case domain.StatusSubmitted:
			out.Submitted++
		case domain.StatusClosed:
			out.Closed++
		}
	}
	return out, nil
}
func EnsureStore(db *store.DB, s domain.Store) error {
	if e := domain.ValidateStore(s); e != nil {
		return e
	}
	return db.SaveStore(s)
}
