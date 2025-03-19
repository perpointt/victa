package repository

import (
	"database/sql"
)

// UserCompanyRepository описывает методы для работы со связями между пользователями и компаниями.
type UserCompanyRepository interface {
	// LinkUserCompany связывает пользователя с компанией.
	LinkUserCompany(userID, companyID int64) error
	// (Опционально) можно добавить методы для получения или удаления связей.
}

type userCompanyRepo struct {
	db *sql.DB
}

// NewUserCompanyRepository создаёт новый репозиторий для работы со связями пользователя и компании.
func NewUserCompanyRepository(db *sql.DB) UserCompanyRepository {
	return &userCompanyRepo{db: db}
}

// LinkUserCompany добавляет связь между пользователем и компанией.
// Если такая связь уже существует, конфликт игнорируется.
func (r *userCompanyRepo) LinkUserCompany(userID, companyID int64) error {
	query := `
		INSERT INTO user_companies (user_id, company_id, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		ON CONFLICT (user_id, company_id) DO NOTHING
	`
	_, err := r.db.Exec(query, userID, companyID)
	return err
}
