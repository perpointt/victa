package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"victa/internal/domain"
)

// UserRepository описывает методы работы с пользователями.
type UserRepository interface {
	CreateUserWithCompany(user *domain.User, companyID *int64) error
	GetAll() ([]domain.User, error)
	GetUsersByCompanyID(companyID int64) ([]domain.User, error)
	GetByID(id int64) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id int64) error
	GetByEmail(email string) (*domain.User, error)
}

type userRepo struct {
	db *sql.DB
}

// NewUserRepository создаёт новую реализацию UserRepository.
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUserWithCompany(user *domain.User, companyID *int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	// Helper для отката транзакции в случае ошибки.
	rollback := func(err error) error {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback error: %v; original error: %w", rbErr, err)
		}
		return err
	}

	// 1. Проверяем, существует ли пользователь с таким email.
	var userExists bool
	checkUserQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	if err = tx.QueryRow(checkUserQuery, user.Email).Scan(&userExists); err != nil {
		return rollback(fmt.Errorf("failed to check user existence: %w", err))
	}
	if userExists {
		return rollback(fmt.Errorf("user with email %s already exists", user.Email))
	}

	// 2. Если companyID передан, проверяем, существует ли такая компания.
	if companyID != nil {
		var companyExists bool
		checkCompanyQuery := `SELECT EXISTS(SELECT 1 FROM companies WHERE id = $1)`
		if err = tx.QueryRow(checkCompanyQuery, *companyID).Scan(&companyExists); err != nil {
			return rollback(fmt.Errorf("failed to check company existence: %w", err))
		}
		if !companyExists {
			return rollback(fmt.Errorf("company with id %d does not exist", *companyID))
		}
	}

	// 3. Создаем пользователя.
	createQuery := `
		INSERT INTO users (email, password, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	if err = tx.QueryRow(createQuery, user.Email, user.Password).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return rollback(fmt.Errorf("failed to create user: %w", err))
	}

	// 4. Если companyID передан, связываем пользователя с компанией с ролью "developer".
	if companyID != nil {
		linkQuery := `
			INSERT INTO user_companies (user_id, company_id, role)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, company_id) DO UPDATE SET role = EXCLUDED.role
		`
		if _, err = tx.Exec(linkQuery, user.ID, *companyID, "developer"); err != nil {
			return rollback(fmt.Errorf("failed to link user with company: %w", err))
		}
	}

	// Фиксируем транзакцию.
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *userRepo) GetAll() ([]domain.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// GetAllByCompanyID возвращает список пользователей, связанных с заданной компанией.
func (r *userRepo) GetUsersByCompanyID(companyID int64) ([]domain.User, error) {
	query := `
		SELECT u.id, u.email, u.password, u.created_at, u.updated_at
		FROM users u
		INNER JOIN user_companies uc ON u.id = uc.user_id
		WHERE uc.company_id = $1
	`
	rows, err := r.db.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepo) GetByID(id int64) (*domain.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE id = $1`
	var user domain.User
	err := r.db.QueryRow(query, id).
		Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, password = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`
	return r.db.QueryRow(query, user.Email, user.Password, user.ID).
		Scan(&user.UpdatedAt)
}

func (r *userRepo) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *userRepo) GetByEmail(email string) (*domain.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`
	var user domain.User
	err := r.db.QueryRow(query, email).
		Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
