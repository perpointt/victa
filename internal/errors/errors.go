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
	ErrInvalidInput        = errors.New("app name and slug must be non-empty")
	ErrSecretNotFound      = errors.New("secret not found")
	ErrDecrypt             = errors.New("invalid key or corrupted data")
	ErrNotCompanyAdmin     = errors.New("operation allowed for company admin only")
	ErrBadFormat           = errors.New("token format is invalid")
	ErrExpired             = errors.New("token expired")
	ErrBadSignature        = errors.New("token signature mismatch")
	ErrUnexpectedMethod    = errors.New("unexpected signing method")
	ErrInvalidToken        = errors.New("invalid token")
	ErrClaimMissing        = errors.New("company_id claim missing")
)
