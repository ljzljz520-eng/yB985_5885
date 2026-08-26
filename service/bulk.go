package service

import (
	"errors"
	"storeinspection/domain"
)

type BulkResult struct {
	Processed, Failed int
	Errors            []string
}

func (s *Service) BulkAssign(ids []string, owner string) BulkResult {
	out := BulkResult{}
	for _, id := range ids {
		out.Processed++
		if e := s.Assign(id, owner, "bulk assignment"); e != nil {
			out.Failed++
			out.Errors = append(out.Errors, e.Error())
		}
	}
	return out
}
func (s *Service) BulkSubmit(ids []string, author, text string) BulkResult {
	out := BulkResult{}
	for _, id := range ids {
		out.Processed++
		if e := s.Submit(id, author, text, ""); e != nil {
			out.Failed++
			out.Errors = append(out.Errors, e.Error())
		}
	}
	return out
}
func (s *Service) ValidateBatch(fs []domain.Finding) error {
	if len(fs) == 0 {
		return errors.New("empty batch")
	}
	seen := map[string]bool{}
	for _, f := range fs {
		if seen[f.ID] {
			return errors.New("duplicate finding")
		}
		seen[f.ID] = true
		if e := domain.ValidateFinding(f); e != nil {
			return e
		}
	}
	return nil
}
func (s *Service) CloseMany(ids []string) BulkResult {
	out := BulkResult{}
	for _, id := range ids {
		out.Processed++
		if e := s.Close(id); e != nil {
			out.Failed++
			out.Errors = append(out.Errors, e.Error())
		}
	}
	return out
}
