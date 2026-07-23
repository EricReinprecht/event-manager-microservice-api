package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/models"
)

type PartyRepository struct {
	db *gorm.DB
}

func NewPartyRepository(
	db *gorm.DB,
) *PartyRepository {

	return &PartyRepository{
		db: db,
	}
}

func (r *PartyRepository) Create(
	ctx context.Context,
	party *models.Party,
) error {

	return r.db.
		WithContext(ctx).
		Create(party).
		Error
}

func (r *PartyRepository) FindAll(
	ctx context.Context,
) ([]models.Party, error) {

	var parties []models.Party

	err := r.db.
		WithContext(ctx).
		Preload("Organizer").
		Find(&parties).
		Error

	return parties, err
}

func (r *PartyRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Party, error) {

	var party models.Party

	err := r.db.
		WithContext(ctx).
		Preload("Organizer").
		First(&party, id).
		Error

	return &party, err
}

func (r *PartyRepository) Update(
	ctx context.Context,
	party *models.Party,
) error {

	return r.db.
		WithContext(ctx).
		Save(party).
		Error
}

func (r *PartyRepository) Delete(
	ctx context.Context,
	party *models.Party,
) error {

	return r.db.
		WithContext(ctx).
		Delete(party).
		Error
}
