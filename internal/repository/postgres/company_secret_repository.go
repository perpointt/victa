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

type CompanySecretRepo struct {
	db        *sql.DB
	stCreate  *sql.Stmt
	stLoad    *sql.Stmt
	stLoadAll *sql.Stmt
}

func NewCompanySecretRepo(db *sql.DB) (*CompanySecretRepo, error) {
	r := &CompanySecretRepo{db: db}
	var err error

	// INSERT + UPSERT: created_at фиксируем при первой вставке, updated_at меняем всегда
	if r.stCreate, err = db.Prepare(`
		INSERT INTO company_secrets (company_id, secret_type, cipher, created_at, updated_at)
		     VALUES ($1, $2, $3, now(), now())
		ON CONFLICT (company_id, secret_type)
		DO UPDATE SET cipher     = EXCLUDED.cipher,
		              updated_at = now()
		RETURNING company_id, secret_type, cipher, created_at, updated_at
	`); err != nil {
		return nil, fmt.Errorf("prepare create secret: %w", err)
	}

	if r.stLoad, err = db.Prepare(`
		SELECT cipher, created_at, updated_at
		  FROM company_secrets
		 WHERE company_id = $1
		   AND secret_type = $2`); err != nil {
		return nil, fmt.Errorf("prepare load secret: %w", err)
	}

	if r.stLoadAll, err = db.Prepare(`
		SELECT secret_type, cipher, created_at, updated_at
		  FROM company_secrets
		 WHERE company_id = $1`); err != nil {
		return nil, fmt.Errorf("prepare load-all secrets: %w", err)
	}

	return r, nil
}

func (r *CompanySecretRepo) Close() error {
	for _, st := range []*sql.Stmt{r.stCreate, r.stLoad, r.stLoadAll} {
		if st != nil {
			if err := st.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Create вставляет или обновляет секрет и возвращает фактические значения из БД.
func (r *CompanySecretRepo) Create(
	ctx context.Context,
	sec *domain.CompanySecret,
) (*domain.CompanySecret, error) {

	var out domain.CompanySecret
	err := r.stCreate.QueryRowContext(
		ctx,
		sec.CompanyID,
		sec.Type,
		sec.Cipher,
	).Scan(
		&out.CompanyID,
		&out.Type,
		&out.Cipher,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create secret: %w", err)
	}
	return &out, nil
}

// GetByCompanyIDAndType один секрет
func (r *CompanySecretRepo) GetByCompanyIDAndType(
	ctx context.Context,
	companyID int64,
	secretType domain.SecretType,
) (*domain.CompanySecret, error) {

	var (
		cipher               []byte
		createdAt, updatedAt time.Time
	)

	err := r.stLoad.QueryRowContext(ctx, companyID, secretType).
		Scan(&cipher, &createdAt, &updatedAt)

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
		UpdatedAt: updatedAt,
	}, nil
}

// GetAllByCompanyID все секреты компании
func (r *CompanySecretRepo) GetAllByCompanyID(
	ctx context.Context,
	companyID int64,
) ([]domain.CompanySecret, error) {

	rows, err := r.stLoadAll.QueryContext(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("query secrets: %w", err)
	}
	defer rows.Close()

	list := make([]domain.CompanySecret, 0, 4)
	for rows.Next() {
		var s domain.CompanySecret
		if err := rows.Scan(&s.Type, &s.Cipher, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan secret: %w", err)
		}
		s.CompanyID = companyID
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
