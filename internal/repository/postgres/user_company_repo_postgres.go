package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"victa/internal/domain"
	appErr "victa/internal/errors"
)

type UserCompanyRepo struct {
	db                      *sql.DB
	stGetAllByCompanyID     *sql.Stmt
	stGetByCompanyAndUserID *sql.Stmt
	stDelete                *sql.Stmt
}

// NewUserCompanyRepo инициализирует репозиторий.
func NewUserCompanyRepo(db *sql.DB) (*UserCompanyRepo, error) {
	r := &UserCompanyRepo{db: db}

	var err error
	if r.stGetAllByCompanyID, err = db.Prepare(`
		SELECT user_id, company_id, role_id
		FROM user_companies
		WHERE company_id = $1
		ORDER BY role_id, user_id`); err != nil {
		return nil, fmt.Errorf("prepare GetAllSecretsByCompanyID: %w", err)
	}

	if r.stGetByCompanyAndUserID, err = db.Prepare(`
		SELECT user_id, company_id, role_id
		FROM user_companies
		WHERE company_id = $1 AND user_id = $2`); err != nil {
		return nil, fmt.Errorf("prepare GetByCompanyAndUserID: %w", err)
	}

	if r.stDelete, err = db.Prepare(`
		DELETE FROM user_companies
		WHERE user_id = $1 AND company_id = $2`); err != nil {
		return nil, fmt.Errorf("prepare Delete: %w", err)
	}

	return r, nil
}

// Close освобождает prepared-statements
func (r *UserCompanyRepo) Close() error {
	if r == nil {
		return nil
	}
	if err := r.stGetAllByCompanyID.Close(); err != nil {
		return err
	}
	if err := r.stGetByCompanyAndUserID.Close(); err != nil {
		return err
	}
	return r.stDelete.Close()
}

// GetAllByCompanyID возвращает все связи пользователей с компанией
func (r *UserCompanyRepo) GetAllByCompanyID(ctx context.Context, companyID int64) ([]domain.UserCompany, error) {
	rows, err := r.stGetAllByCompanyID.QueryContext(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("query GetAllSecretsByCompanyID: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	list := make([]domain.UserCompany, 0, 16) // типовой отдел
	for rows.Next() {
		var uc domain.UserCompany
		if err := rows.Scan(&uc.UserID, &uc.CompanyID, &uc.RoleID); err != nil {
			return nil, fmt.Errorf("scan user_company: %w", err)
		}
		list = append(list, uc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return list, nil
}

// GetByCompanyAndUserID возвращает одну связь или ErrUserCompanyNotFound
func (r *UserCompanyRepo) GetByCompanyAndUserID(ctx context.Context, companyID, userID int64) (*domain.UserCompany, error) {
	var uc domain.UserCompany
	err := r.stGetByCompanyAndUserID.QueryRowContext(ctx, companyID, userID).
		Scan(&uc.UserID, &uc.CompanyID, &uc.RoleID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, appErr.ErrUserCompanyNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get relation: %w", err)
	}
	return &uc, nil
}

// Delete удаляет связь или возвращает ErrUserCompanyNotFound
func (r *UserCompanyRepo) Delete(ctx context.Context, userID, companyID int64) error {
	res, err := r.stDelete.ExecContext(ctx, userID, companyID)
	if err != nil {
		return fmt.Errorf("delete relation: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rowsAffected: %w", err)
	}
	if affected == 0 {
		return appErr.ErrUserCompanyNotFound
	}
	return nil
}
