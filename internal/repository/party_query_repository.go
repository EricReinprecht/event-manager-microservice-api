package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type PartyQueryRepository struct {
	db database.DBExecutor
}

func NewPartyQueryRepository(
	db database.DBExecutor,
) *PartyQueryRepository {

	return &PartyQueryRepository{
		db: db,
	}
}

func (r *PartyQueryRepository) FindAll(
	ctx context.Context,
) ([]models.Party, error) {

	var parties []models.Party

	err := r.db.
		WithContext(ctx).
		Preload("Organizer").
		Preload("Category").
		Preload("Thumbnail").
		Find(&parties).
		Error()

	return parties, err
}

func (r *PartyQueryRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Party, error) {

	var party models.Party

	err := r.db.
		WithContext(ctx).
		Preload("Organizer").
		Preload("Category").
		Preload("Thumbnail").
		Preload("Images").
		First(
			&party,
			"id = ?",
			id,
		).
		Error()

	if err != nil {
		return nil, err
	}

	return &party, nil
}

func (r *PartyQueryRepository) FindForUser(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	startAt string,
	endAt string,
	role string,
	page int,
	limit int,
) ([]models.Party, int64, error) {

	var parties []models.Party

	var total int64

	query := r.db.
		WithContext(ctx).
		Model(&models.Party{})

	switch role {

	case "organized":

		query = query.Where(
			"organizer_id = ?",
			userID,
		)

	case "member":

		query = query.
			Joins(
				"JOIN party_members ON party_members.party_id = parties.id",
			).
			Where(
				"party_members.user_id = ?",
				userID,
			)
	}

	if name != "" {

		query = query.Where(
			"title ILIKE ?",
			"%"+name+"%",
		)
	}

	if startAt != "" {

		query = query.Where(
			"start_at >= ?",
			startAt,
		)
	}

	if endAt != "" {

		query = query.Where(
			"end_at <= ?",
			endAt,
		)
	}

	if err := query.Count(&total).Error(); err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	err := query.
		Order("start_at DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Preload("Category").
		Preload("Thumbnail").
		Find(&parties).
		Error()

	return parties, total, err
}
