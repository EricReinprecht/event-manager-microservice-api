package repository

import (
	"context"
	"strings"

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
		Preload("Categories").
		Preload("Thumbnail").
		Preload("Images").
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
		Preload("Categories").
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
		Model(&models.Party{}).
		Distinct("parties.*")

	switch role {

	case "organized":

		query = query.Where(
			"parties.organizer_id = ?",
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
			"parties.title ILIKE ?",
			"%"+name+"%",
		)
	}

	if startAt != "" {

		query = query.Where(
			"parties.start_at >= ?",
			startAt,
		)
	}

	if endAt != "" {

		query = query.Where(
			"parties.end_at <= ?",
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
		Order(
			"parties.start_at DESC",
		).
		Limit(limit).
		Offset((page - 1) * limit).
		Preload("Organizer").
		Preload("Categories").
		Preload("Thumbnail").
		Preload("Images").
		Find(&parties).
		Error()

	return parties, total, err
}

func (r *PartyQueryRepository) FindOrganizedByUser(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	startAt string,
	endAt string,
	sorts string,
	page int,
	limit int,
) ([]models.Party, int64, error) {

	var parties []models.Party
	var total int64

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	query := r.db.
		WithContext(ctx).
		Model(&models.Party{}).
		Where(
			"parties.organizer_id = ?",
			userID,
		)

	if name != "" {

		query = query.Where(
			"LOWER(parties.title) LIKE ?",
			"%"+strings.ToLower(name)+"%",
		)
	}

	if startAt != "" {

		query = query.Where(
			"(parties.start_at AT TIME ZONE parties.timezone)::date = ?",
			startAt,
		)
	}

	if endAt != "" {

		query = query.Where(
			"(parties.end_at AT TIME ZONE parties.timezone)::date = ?",
			endAt,
		)
	}

	if err := query.Count(&total).Error(); err != nil {
		return nil, 0, err
	}

	query = applyPartySorts(
		query,
		sorts,
	)

	err := query.
		Preload("Organizer").
		Preload("Categories").
		Preload("Thumbnail").
		Preload("Images").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&parties).
		Error()

	return parties, total, err
}

func applyPartySorts(
	query database.DBExecutor,
	sorts string,
) database.DBExecutor {

	allowed := map[string]string{
		"title":    "parties.title",
		"startAt":  "parties.start_at",
		"endAt":    "parties.end_at",
		"location": "parties.location_name",
	}

	if sorts == "" {

		return query.Order(
			"parties.start_at DESC",
		)
	}

	for _, sort := range strings.Split(
		sorts,
		",",
	) {

		parts := strings.Split(
			sort,
			":",
		)

		if len(parts) != 2 {
			continue
		}

		column, ok := allowed[parts[0]]

		if !ok {
			continue
		}

		direction := strings.ToLower(parts[1])

		if direction != "asc" &&
			direction != "desc" {

			direction = "asc"
		}

		query = query.Order(
			column + " " + direction,
		)
	}

	return query
}
