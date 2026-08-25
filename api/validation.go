package api

import (
	"errors"
	"net/http"
	"storeinspection/domain"
)

func ValidateRequestMethod(r *http.Request, expected string) error {
	if r.Method != expected {
		return errors.New("method not allowed")
	}
	return nil
}
func ValidateFilter(f domain.InspectionFilter) error {
	if f.Status != "" && !domain.ValidInspectionStatus(f.Status) {
		return errors.New("invalid status filter")
	}
	return nil
}
func ErrorCode(err error) string {
	if err == nil {
		return ""
	}
	switch err.Error() {
	case "storage unavailable":
		return "storage_unavailable"
	case "owner required":
		return "owner_required"
	default:
		return "business_error"
	}
}
func StatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if ErrorCode(err) == "storage_unavailable" {
		return http.StatusServiceUnavailable
	}
	return http.StatusUnprocessableEntity
}
