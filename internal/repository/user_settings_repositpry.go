package repository

import (
	"database/sql"
	"errors"

	"victa/internal/domain"
)

// UserSettingsRepository определяет методы для работы с настройками пользователя
type UserSettingsRepository interface {
	Create(s *domain.UserSettings) (*domain.UserSettings, error)
	FindByUserID(userID int64) (*domain.UserSettings, error)
}

// PostgresUserSettingsRepo реализует UserSettingsRepository через Postgres
type PostgresUserSettingsRepo struct {
	DB *sql.DB
}

// NewPostgresUserSettingsRepo создаёт репозиторий настроек пользователя
func NewPostgresUserSettingsRepo(db *sql.DB) *PostgresUserSettingsRepo {
	return &PostgresUserSettingsRepo{DB: db}
}

// Create вставляет новые настройки пользователя и возвращает сущность
func (r *PostgresUserSettingsRepo) Create(s *domain.UserSettings) (*domain.UserSettings, error) {
	var (
		newSettings domain.UserSettings
		acNull      sql.NullInt64
	)
	err := r.DB.QueryRow(
		`INSERT INTO user_settings (user_id, active_company_id)
VALUES ($1, $2)
RETURNING user_id, active_company_id`,
		s.UserId, s.ActiveCompanyId,
	).Scan(
		&newSettings.UserId,
		&acNull,
	)
	if err != nil {
		return nil, err
	}
	if acNull.Valid {
		v := acNull.Int64
		newSettings.ActiveCompanyId = &v
	}
	return &newSettings, nil
}

// FindByUserID возвращает настройки для данного user_id или nil, если не найдено
func (r *PostgresUserSettingsRepo) FindByUserID(userID int64) (*domain.UserSettings, error) {
	var (
		settings domain.UserSettings
		acNull   sql.NullInt64
	)
	err := r.DB.QueryRow(
		`SELECT user_id, active_company_id
FROM user_settings
WHERE user_id = $1`,
		userID,
	).Scan(
		&settings.UserId,
		&acNull,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if acNull.Valid {
		v := acNull.Int64
		settings.ActiveCompanyId = &v
	}
	return &settings, nil
}
