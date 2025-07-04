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
	// GetByID возвращает пользователя по внутреннему ID или nil, если не найден
	GetByID(id int64) (*domain.User, error)
	// GetByTgID возвращает пользователя по Telegram ID или nil, если не найден
	GetByTgID(tgID int64) (*domain.User, error)
	// GetAllByCompanyID возвращает всех пользователей, связанных с указанной компанией
	GetAllByCompanyID(companyID int64) ([]domain.User, error)
}

// PostgresUserRepo реализует UserRepository через Postgres
type PostgresUserRepo struct {
	DB *sql.DB
}

// NewPostgresUserRepo создаёт репозиторий пользователей
func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{DB: db}
}

// Create вставляет нового пользователя и возвращает сущность
func (r *PostgresUserRepo) Create(u *domain.User) (*domain.User, error) {
	newUser := &domain.User{}
	err := r.DB.QueryRow(
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
	return newUser, nil
}

// GetByID ищет пользователя по внутреннему ID
func (r *PostgresUserRepo) GetByID(id int64) (*domain.User, error) {
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

// GetByTgID ищет пользователя по Telegram ID
func (r *PostgresUserRepo) GetByTgID(tgID int64) (*domain.User, error) {
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

// GetAllByCompanyID возвращает всех пользователей, которые состоят в компании companyID,
// отсортированных по дате создания (сначала самые новые)
func (r *PostgresUserRepo) GetAllByCompanyID(companyID int64) ([]domain.User, error) {
	rows, err := r.DB.Query(
		`SELECT u.id, u.tg_id, u.name, u.created_at, u.updated_at
         FROM users u
         JOIN user_companies uc ON u.id = uc.user_id
         WHERE uc.company_id = $1
         ORDER BY u.created_at DESC`, // сортировка по дате создания
		companyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.ID,
			&u.TgId,
			&u.Name,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
