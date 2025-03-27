package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"victa/internal/domain"
)

// CompanyRepository описывает методы для работы с компаниями.
type CompanyRepository interface {
	CreateCompanyWithUser(company *domain.Company, userID int64) error
	GetAll() ([]domain.Company, error)
	GetAllWithUser(userID int64) ([]domain.Company, error)
	GetByIdWithUser(userID, companyID int64) (*domain.Company, error)
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

// CreateCompanyWithUser создает компанию и связывает её с пользователем в рамках транзакции.
// Для создателя компании роль устанавливается как "admin". Перед созданием проверяется, что
// среди компаний пользователя нет компании с таким же именем.
func (r *companyRepo) CreateCompanyWithUser(company *domain.Company, userID int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Проверяем, существует ли уже компания с таким именем у пользователя.
	checkQuery := `
		SELECT EXISTS(
			SELECT 1 
			FROM companies c
			INNER JOIN user_companies uc ON c.id = uc.company_id
			WHERE uc.user_id = $1 AND c.name = $2
		)
	`
	var exists bool
	if err := tx.QueryRow(checkQuery, userID, company.Name).Scan(&exists); err != nil {
		tx.Rollback()
		return err
	}
	if exists {
		tx.Rollback()
		return fmt.Errorf("company with name %q already exists for user %d", company.Name, userID)
	}

	// Создаем новую компанию.
	createQuery := `
		INSERT INTO companies (name, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	if err := tx.QueryRow(createQuery, company.Name).
		Scan(&company.ID, &company.CreatedAt, &company.UpdatedAt); err != nil {
		tx.Rollback()
		return err
	}

	// Устанавливаем связь с ролью "admin" для создателя.
	linkQuery := `
		INSERT INTO user_companies (user_id, company_id, role)
		VALUES ($1, $2, 'admin')
	`
	if _, err := tx.Exec(linkQuery, userID, company.ID); err != nil {
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

// GetByIdWithUser возвращает компанию по её идентификатору,
// если она связана с указанным пользователем.
func (r *companyRepo) GetByIdWithUser(userID, companyID int64) (*domain.Company, error) {
	query := `
		SELECT c.id, c.name, c.created_at, c.updated_at
		FROM companies c
		INNER JOIN user_companies uc ON c.id = uc.company_id
		WHERE uc.user_id = $1 AND c.id = $2
	`
	var comp domain.Company
	err := r.db.QueryRow(query, userID, companyID).Scan(&comp.ID, &comp.Name, &comp.CreatedAt, &comp.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("company not found or access denied")
		}
		return nil, err
	}
	return &comp, nil
}

// GetAllWithUser возвращает список компаний, связанных с указанным пользователем.
func (r *companyRepo) GetAllWithUser(userID int64) ([]domain.Company, error) {
	query := `
		SELECT c.id, c.name, c.created_at, c.updated_at
		FROM companies c
		INNER JOIN user_companies uc ON c.id = uc.company_id
		WHERE uc.user_id = $1
	`
	rows, err := r.db.Query(query, userID)
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
