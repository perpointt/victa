package repository

import (
	"database/sql"
	"errors"
	"victa/internal/domain"
)

// CompanyIntegrationRepository определяет методы для работы с company_integrations.
type CompanyIntegrationRepository interface {
	// GetByID возвращает настройки интеграций для заданной компании или nil, если их нет.
	GetByID(companyID int64) (*domain.CompanyIntegration, error)
	// CreateOrUpdate создаёт или обновляет настройки интеграций компании и возвращает их.
	CreateOrUpdate(ci *domain.CompanyIntegration) (*domain.CompanyIntegration, error)
}

// PostgresCompanyIntegrationRepo реализует CompanyIntegrationRepository через Postgres.
type PostgresCompanyIntegrationRepo struct {
	DB *sql.DB
}

// NewPostgresCompanyIntegrationRepo создаёт репозиторий company_integrations.
func NewPostgresCompanyIntegrationRepo(db *sql.DB) *PostgresCompanyIntegrationRepo {
	return &PostgresCompanyIntegrationRepo{DB: db}
}

// GetByID возвращает настройки интеграций по company_id.
func (r *PostgresCompanyIntegrationRepo) GetByID(companyID int64) (*domain.CompanyIntegration, error) {
	var ci domain.CompanyIntegration
	err := r.DB.QueryRow(
		`SELECT company_id, codemagic_api_key, notification_bot_token, deploy_notification_chat_id, issues_notification_chat_id
         FROM company_integrations
         WHERE company_id = $1`,
		companyID,
	).Scan(
		&ci.CompanyID,
		&ci.CodemagicAPIKey,
		&ci.NotificationBotToken,
		&ci.DeployNotificationChatID,
		&ci.IssuesNotificationChatID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return &ci, nil
}

// CreateOrUpdate вставляет или обновляет запись company_integrations.
func (r *PostgresCompanyIntegrationRepo) CreateOrUpdate(ci *domain.CompanyIntegration) (*domain.CompanyIntegration, error) {
	row := r.DB.QueryRow(
		`INSERT INTO company_integrations
         (company_id, codemagic_api_key, notification_bot_token, deploy_notification_chat_id, issues_notification_chat_id)
     VALUES ($1, $2, $3, $4, $5)
     ON CONFLICT (company_id) DO UPDATE
       SET codemagic_api_key            = EXCLUDED.codemagic_api_key,
           notification_bot_token       = EXCLUDED.notification_bot_token,
           deploy_notification_chat_id  = EXCLUDED.deploy_notification_chat_id,
           issues_notification_chat_id  = EXCLUDED.issues_notification_chat_id
     RETURNING company_id, codemagic_api_key, notification_bot_token, deploy_notification_chat_id, issues_notification_chat_id`,
		ci.CompanyID,
		ci.CodemagicAPIKey,
		ci.NotificationBotToken,
		ci.DeployNotificationChatID,
		ci.IssuesNotificationChatID,
	)

	var updated domain.CompanyIntegration
	if err := row.Scan(
		&updated.CompanyID,
		&updated.CodemagicAPIKey,
		&updated.NotificationBotToken,
		&updated.DeployNotificationChatID,
		&updated.IssuesNotificationChatID,
	); err != nil {
		return nil, err
	}
	return &updated, nil
}
