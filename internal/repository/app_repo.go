package repository

import (
	"database/sql"
	"errors"
	"victa/internal/domain"
)

// AppRepository описывает методы для работы с приложениями.
type AppRepository interface {
	Create(app *domain.App) error
	GetAll() ([]domain.App, error)
	GetByID(id int64) (*domain.App, error)
	Update(app *domain.App) error
	Delete(id int64) error
}

type appRepo struct {
	db *sql.DB
}

// NewAppRepository возвращает реализацию AppRepository.
func NewAppRepository(db *sql.DB) AppRepository {
	return &appRepo{db: db}
}

func (r *appRepo) Create(app *domain.App) error {
	query := `
		INSERT INTO apps (company_id, name, platform, store_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, app.CompanyID, app.Name, app.Platform, app.StoreURL).
		Scan(&app.ID, &app.CreatedAt, &app.UpdatedAt)
}

func (r *appRepo) GetAll() ([]domain.App, error) {
	query := `SELECT id, company_id, name, platform, store_url, created_at, updated_at FROM apps`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []domain.App
	for rows.Next() {
		var a domain.App
		if err := rows.Scan(&a.ID, &a.CompanyID, &a.Name, &a.Platform, &a.StoreURL, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, nil
}

func (r *appRepo) GetByID(id int64) (*domain.App, error) {
	query := `SELECT id, company_id, name, platform, store_url, created_at, updated_at FROM apps WHERE id = $1`
	var a domain.App
	err := r.db.QueryRow(query, id).
		Scan(&a.ID, &a.CompanyID, &a.Name, &a.Platform, &a.StoreURL, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("app not found")
		}
		return nil, err
	}
	return &a, nil
}

func (r *appRepo) Update(app *domain.App) error {
	query := `
		UPDATE apps
		SET company_id = $1, name = $2, platform = $3, store_url = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at
	`
	return r.db.QueryRow(query, app.CompanyID, app.Name, app.Platform, app.StoreURL, app.ID).
		Scan(&app.UpdatedAt)
}

func (r *appRepo) Delete(id int64) error {
	query := `DELETE FROM apps WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
