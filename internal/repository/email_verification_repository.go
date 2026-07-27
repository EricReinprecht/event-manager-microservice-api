package repository

import (
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
	token string,
) (*models.EmailVerification, error) {

	var verification models.EmailVerification

	err := r.executor.
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
