package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type RefreshTokenRepository struct {
	db database.DBExecutor
}

func NewRefreshTokenRepository(
	db database.DBExecutor,
) *RefreshTokenRepository {

	return &RefreshTokenRepository{
		db: db,
	}
}

func (r *RefreshTokenRepository) Create(
	token *models.RefreshToken,
) error {

	return r.db.
		Create(token).
		Error()
}

func (r *RefreshTokenRepository) FindActiveByHash(
	ctx context.Context,
	hash string,
) (*models.RefreshToken, error) {

	var token models.RefreshToken

	err := r.db.
		WithContext(ctx).
		Where(
			"token_hash = ? AND revoked_at IS NULL AND expires_at > ?",
			hash,
			time.Now(),
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

func (r *RefreshTokenRepository) FindByHash(
	ctx context.Context,
	hash string,
) (*models.RefreshToken, error) {

	var token models.RefreshToken

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

func (r *RefreshTokenRepository) Revoke(
	id uuid.UUID,
) error {

	return r.db.
		Model(
			&models.RefreshToken{},
		).
		Where(
			"id = ?",
			id,
		).
		Updates(
			map[string]any{
				"revoked_at": time.Now(),
			},
		).
		Error()
}

func (r *RefreshTokenRepository) Transaction(
	ctx context.Context,
	fn func(tx database.DBExecutor) error,
) error {

	tx := r.db.
		WithContext(ctx).
		Begin()

	if err := tx.Error(); err != nil {
		return err
	}

	if err := fn(tx); err != nil {

		_ = tx.Rollback()

		return err
	}

	return tx.Commit()
}

func (r *RefreshTokenRepository) RevokeByHash(
	ctx context.Context,
	hash string,
) error {

	now := time.Now()

	return r.db.
		WithContext(ctx).
		Model(
			&models.RefreshToken{},
		).
		Where(
			"token_hash = ?",
			hash,
		).
		Updates(
			map[string]any{
				"revoked_at": now,
			},
		).
		Error()
}

func (r *RefreshTokenRepository) RevokeAllForUser(
	userID uuid.UUID,
) error {

	now := time.Now()

	return r.db.
		Model(
			&models.RefreshToken{},
		).
		Where(
			"user_id = ? AND revoked_at IS NULL",
			userID,
		).
		Updates(
			map[string]any{
				"revoked_at": now,
			},
		).
		Error()
}

func (r *RefreshTokenRepository) RevokeFamily(
	ctx context.Context,
	familyID uuid.UUID,
) error {

	return r.db.
		WithContext(ctx).
		Model(
			&models.RefreshToken{},
		).
		Where(
			"family_id = ?",
			familyID,
		).
		Updates(
			map[string]any{
				"revoked_at": time.Now(),
			},
		).
		Error()
}

func (r *RefreshTokenRepository) FindSessionsByUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]models.RefreshToken, error) {

	var tokens []models.RefreshToken

	err := r.db.
		WithContext(ctx).
		Where(
			"user_id = ?",
			userID,
		).
		Order(
			"created_at DESC",
		).
		Find(
			&tokens,
		).
		Error()

	return tokens, err
}

func (r *RefreshTokenRepository) RevokeAllByUser(
	ctx context.Context,
	userID uuid.UUID,
) error {

	return r.db.
		WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where(
			"user_id = ? AND revoked_at IS NULL",
			userID,
		).
		Updates(
			map[string]any{
				"revoked_at": time.Now(),
			},
		).
		Error()
}
