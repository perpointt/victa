package repository

import (
	"database/sql"
	"errors"
	"victa/internal/domain"
)

// CompanyRepository описывает методы для работы с компаниями.
type CompanyRepository interface {
	Create(company domain.Company, userID int64) (*domain.Company, error)
	GetAllByUserId(userID int64) ([]domain.Company, error)
	GetById(companyID int64) (*domain.Company, error)
}

// PostgresCompanyRepo реализует CompanyRepository через Postgres
type PostgresCompanyRepo struct {
	DB *sql.DB
}

// NewPostgresCompanyRepo возвращает реализацию CompanyRepository.
func NewPostgresCompanyRepo(db *sql.DB) *PostgresCompanyRepo {
	return &PostgresCompanyRepo{DB: db}
}

// Create вставляет новую компанию и связывает её с пользователем-администратором в одной транзакции
// и возвращает только что созданную сущность компании.
func (r *PostgresCompanyRepo) Create(company domain.Company, userID int64) (*domain.Company, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1) Вставляем компанию и возвращаем её ID и метки времени
	created := &domain.Company{}
	err = tx.QueryRow(
		`INSERT INTO companies (name, created_at, updated_at)
         VALUES ($1, now(), now())
         RETURNING id, name, created_at, updated_at`,
		company.Name,
	).Scan(&created.ID, &created.Name, &created.CreatedAt, &created.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// 2) Связываем пользователя с созданной компанией как admin
	_, err = tx.Exec(
		`INSERT INTO user_companies (user_id, company_id, role_id)
         VALUES ($1, $2, (SELECT id FROM roles WHERE slug = 'admin'))`,
		userID, created.ID,
	)
	if err != nil {
		return nil, err
	}

	// 3) Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return created, nil
}

// GetAllByUserId возвращает все компании, к которым привязан пользователь
func (r *PostgresCompanyRepo) GetAllByUserId(userID int64) ([]domain.Company, error) {
	rows, err := r.DB.Query(
		`SELECT c.id, c.name, c.created_at, c.updated_at
         FROM companies c
         JOIN user_companies uc ON c.id = uc.company_id
         WHERE uc.user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Company
	for rows.Next() {
		var c domain.Company
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

// GetById возвращает компанию по её ID, или nil если не найдена
func (r *PostgresCompanyRepo) GetById(companyID int64) (*domain.Company, error) {
	var c domain.Company
	err := r.DB.QueryRow(
		`SELECT id, name, created_at, updated_at
         FROM companies
         WHERE id = $1`,
		companyID,
	).Scan(
		&c.ID,
		&c.Name,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
