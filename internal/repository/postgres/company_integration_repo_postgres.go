package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"victa/internal/domain"

	appErr "victa/internal/errors"
)

// CompanyIntegrationRepo хранит prepared‑statements.
type CompanyIntegrationRepo struct {
	db        *sql.DB
	stGetByID *sql.Stmt
	stUpsert  *sql.Stmt
}

// NewCompanyIntegrationRepo инициализирует репозиторий.
func NewCompanyIntegrationRepo(db *sql.DB) (*CompanyIntegrationRepo, error) {
	r := &CompanyIntegrationRepo{db: db}

	var err error
	if r.stGetByID, err = db.Prepare(`
		SELECT company_id,
		       codemagic_api_key,
		       notification_bot_token,
		       deploy_notification_chat_id,
		       issues_notification_chat_id,
		       errors_notification_chat_id
		  FROM company_integrations
		 WHERE company_id = $1`); err != nil {
		return nil, fmt.Errorf("prepare getByID: %w", err)
	}

	if r.stUpsert, err = db.Prepare(`
		INSERT INTO company_integrations (
		      company_id,
		      codemagic_api_key,
		      notification_bot_token,
		      deploy_notification_chat_id,
		      issues_notification_chat_id,
		      errors_notification_chat_id)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (company_id) DO UPDATE
		    SET codemagic_api_key           = EXCLUDED.codemagic_api_key,
		        notification_bot_token      = EXCLUDED.notification_bot_token,
		        deploy_notification_chat_id = EXCLUDED.deploy_notification_chat_id,
		        issues_notification_chat_id = EXCLUDED.issues_notification_chat_id,
		        errors_notification_chat_id = EXCLUDED.errors_notification_chat_id
		RETURNING company_id,
		          codemagic_api_key,
		          notification_bot_token,
		          deploy_notification_chat_id,
		          issues_notification_chat_id,
		          errors_notification_chat_id`); err != nil {
		return nil, fmt.Errorf("prepare upsert: %w", err)
	}

	return r, nil
}

// Close закрывает prepared‑statements.
func (r *CompanyIntegrationRepo) Close() error {
	if r == nil {
		return nil
	}
	if err := r.stGetByID.Close(); err != nil {
		return err
	}
	return r.stUpsert.Close()
}

// GetByID возвращает настройки интеграций или ErrIntegrationNotFound.
func (r *CompanyIntegrationRepo) GetByID(ctx context.Context, companyID int64) (*domain.CompanyIntegration, error) {
	var ci domain.CompanyIntegration
	err := r.stGetByID.QueryRowContext(ctx, companyID).Scan(
		&ci.CompanyID,
		&ci.CodemagicAPIKey,
		&ci.NotificationBotToken,
		&ci.DeployNotificationChatID,
		&ci.IssuesNotificationChatID,
		&ci.ErrorsNotificationChatID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, appErr.ErrIntegrationNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get integration: %w", err)
	}
	return &ci, nil
}

// CreateOrUpdate делает upsert и возвращает актуальные данные.
func (r *CompanyIntegrationRepo) CreateOrUpdate(ctx context.Context, ci *domain.CompanyIntegration) (*domain.CompanyIntegration, error) {
	row := r.stUpsert.QueryRowContext(ctx,
		ci.CompanyID,
		ci.CodemagicAPIKey,
		ci.NotificationBotToken,
		ci.DeployNotificationChatID,
		ci.IssuesNotificationChatID,
		ci.ErrorsNotificationChatID,
	)

	var updated domain.CompanyIntegration
	if err := row.Scan(
		&updated.CompanyID,
		&updated.CodemagicAPIKey,
		&updated.NotificationBotToken,
		&updated.DeployNotificationChatID,
		&updated.IssuesNotificationChatID,
		&updated.ErrorsNotificationChatID,
	); err != nil {
		return nil, fmt.Errorf("upsert integration: %w", err)
	}
	return &updated, nil
}
