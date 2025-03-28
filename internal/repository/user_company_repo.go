package repository

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

// UserCompanyRepository описывает методы для работы со связями между пользователями и компаниями.
type UserCompanyRepository interface {
	LinkUserWithCompany(userIDs []int64, companyID int64) error
	UnlinkUserFromCompany(userIDs []int64, companyID int64) error
	GetUserRole(userID, companyID int64) (string, error)
}

type userCompanyRepo struct {
	db *sql.DB
}

// NewUserCompanyRepository создаёт новый репозиторий для работы со связями пользователя и компании.
func NewUserCompanyRepository(db *sql.DB) UserCompanyRepository {
	return &userCompanyRepo{db: db}
}

// LinkUserWithCompany добавляет связь между пользователями и компанией, используя массив userIDs.
// Если такая связь уже существует, обновляется роль.
func (r *userCompanyRepo) LinkUserWithCompany(userIDs []int64, companyID int64) error {
	query := `
		INSERT INTO user_companies (user_id, company_id, role)
		SELECT unnest($1::bigint[]), $2, 'developer'
		ON CONFLICT (user_id, company_id) DO UPDATE SET role = EXCLUDED.role
	`
	_, err := r.db.Exec(query, pq.Array(userIDs), companyID)
	if err != nil {
		return fmt.Errorf("failed to link users with company: %w", err)
	}
	return nil
}

// UnlinkUserFromCompany удаляет связь между пользователями и компанией, используя массив userIDs.
func (r *userCompanyRepo) UnlinkUserFromCompany(userIDs []int64, companyID int64) error {
	query := `
		DELETE FROM user_companies 
		WHERE user_id = ANY($1) AND company_id = $2
	`
	_, err := r.db.Exec(query, pq.Array(userIDs), companyID)
	if err != nil {
		return fmt.Errorf("failed to unlink users from company: %w", err)
	}
	return nil
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
