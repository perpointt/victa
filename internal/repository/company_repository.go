package repository

import (
	"database/sql"
	"errors"
	"victa/internal/domain"
)

// CompanyRepository описывает методы для работы с компаниями.
type CompanyRepository interface {
	Create(company domain.Company, userID int64) (*domain.Company, error)
	Update(company domain.Company) (*domain.Company, error)
	Delete(companyID int64) error
	GetAllByUserID(userID int64) ([]domain.Company, error)
	GetByID(companyID int64) (*domain.Company, error)
	GetUserRole(userID, companyID int64) (*int64, error)
	AddUserToCompany(userID, companyID int64, roleSlug string) error
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

// Update изменяет название компании и возвращает обновлённую сущность
func (r *PostgresCompanyRepo) Update(company domain.Company) (*domain.Company, error) {
	updated := &domain.Company{}
	err := r.DB.QueryRow(
		`UPDATE companies
         SET name = $1, updated_at = now()
         WHERE id = $2
         RETURNING id, name, created_at, updated_at`,
		company.Name, company.ID,
	).Scan(
		&updated.ID,
		&updated.Name,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// Delete удаляет компанию по ID; все user_companies удалятся автоматически благодаря ON DELETE CASCADE
func (r *PostgresCompanyRepo) Delete(companyID int64) error {
	res, err := r.DB.Exec(
		`DELETE FROM companies WHERE id = $1`,
		companyID,
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

// GetAllByUserID возвращает все компании, к которым привязан пользователь,
func (r *PostgresCompanyRepo) GetAllByUserID(userID int64) ([]domain.Company, error) {
	rows, err := r.DB.Query(
		`SELECT c.id, c.name, c.created_at, c.updated_at
         FROM companies c
         JOIN user_companies uc ON c.id = uc.company_id
         WHERE uc.user_id = $1
         ORDER BY c.created_at DESC`, // сортируем по дате создания
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

// GetByID возвращает компанию по её ID, или nil если не найдена
func (r *PostgresCompanyRepo) GetByID(companyID int64) (*domain.Company, error) {
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

func (r *PostgresCompanyRepo) GetUserRole(userID, companyID int64) (*int64, error) {
	var roleID *int64
	err := r.DB.QueryRow(
		`SELECT r.id
         FROM user_companies uc
         JOIN roles r ON uc.role_id = r.id
         WHERE uc.user_id = $1 AND uc.company_id = $2`,
		userID, companyID,
	).Scan(&roleID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return roleID, err
}

func (r *PostgresCompanyRepo) AddUserToCompany(userID, companyID int64, roleSlug string) error {
	_, err := r.DB.Exec(
		`INSERT INTO user_companies (user_id, company_id, role_id)
         SELECT $1, $2, id
           FROM roles
          WHERE slug = $3`,
		userID, companyID, roleSlug,
	)
	return err
}
