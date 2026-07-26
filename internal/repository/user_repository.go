package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type UserRepository struct {
	executor database.DBExecutor
}

func NewUserRepository(
	executor database.DBExecutor,
) *UserRepository {

	return &UserRepository{
		executor: executor,
	}

}

func (r *UserRepository) Create(
	ctx context.Context,
	user *models.User,
) error {

	return r.executor.
		WithContext(ctx).
		Create(user).
		Error()

}

func (r *UserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {

	var user models.User

	err := r.executor.
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

func (r *UserRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.User, error) {

	var user models.User

	err := r.executor.
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
