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

// CompanySecretRepo реализует хранение company_secrets.
type CompanySecretRepo struct {
	db        *sql.DB
	stSave    *sql.Stmt
	stLoad    *sql.Stmt
	stLoadAll *sql.Stmt
}

// NewCompanySecretRepo подготавливает выражения.
func NewCompanySecretRepo(db *sql.DB) (*CompanySecretRepo, error) {
	r := &CompanySecretRepo{db: db}
	var err error

	if r.stSave, err = db.Prepare(`
		INSERT INTO company_secrets (company_id, secret_type, cipher, created_at)
		     VALUES ($1, $2, $3, now())
		ON CONFLICT (company_id, secret_type)
		DO UPDATE SET cipher = EXCLUDED.cipher,
		              created_at = now()`); err != nil {
		return nil, fmt.Errorf("prepare save secret: %w", err)
	}

	if r.stLoad, err = db.Prepare(`
		SELECT cipher, created_at
		  FROM company_secrets
		 WHERE company_id = $1
		   AND secret_type = $2`); err != nil {
		return nil, fmt.Errorf("prepare load secret: %w", err)
	}

	if r.stLoadAll, err = db.Prepare(`
		SELECT secret_type, cipher, created_at
		  FROM company_secrets
		 WHERE company_id = $1`); err != nil {
		return nil, fmt.Errorf("prepare load‑all secrets: %w", err)
	}
	return r, nil
}

// Close освобождает prepared statements.
func (r *CompanySecretRepo) Close() error {
	for _, st := range []*sql.Stmt{r.stSave, r.stLoad, r.stLoadAll} {
		if st != nil {
			if err := st.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Create кладёт / обновляет секрет.
func (r *CompanySecretRepo) Create(ctx context.Context, sec *domain.CompanySecret) error {
	_, err := r.stSave.ExecContext(
		ctx,
		sec.CompanyID,
		sec.Type,
		sec.Cipher,
	)
	return err
}

// GetByCompanyIDAndType возвращает один секрет.
func (r *CompanySecretRepo) GetByCompanyIDAndType(
	ctx context.Context,
	companyID int64,
	secretType domain.SecretType,
) (*domain.CompanySecret, error) {

	var (
		cipher    []byte
		createdAt time.Time
	)

	// используем подготовленный statement без SQL‑текста
	err := r.stLoad.QueryRowContext(ctx, companyID, secretType).
		Scan(&cipher, &createdAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, appErr.ErrSecretNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("load secret: %w", err)
	}

	return &domain.CompanySecret{
		CompanyID: companyID,
		Type:      secretType,
		Cipher:    cipher,
		CreatedAt: createdAt,
	}, nil
}

// GetAllByCompanyID возвращает все секреты компании.
func (r *CompanySecretRepo) GetAllByCompanyID(
	ctx context.Context,
	companyID int64,
) ([]domain.CompanySecret, error) {
	rows, err := r.stLoadAll.QueryContext(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("query secrets: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	list := make([]domain.CompanySecret, 0, 4)
	for rows.Next() {
		var s domain.CompanySecret
		if err := rows.Scan(&s.Type, &s.Cipher); err != nil {
			return nil, fmt.Errorf("scan secret: %w", err)
		}
		list = append(list, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	if len(list) == 0 {
		return nil, appErr.ErrSecretNotFound
	}
	return list, nil
}
