package repository

import (
	"database/sql"
	"errors"

	"victa/internal/domain"
)

// UserRepository определяет методы для работе с пользователями
type UserRepository interface {
	// Create сохраняет нового пользователя и возвращает созданную сущность
	Create(u *domain.User) (*domain.User, error)
	// FindByID возвращает пользователя по внутреннему ID или nil, если не найден
	FindByID(id int64) (*domain.User, error)
	// FindByTgID возвращает пользователя по Telegram ID или nil, если не найден
	FindByTgID(tgID int64) (*domain.User, error)
}

// PostgresUserRepo реализует UserRepository через Postgres
type PostgresUserRepo struct {
	DB *sql.DB
}

// NewPostgresUserRepo создаёт репозиторий пользователей
func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{DB: db}
}

// Create вставляет нового пользователя и его настройки в рамках одной транзакции
func (r *PostgresUserRepo) Create(u *domain.User) (*domain.User, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	newUser := &domain.User{}
	err = tx.QueryRow(
		`INSERT INTO users (tg_id, name, created_at, updated_at)
         VALUES ($1, $2, now(), now())
         RETURNING id, tg_id, name, created_at, updated_at`,
		u.TgId, u.Name,
	).Scan(
		&newUser.ID,
		&newUser.TgId,
		&newUser.Name,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(
		`INSERT INTO user_settings (user_id, active_company_id)
         VALUES ($1, NULL)`,
		newUser.ID,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return newUser, nil
}

// FindByID ищет пользователя по внутреннему ID
func (r *PostgresUserRepo) FindByID(id int64) (*domain.User, error) {
	u := &domain.User{}
	err := r.DB.QueryRow(
		`SELECT id, tg_id, name, created_at, updated_at
         FROM users
         WHERE id = $1`,
		id,
	).Scan(
		&u.ID,
		&u.TgId,
		&u.Name,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// FindByTgID ищет пользователя по Telegram ID
func (r *PostgresUserRepo) FindByTgID(tgID int64) (*domain.User, error) {
	u := &domain.User{}
	err := r.DB.QueryRow(
		`SELECT id, tg_id, name, created_at, updated_at
         FROM users
         WHERE tg_id = $1`,
		tgID,
	).Scan(
		&u.ID,
		&u.TgId,
		&u.Name,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}
