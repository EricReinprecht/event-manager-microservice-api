package password_reset_token_repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type PasswordResetTokenRepository struct {
	db database.DBExecutor
}

func NewPasswordResetTokenRepository(
	db database.DBExecutor,
) *PasswordResetTokenRepository {

	return &PasswordResetTokenRepository{
		db: db,
	}
}

func (r *PasswordResetTokenRepository) Create(
	ctx context.Context,
	token *models.PasswordResetToken,
) error {

	return r.db.
		WithContext(ctx).
		Create(
			token,
		).
		Error()
}

func (r *PasswordResetTokenRepository) FindByHash(
	ctx context.Context,
	hash string,
) (*models.PasswordResetToken, error) {

	var token models.PasswordResetToken

	err := r.db.
		WithContext(ctx).
		Where(
			"token_hash = ?",
			hash,
		).
		First(
			&token,
		).
		Error()

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *PasswordResetTokenRepository) MarkUsed(
	ctx context.Context,
	id uuid.UUID,
) error {

	now := time.Now()

	return r.db.
		WithContext(ctx).
		Model(
			&models.PasswordResetToken{},
		).
		Where(
			"id = ?",
			id,
		).
		Updates(
			map[string]interface{}{
				"used_at": &now,
			},
		).
		Error()
}

func (r *PasswordResetTokenRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {

	return r.db.
		WithContext(ctx).
		Delete(
			&models.PasswordResetToken{
				ID: id,
			},
		).
		Error()
}

func (r *PasswordResetTokenRepository) DeleteForUser(
	ctx context.Context,
	userID uuid.UUID,
) error {

	return r.db.
		WithContext(ctx).
		Where(
			"user_id = ?",
			userID,
		).
		Delete(
			&models.PasswordResetToken{},
		).
		Error()
}

func (r *PasswordResetTokenRepository) InvalidateForUser(
	ctx context.Context,
	userID uuid.UUID,
) error {

	now := time.Now()

	return r.db.
		WithContext(ctx).
		Model(
			&models.PasswordResetToken{},
		).
		Where(
			"user_id = ? AND used_at IS NULL AND invalidated_at IS NULL",
			userID,
		).
		Updates(
			map[string]interface{}{
				"invalidated_at": &now,
			},
		).
		Error()
}

func (r *PasswordResetTokenRepository) FindLatestByUser(
	ctx context.Context,
	userID uuid.UUID,
) (*models.PasswordResetToken, error) {

	var token models.PasswordResetToken

	err := r.db.
		WithContext(ctx).
		Where(
			"user_id = ?",
			userID,
		).
		Order(
			"created_at DESC",
		).
		First(
			&token,
		).
		Error()

	if err != nil {
		return nil, err
	}

	return &token, nil
}
