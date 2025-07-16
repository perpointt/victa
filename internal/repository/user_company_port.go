package repository

import (
	"context"
	"victa/internal/domain"
)

type UserCompanyRepository interface {
	GetAllByCompanyID(ctx context.Context, companyID int64) ([]domain.UserCompany, error)
	GetByCompanyAndUserID(ctx context.Context, companyID, userID int64) (*domain.UserCompany, error)
	Delete(ctx context.Context, userID, companyID int64) error
}
