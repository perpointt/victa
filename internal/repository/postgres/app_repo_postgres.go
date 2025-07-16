package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	appErr "victa/internal/errors"

	"victa/internal/domain"
)

// AppRepo реализует AppRepository через prepared‑statements.
type AppRepo struct {
	db                  *sql.DB
	stGetByID           *sql.Stmt
	stGetAllByCompanyID *sql.Stmt
	stCreate            *sql.Stmt
	stUpdate            *sql.Stmt
	stDelete            *sql.Stmt
}

// NewAppRepo подготавливает выражения; при ошибке сразу вернёт её.
func NewAppRepo(db *sql.DB) (*AppRepo, error) {
	r := &AppRepo{db: db}
	var err error

	if r.stGetByID, err = db.Prepare(`
		SELECT id, company_id, name, slug, created_at, updated_at
		  FROM apps
		 WHERE id = $1`); err != nil {
		return nil, fmt.Errorf("prepare getByID: %w", err)
	}

	if r.stGetAllByCompanyID, err = db.Prepare(`
		SELECT id, company_id, name, slug, created_at, updated_at
		  FROM apps
		 WHERE company_id = $1
		 ORDER BY created_at DESC`); err != nil {
		return nil, fmt.Errorf("prepare getAllByCompanyID: %w", err)
	}

	if r.stCreate, err = db.Prepare(`
		INSERT INTO apps (company_id, name, slug, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id, company_id, name, slug, created_at, updated_at`); err != nil {
		return nil, fmt.Errorf("prepare create: %w", err)
	}

	if r.stUpdate, err = db.Prepare(`
		UPDATE apps
		   SET name = $1,
		       slug = $2,
		       updated_at = $3
		 WHERE id = $4
		RETURNING id, company_id, name, slug, created_at, updated_at`); err != nil {
		return nil, fmt.Errorf("prepare update: %w", err)
	}

	if r.stDelete, err = db.Prepare(`DELETE FROM apps WHERE id = $1`); err != nil {
		return nil, fmt.Errorf("prepare delete: %w", err)
	}

	return r, nil
}

// Close освобождает prepared‑statements.
func (r *AppRepo) Close() error {
	for _, st := range []*sql.Stmt{
		r.stGetByID, r.stGetAllByCompanyID, r.stCreate, r.stUpdate, r.stDelete,
	} {
		if st != nil {
			if err := st.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetByID возвращает приложение или ErrAppNotFound.
func (r *AppRepo) GetByID(ctx context.Context, appID int64) (*domain.App, error) {
	var a domain.App
	err := r.stGetByID.QueryRowContext(ctx, appID).
		Scan(&a.ID, &a.CompanyID, &a.Name, &a.Slug, &a.CreatedAt, &a.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, appErr.ErrAppNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get app by id: %w", err)
	}
	return &a, nil
}

// GetAllByCompanyID возвращает приложения компании.
func (r *AppRepo) GetAllByCompanyID(ctx context.Context, companyID int64) ([]domain.App, error) {
	rows, err := r.stGetAllByCompanyID.QueryContext(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("query apps: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	list := make([]domain.App, 0, 8)
	for rows.Next() {
		var a domain.App
		if err := rows.Scan(&a.ID, &a.CompanyID, &a.Name, &a.Slug, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan app: %w", err)
		}
		list = append(list, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return list, nil
}

// Create сохраняет новое приложение.
func (r *AppRepo) Create(ctx context.Context, app *domain.App) (*domain.App, error) {
	now := time.Now().UTC()
	var a domain.App
	err := r.stCreate.QueryRowContext(ctx,
		app.CompanyID, app.Name, app.Slug, now).
		Scan(&a.ID, &a.CompanyID, &a.Name, &a.Slug, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("create app: %w", err)
	}
	return &a, nil
}

// Update изменяет имя и slug.
func (r *AppRepo) Update(ctx context.Context, app *domain.App) (*domain.App, error) {
	now := time.Now().UTC()
	var a domain.App
	err := r.stUpdate.QueryRowContext(ctx,
		app.Name, app.Slug, now, app.ID).
		Scan(&a.ID, &a.CompanyID, &a.Name, &a.Slug, &a.CreatedAt, &a.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, appErr.ErrCompanyNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update app: %w", err)
	}
	return &a, nil
}

// Delete удаляет приложение.
func (r *AppRepo) Delete(ctx context.Context, appID int64) error {
	res, err := r.stDelete.ExecContext(ctx, appID)
	if err != nil {
		return fmt.Errorf("delete app: %w", err)
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
