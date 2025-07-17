package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"victa/internal/domain"

	appErr "victa/internal/errors"
)

// CompanyRepo реализует CompanyRepository через prepared-statements.
type CompanyRepo struct {
	db *sql.DB

	stCreateCompany  *sql.Stmt
	stLinkAdmin      *sql.Stmt
	stAddUser        *sql.Stmt
	stUpdate         *sql.Stmt
	stDelete         *sql.Stmt
	stGetAllByUserID *sql.Stmt
	stGetByID        *sql.Stmt
	stGetUserRole    *sql.Stmt
}

// NewCompanyRepo инициализирует репозиторий.
func NewCompanyRepo(db *sql.DB) (*CompanyRepo, error) {
	r := &CompanyRepo{db: db}
	var err error

	if r.stCreateCompany, err = db.Prepare(`
		INSERT INTO companies (name, created_at, updated_at)
		VALUES ($1, $2, $2)
		RETURNING id, name, created_at, updated_at`); err != nil {
		return nil, fmt.Errorf("prepare create company: %w", err)
	}
	if r.stLinkAdmin, err = db.Prepare(`
		INSERT INTO user_companies (user_id, company_id, role_id)
		VALUES ($1, $2, (SELECT id FROM roles WHERE slug = 'admin'))`); err != nil {
		return nil, fmt.Errorf("prepare link admin: %w", err)
	}
	if r.stAddUser, err = db.Prepare(`
		INSERT INTO user_companies (user_id, company_id, role_id)
		SELECT $1, $2, id FROM roles WHERE slug = $3`); err != nil {
		return nil, fmt.Errorf("prepare add user: %w", err)
	}
	if r.stUpdate, err = db.Prepare(`
		UPDATE companies
		SET name = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, name, created_at, updated_at`); err != nil {
		return nil, fmt.Errorf("prepare update company: %w", err)
	}
	if r.stDelete, err = db.Prepare(`DELETE FROM companies WHERE id = $1`); err != nil {
		return nil, fmt.Errorf("prepare delete company: %w", err)
	}
	if r.stGetAllByUserID, err = db.Prepare(`
		SELECT c.id, c.name, c.created_at, c.updated_at
		FROM companies c
		JOIN user_companies uc ON c.id = uc.company_id
		WHERE uc.user_id = $1
		ORDER BY c.created_at DESC`); err != nil {
		return nil, fmt.Errorf("prepare get all by user: %w", err)
	}
	if r.stGetByID, err = db.Prepare(`
		SELECT id, name, created_at, updated_at
		FROM companies
		WHERE id = $1`); err != nil {
		return nil, fmt.Errorf("prepare get by id: %w", err)
	}
	if r.stGetUserRole, err = db.Prepare(`
		SELECT r.id
		FROM user_companies uc
		JOIN roles r ON uc.role_id = r.id
		WHERE uc.user_id = $1 AND uc.company_id = $2`); err != nil {
		return nil, fmt.Errorf("prepare get user role: %w", err)
	}

	return r, nil
}

// Close освобождает ресурсы prepared-statement.
func (r *CompanyRepo) Close() error {
	if r == nil {
		return nil
	}
	for _, st := range []*sql.Stmt{
		r.stCreateCompany, r.stLinkAdmin, r.stAddUser,
		r.stUpdate, r.stDelete,
		r.stGetAllByUserID, r.stGetByID, r.stGetUserRole,
	} {
		if st != nil {
			if err := st.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Create создаёт компанию + делает userID админом в одной транзакции.
func (r *CompanyRepo) Create(ctx context.Context, company domain.Company, userID int64) (*domain.Company, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	now := time.Now().UTC()
	created := new(domain.Company)

	if err = tx.StmtContext(ctx, r.stCreateCompany).
		QueryRowContext(ctx, company.Name, now).
		Scan(&created.ID, &created.Name, &created.CreatedAt, &created.UpdatedAt); err != nil {
		return nil, fmt.Errorf("insert company: %w", err)
	}

	if _, err = tx.StmtContext(ctx, r.stLinkAdmin).
		ExecContext(ctx, userID, created.ID); err != nil {
		return nil, fmt.Errorf("link admin: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	return created, nil
}

// Update меняет название компании.
func (r *CompanyRepo) Update(ctx context.Context, company domain.Company) (*domain.Company, error) {
	updated := new(domain.Company)
	err := r.stUpdate.QueryRowContext(ctx, company.Name, time.Now().UTC(), company.ID).
		Scan(&updated.ID, &updated.Name, &updated.CreatedAt, &updated.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, appErr.ErrCompanyNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update company: %w", err)
	}
	return updated, nil
}

// Delete удаляет компанию; каскад гарантирован на уровне FK.
func (r *CompanyRepo) Delete(ctx context.Context, companyID int64) error {
	res, err := r.stDelete.ExecContext(ctx, companyID)
	if err != nil {
		return fmt.Errorf("delete company: %w", err)
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rowsAffected: %w", err)
	}
	if aff == 0 {
		return appErr.ErrCompanyNotFound
	}
	return nil
}

// GetAllByUserID возвращает компании пользователя.
func (r *CompanyRepo) GetAllByUserID(ctx context.Context, userID int64) ([]domain.Company, error) {
	rows, err := r.stGetAllByUserID.QueryContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query companies: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	list := make([]domain.Company, 0, 8)
	for rows.Next() {
		var c domain.Company
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan company: %w", err)
		}
		list = append(list, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return list, nil
}

// GetByID возвращает компанию или ErrCompanyNotFound.
func (r *CompanyRepo) GetByID(ctx context.Context, companyID int64) (*domain.Company, error) {
	var c domain.Company
	err := r.stGetByID.QueryRowContext(ctx, companyID).
		Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, appErr.ErrCompanyNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get company by id: %w", err)
	}
	return &c, nil
}

// GetUserRole возвращает роль пользователя в компании.
func (r *CompanyRepo) GetUserRole(ctx context.Context, userID, companyID int64) (*int64, error) {
	var roleID int64
	err := r.stGetUserRole.QueryRowContext(ctx, userID, companyID).Scan(&roleID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, appErr.ErrRelationNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user role: %w", err)
	}
	return &roleID, nil
}

// AddUserToCompany даёт пользователю роль по слагу.
func (r *CompanyRepo) AddUserToCompany(ctx context.Context, userID, companyID int64, roleSlug string) error {
	res, err := r.stAddUser.ExecContext(ctx, userID, companyID, roleSlug)
	if err != nil {
		return fmt.Errorf("add user: %w", err)
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rowsAffected: %w", err)
	}
	if aff == 0 {
		return appErr.ErrRoleNotFound
	}
	return nil
}
