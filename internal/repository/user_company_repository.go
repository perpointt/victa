package repository

import (
	"database/sql"
	"victa/internal/domain"
)

type UserCompanyRepository interface {
	// GetAllByCompanyID возвращает все связи пользователей с указанной компанией
	GetAllByCompanyID(companyID int64) ([]domain.UserCompany, error)
}

type PostgresUserCompanyRepository struct {
	DB *sql.DB
}

func NewPostgresUserCompanyRepository(db *sql.DB) *PostgresUserCompanyRepository {
	return &PostgresUserCompanyRepository{DB: db}
}

// GetAllByCompanyID возвращает все записи из user_companies для companyID
func (r *PostgresUserCompanyRepository) GetAllByCompanyID(companyID int64) ([]domain.UserCompany, error) {
	rows, err := r.DB.Query(
		`SELECT user_id, company_id, role_id
         FROM user_companies
         WHERE company_id = $1
         ORDER BY role_id, user_id`,
		companyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.UserCompany
	for rows.Next() {
		var uc domain.UserCompany
		if err := rows.Scan(
			&uc.UserId,
			&uc.CompanyID,
			&uc.RoleID,
		); err != nil {
			return nil, err
		}
		list = append(list, uc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}
