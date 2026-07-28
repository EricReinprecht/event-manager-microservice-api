package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type EmailVerificationRepository struct {
	executor database.DBExecutor
}

func NewEmailVerificationRepository(
	executor database.DBExecutor,
) *EmailVerificationRepository {

	return &EmailVerificationRepository{
		executor: executor,
	}

}

func (r *EmailVerificationRepository) Create(
	verification *models.EmailVerification,
) error {

	return r.executor.
		Create(
			verification,
		).
		Error()

}

func (r *EmailVerificationRepository) FindByToken(
	ctx context.Context,
	token string,
) (*models.EmailVerification, error) {

	var verification models.EmailVerification

	err := r.executor.
		WithContext(ctx).
		Where(
			"token = ?",
			token,
		).
		First(
			&verification,
		).
		Error()

	if err != nil {
		return nil, err
	}

	return &verification, nil
}

func (r *EmailVerificationRepository) DeleteExpired() error {

	return r.executor.
		Where(
			"expires_at < ?",
			time.Now(),
		).
		Delete(
			&models.EmailVerification{},
		).
		Error()

}

func (r *EmailVerificationRepository) Delete(
	id uuid.UUID,
) error {

	return r.executor.
		Where(
			"id = ?",
			id,
		).
		Delete(
			&models.EmailVerification{},
		).
		Error()
}

func (r *EmailVerificationRepository) Update(
	verification *models.EmailVerification,
) error {

	return r.executor.
		Save(
			verification,
		).
		Error()
}

func (r *EmailVerificationRepository) DeleteByUser(
	userID uuid.UUID,
) error {

	return r.executor.
		Where(
			"user_id = ?",
			userID,
		).
		Delete(
			&models.EmailVerification{},
		).
		Error()
}

func (r *EmailVerificationRepository) InvalidateForUser(
	ctx context.Context,
	userID uuid.UUID,
) error {

	now := time.Now()

	return r.executor.
		WithContext(ctx).
		Model(&models.EmailVerification{}).
		Where(
			"user_id = ? AND used_at IS NULL",
			userID,
		).
		Updates(
			map[string]interface{}{
				"used_at": &now,
			},
		).
		Error()
}

func (r *EmailVerificationRepository) FindLatestByUser(
	ctx context.Context,
	userID uuid.UUID,
) (*models.EmailVerification, error) {

	var verification models.EmailVerification

	err := r.executor.
		WithContext(ctx).
		Where(
			"user_id = ?",
			userID,
		).
		Order(
			"created_at DESC",
		).
		First(
			&verification,
		).
		Error()

	if err != nil {
		return nil, err
	}

	return &verification, nil
}
