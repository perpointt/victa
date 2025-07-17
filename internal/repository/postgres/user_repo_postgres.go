package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"victa/internal/domain"
	appErr "victa/internal/errors"
)

// UserRepo реализует UserRepository через database/sql + prepared statements
type UserRepo struct {
	db                  *sql.DB
	stCreate            *sql.Stmt
	stUpdate            *sql.Stmt
	stGetByID           *sql.Stmt
	stGetByTgID         *sql.Stmt
	stGetAllByCompanyID *sql.Stmt
}

// NewUserRepo инициализирует репозиторий.
func NewUserRepo(db *sql.DB) (*UserRepo, error) {
	r := &UserRepo{db: db}

	var err error
	if r.stCreate, err = db.Prepare(`
		INSERT INTO users (tg_id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $3)
		RETURNING id, tg_id, name, created_at, updated_at`); err != nil {
		return nil, fmt.Errorf("prepare create: %w", err)
	}

	if r.stUpdate, err = db.Prepare(`
		UPDATE users
		SET name = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, tg_id, name, created_at, updated_at`); err != nil {
		return nil, fmt.Errorf("prepare update: %w", err)
	}

	if r.stGetByID, err = db.Prepare(`
		SELECT id, tg_id, name, created_at, updated_at
		FROM users
		WHERE id = $1`); err != nil {
		return nil, fmt.Errorf("prepare getByID: %w", err)
	}

	if r.stGetByTgID, err = db.Prepare(`
		SELECT id, tg_id, name, created_at, updated_at
		FROM users
		WHERE tg_id = $1`); err != nil {
		return nil, fmt.Errorf("prepare getByTgID: %w", err)
	}

	if r.stGetAllByCompanyID, err = db.Prepare(`
		SELECT u.id, u.tg_id, u.name, u.created_at, u.updated_at
		FROM users u
		JOIN user_companies uc ON u.id = uc.user_id
		WHERE uc.company_id = $1
		ORDER BY u.created_at DESC`); err != nil {
		return nil, fmt.Errorf("prepare getAllByCompanyID: %w", err)
	}

	return r, nil
}

// Close освобождает ресурсы prepared-statement.
func (r *UserRepo) Close() error {
	if r == nil {
		return nil
	}
	for _, st := range []*sql.Stmt{
		r.stCreate, r.stUpdate, r.stGetByID,
		r.stUpdate, r.stGetAllByCompanyID,
	} {
		if st != nil {
			if err := st.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Create сохраняет нового пользователя
func (r *UserRepo) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	now := time.Now().UTC()
	result := new(domain.User)

	if err := r.stCreate.QueryRowContext(ctx, user.TgID, user.Name, now).
		Scan(&result.ID, &result.TgID, &result.Name, &result.CreatedAt, &result.UpdatedAt); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return result, nil
}

// Update изменяет имя пользователя
func (r *UserRepo) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	updated := new(domain.User)
	err := r.stUpdate.QueryRowContext(ctx, user.Name, time.Now().UTC(), user.ID).
		Scan(&updated.ID, &updated.TgID, &updated.Name, &updated.CreatedAt, &updated.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, appErr.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return updated, nil
}

// GetByID возвращает пользователя по внутреннему ID
func (r *UserRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	u := new(domain.User)
	err := r.stGetByID.QueryRowContext(ctx, id).
		Scan(&u.ID, &u.TgID, &u.Name, &u.CreatedAt, &u.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, appErr.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return u, nil
}

// GetByTgID возвращает пользователя по Telegram ID
func (r *UserRepo) GetByTgID(ctx context.Context, tgID int64) (*domain.User, error) {
	u := new(domain.User)
	err := r.stGetByTgID.QueryRowContext(ctx, tgID).
		Scan(&u.ID, &u.TgID, &u.Name, &u.CreatedAt, &u.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, appErr.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by tgID: %w", err)
	}
	return u, nil
}

// GetAllByCompanyID список сотрудников компании, новые сверху
func (r *UserRepo) GetAllByCompanyID(ctx context.Context, companyID int64) ([]domain.User, error) {
	rows, err := r.stGetAllByCompanyID.QueryContext(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("query users by company: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	users := make([]domain.User, 0, 16) // средний отдел ≈ 10–15 чел.
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.TgID, &u.Name, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return users, nil
}
