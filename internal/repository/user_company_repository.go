package repository

import (
	"database/sql"
	"errors"
	"victa/internal/domain"
)

type UserCompanyRepository interface {
	// GetAllByCompanyID возвращает все связи пользователей с указанной компанией
	GetAllByCompanyID(companyID int64) ([]domain.UserCompany, error)
	GetByCompanyAndUserID(companyID, id int64) (*domain.UserCompany, error)
	Delete(userID, companyID int64) error
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
			&uc.UserID,
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

// GetByCompanyAndUserID возвращает связь user_companies для данного companyID и ID.
// Если запись не найдена, возвращает (nil, nil).
func (r *PostgresUserCompanyRepository) GetByCompanyAndUserID(companyID, userID int64) (*domain.UserCompany, error) {
	var uc domain.UserCompany
	err := r.DB.QueryRow(
		`SELECT user_id, company_id, role_id
         FROM user_companies
         WHERE company_id = $1 AND user_id = $2`,
		companyID, userID,
	).Scan(&uc.UserID, &uc.CompanyID, &uc.RoleID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &uc, nil
}

// Delete удаляет связь user_companies для заданного userID и companyID.
// Если записи нет, возвращает sql.ErrNoRows.
func (r *PostgresUserCompanyRepository) Delete(userID, companyID int64) error {
	res, err := r.DB.Exec(
		`DELETE FROM user_companies
         WHERE user_id = $1 AND company_id = $2`,
		userID, companyID,
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
