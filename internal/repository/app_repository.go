package repository

import (
	"database/sql"
	"errors"

	"victa/internal/domain"
)

// AppRepository описывает CRUD для приложений.
type AppRepository interface {
	// GetByID возвращает приложение по его ID или nil, если не найдено.
	GetByID(id int64) (*domain.App, error)
	// GetAllByCompanyID возвращает все приложения компании, отсортированные по дате создания (DESC).
	GetAllByCompanyID(companyID int64) ([]domain.App, error)
	// Create сохраняет новое приложение и возвращает созданную сущность.
	Create(app *domain.App) (*domain.App, error)
	// Update изменяет имя и slug приложения и возвращает обновлённую сущность.
	Update(app *domain.App) (*domain.App, error)
	// Delete удаляет приложение по ID.
	Delete(id int64) error
}

// PostgresAppRepo реализует AppRepository через Postgres.
type PostgresAppRepo struct {
	DB *sql.DB
}

// NewPostgresAppRepo создаёт репозиторий приложений.
func NewPostgresAppRepo(db *sql.DB) *PostgresAppRepo {
	return &PostgresAppRepo{DB: db}
}

// GetByID возвращает приложение по его ID или nil, если не найдено.
func (r *PostgresAppRepo) GetByID(id int64) (*domain.App, error) {
	var a domain.App
	err := r.DB.QueryRow(
		`SELECT id, company_id, name, slug, created_at, updated_at
         FROM apps
         WHERE id = $1`, id,
	).Scan(
		&a.ID,
		&a.CompanyID,
		&a.Name,
		&a.Slug,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// GetAllByCompanyID возвращает все приложения компании, отсортированные по дате создания (DESC).
func (r *PostgresAppRepo) GetAllByCompanyID(companyID int64) ([]domain.App, error) {
	rows, err := r.DB.Query(
		`SELECT id, company_id, name, slug, created_at, updated_at
         FROM apps
         WHERE company_id = $1
         ORDER BY created_at DESC`, companyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.App
	for rows.Next() {
		var a domain.App
		if err := rows.Scan(
			&a.ID,
			&a.CompanyID,
			&a.Name,
			&a.Slug,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

// Create сохраняет новое приложение и возвращает созданную сущность.
func (r *PostgresAppRepo) Create(app *domain.App) (*domain.App, error) {
	var a domain.App
	err := r.DB.QueryRow(
		`INSERT INTO apps (company_id, name, slug, created_at, updated_at)
         VALUES ($1, $2, $3, now(), now())
         RETURNING id, company_id, name, slug, created_at, updated_at`,
		app.CompanyID, app.Name, app.Slug,
	).Scan(
		&a.ID,
		&a.CompanyID,
		&a.Name,
		&a.Slug,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Update изменяет имя и slug приложения и возвращает обновлённую сущность.
func (r *PostgresAppRepo) Update(app *domain.App) (*domain.App, error) {
	var a domain.App
	err := r.DB.QueryRow(
		`UPDATE apps
         SET name = $1, slug = $2, updated_at = now()
         WHERE id = $3
         RETURNING id, company_id, name, slug, created_at, updated_at`,
		app.Name, app.Slug, app.ID,
	).Scan(
		&a.ID,
		&a.CompanyID,
		&a.Name,
		&a.Slug,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Delete удаляет приложение по ID.
func (r *PostgresAppRepo) Delete(id int64) error {
	res, err := r.DB.Exec(
		`DELETE FROM apps WHERE id = $1`, id,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
