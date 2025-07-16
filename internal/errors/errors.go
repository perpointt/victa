package errors

import "errors"

var (
	ErrCompanyNotFound     = errors.New("company not found")
	ErrRelationNotFound    = errors.New("user–company relation not found")
	ErrRoleNotFound        = errors.New("role slug not found")
	ErrAppNotFound         = errors.New("app not found")
	ErrIntegrationNotFound = errors.New("company integration not found")
	ErrUserCompanyNotFound = errors.New("user-company relation not found")
	ErrUserNotFound        = errors.New("user not found")
)
