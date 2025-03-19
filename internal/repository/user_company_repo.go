package repository

import (
	"database/sql"
)

// UserCompanyRepository описывает методы для работы со связями между пользователями и компаниями.
type UserCompanyRepository interface {
	LinkUserCompanyWithRole(userID, companyID int64, role string) error
	GetUserRole(userID, companyID int64) (string, error)
	RemoveUserFromCompany(userID, companyID int64) error
}

type userCompanyRepo struct {
	db *sql.DB
}

// NewUserCompanyRepository создаёт новый репозиторий для работы со связями пользователя и компании.
func NewUserCompanyRepository(db *sql.DB) UserCompanyRepository {
	return &userCompanyRepo{db: db}
}

// LinkUserCompanyWithRole добавляет связь между пользователем и компанией.
// Если такая связь уже существует, конфликт игнорируется.
func (r *userCompanyRepo) LinkUserCompanyWithRole(userID, companyID int64, role string) error {
	query := `
		INSERT INTO user_companies (user_id, company_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, company_id) DO UPDATE SET role = EXCLUDED.role
	`
	_, err := r.db.Exec(query, userID, companyID, role)
	return err
}

// GetUserRole возвращает роль пользователя в указанной компании.
func (r *userCompanyRepo) GetUserRole(userID, companyID int64) (string, error) {
	query := `SELECT role FROM user_companies WHERE user_id = $1 AND company_id = $2`
	var role string
	err := r.db.QueryRow(query, userID, companyID).Scan(&role)
	if err != nil {
		return "", err
	}
	return role, nil
}

// RemoveUserFromCompany удаляет связь между пользователем и компанией.
func (r *userCompanyRepo) RemoveUserFromCompany(userID, companyID int64) error {
	query := `DELETE FROM user_companies WHERE user_id = $1 AND company_id = $2`
	_, err := r.db.Exec(query, userID, companyID)
	return err
}
