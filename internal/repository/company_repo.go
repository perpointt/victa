package repository

import (
	"database/sql"
	"errors"
	"victa/internal/domain"
)

// CompanyRepository описывает методы для работы с компаниями.
type CompanyRepository interface {
	Create(company *domain.Company) error
	CreateAndLink(company *domain.Company, userID int64) error
	GetAll() ([]domain.Company, error)
	GetByID(id int64) (*domain.Company, error)
	Update(company *domain.Company) error
	Delete(id int64) error
}

type companyRepo struct {
	db *sql.DB
}

// NewCompanyRepository возвращает реализацию CompanyRepository.
func NewCompanyRepository(db *sql.DB) CompanyRepository {
	return &companyRepo{db: db}
}

func (r *companyRepo) Create(company *domain.Company) error {
	query := `
		INSERT INTO companies (name, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, company.Name).
		Scan(&company.ID, &company.CreatedAt, &company.UpdatedAt)
}

// CreateAndLink создает компанию и связывает её с пользователем в рамках транзакции.
func (r *companyRepo) CreateAndLink(company *domain.Company, userID int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	createQuery := `
		INSERT INTO companies (name, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(createQuery, company.Name).
		Scan(&company.ID, &company.CreatedAt, &company.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	linkQuery := `
		INSERT INTO user_companies (user_id, company_id, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
	`
	_, err = tx.Exec(linkQuery, userID, company.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *companyRepo) GetAll() ([]domain.Company, error) {
	query := `SELECT id, name, created_at, updated_at FROM companies`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []domain.Company
	for rows.Next() {
		var comp domain.Company
		if err := rows.Scan(&comp.ID, &comp.Name, &comp.CreatedAt, &comp.UpdatedAt); err != nil {
			return nil, err
		}
		companies = append(companies, comp)
	}
	return companies, nil
}

func (r *companyRepo) GetByID(id int64) (*domain.Company, error) {
	query := `SELECT id, name, created_at, updated_at FROM companies WHERE id = $1`
	var comp domain.Company
	err := r.db.QueryRow(query, id).Scan(&comp.ID, &comp.Name, &comp.CreatedAt, &comp.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("company not found")
		}
		return nil, err
	}
	return &comp, nil
}

func (r *companyRepo) Update(company *domain.Company) error {
	query := `
		UPDATE companies 
		SET name = $1, updated_at = NOW() 
		WHERE id = $2
		RETURNING updated_at
	`
	return r.db.QueryRow(query, company.Name, company.ID).Scan(&company.UpdatedAt)
}

func (r *companyRepo) Delete(id int64) error {
	query := `DELETE FROM companies WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
