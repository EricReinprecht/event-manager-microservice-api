package repository

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type UserRepository struct {
	db database.DBExecutor
}

func NewUserRepository(
	db database.DBExecutor,
) *UserRepository {

	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *models.User,
) error {

	return r.db.
		WithContext(ctx).
		Create(user).
		Error()

}

func (r *UserRepository) Update(
	ctx context.Context,
	user *models.User,
) error {

	return r.db.
		WithContext(ctx).
		Save(user).
		Error()
}

func (r *UserRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.User, error) {

	var user models.User

	err := r.db.
		WithContext(ctx).
		First(
			&user,
			"id = ?",
			id,
		).
		Error()

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (r *UserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {

	var user models.User

	err := r.db.
		WithContext(ctx).
		Where(
			"email = ?",
			email,
		).
		First(
			&user,
		).
		Error()

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (r *UserRepository) FindByUsername(
	ctx context.Context,
	username string,
) (*models.User, error) {

	var user models.User

	err := r.db.
		WithContext(ctx).
		Where(
			"username = ?",
			username,
		).
		First(
			&user,
		).
		Error()

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByIdentifier(
	ctx context.Context,
	identifier string,
) (*models.User, error) {

	identifier = strings.TrimSpace(
		strings.ToLower(identifier),
	)

	var user models.User

	err := r.db.
		WithContext(ctx).
		Where(
			"LOWER(email) = ? OR LOWER(username) = ?",
			identifier,
			identifier,
		).
		First(&user).
		Error()

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Transaction(
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

func (r *UserRepository) Delete(
	ctx context.Context,
	user *models.User,
) error {

	return r.db.
		WithContext(ctx).
		Delete(user).
		Error()
}
