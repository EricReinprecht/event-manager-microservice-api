package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"gorm.io/gorm"
)

type TicketCategoryRepository struct {
	db database.DBExecutor
}

func NewTicketCategoryRepository(
	db database.DBExecutor,
) *TicketCategoryRepository {

	return &TicketCategoryRepository{
		db: db,
	}
}

func (r *TicketCategoryRepository) Create(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	return r.db.
		WithContext(ctx).
		Create(category).
		Error()
}

func (r *TicketCategoryRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.TicketCategory, error) {

	var category models.TicketCategory

	err := r.db.
		WithContext(ctx).
		Preload("Party").
		First(&category, "id = ?", id).
		Error()

	return &category, err
}

func (r *TicketCategoryRepository) FindByParty(
	ctx context.Context,
	partyID uuid.UUID,
) ([]models.TicketCategory, error) {

	var categories []models.TicketCategory

	err := r.db.
		WithContext(ctx).
		Where("party_id = ?", partyID).
		Find(&categories).
		Error()

	return categories, err
}

func (r *TicketCategoryRepository) Update(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	tx := r.db.WithContext(ctx).Begin()

	err := tx.
		Where(
			"ticket_category_id = ?",
			category.ID,
		).
		Delete(&models.TicketAccessWindow{}).
		Error()

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.
		Save(category).
		Error()

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *TicketCategoryRepository) Delete(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	return r.db.
		WithContext(ctx).
		Delete(category).
		Error()
}

func (r *TicketCategoryRepository) FindByIDForUpdate(
	ctx context.Context,
	id uuid.UUID,
) (*models.TicketCategory, error) {

	var category models.TicketCategory

	err := r.db.
		WithContext(ctx).
		Raw(
			`
			SELECT *
			FROM ticket_categories
			WHERE id = ?
			AND deleted_at IS NULL
			FOR UPDATE
			`,
			id,
		).
		Scan(&category).
		Error()

	if err != nil {
		return nil, err
	}

	if category.ID == uuid.Nil {
		return nil, gorm.ErrRecordNotFound
	}

	return &category, nil
}

func (r *TicketCategoryRepository) Replace(
	tx database.DBExecutor,
	partyID uuid.UUID,
	categories []models.TicketCategory,
) error {

	var existing []models.TicketCategory

	if err := tx.
		Where(
			"party_id = ?",
			partyID,
		).
		Find(
			&existing,
		).
		Error(); err != nil {

		return err
	}

	keepCategories := make(map[uuid.UUID]bool)

	for _, category := range categories {

		//
		// CREATE CATEGORY
		//
		if category.ID == uuid.Nil {

			category.ID = uuid.New()

			category.PartyID = partyID

			keepCategories[category.ID] = true

			windows := category.AccessWindows

			category.AccessWindows = nil

			if err := tx.
				Create(
					&category,
				).
				Error(); err != nil {

				return err
			}

			for _, window := range windows {

				window.ID = uuid.New()

				window.TicketCategoryID = category.ID

				if err := tx.
					Create(
						&window,
					).
					Error(); err != nil {

					return err
				}
			}

			continue
		}

		//
		// UPDATE CATEGORY
		//
		keepCategories[category.ID] = true

		if err := tx.
			Model(
				&models.TicketCategory{},
			).
			Where(
				"id = ?",
				category.ID,
			).
			Updates(
				map[string]interface{}{
					"name":                     category.Name,
					"price":                    category.Price,
					"capacity":                 category.Capacity,
					"requires_verification":    category.RequiresVerification,
					"refund_requires_approval": category.RefundRequiresApproval,
					"refund_policy_id":         category.RefundPolicyID,
				},
			).
			Error(); err != nil {

			return err
		}

		//
		// ACCESS WINDOWS
		//
		var existingWindows []models.TicketAccessWindow

		if err := tx.
			Where(
				"ticket_category_id = ?",
				category.ID,
			).
			Find(
				&existingWindows,
			).
			Error(); err != nil {

			return err
		}

		keepWindows := make(map[uuid.UUID]bool)

		for _, window := range category.AccessWindows {

			//
			// CREATE WINDOW
			//
			if window.ID == uuid.Nil {

				window.ID = uuid.New()

				window.TicketCategoryID = category.ID

				if err := tx.
					Create(
						&window,
					).
					Error(); err != nil {

					return err
				}

				continue
			}

			//
			// UPDATE WINDOW
			//
			keepWindows[window.ID] = true

			if err := tx.
				Model(
					&models.TicketAccessWindow{},
				).
				Where(
					"id = ?",
					window.ID,
				).
				Updates(
					map[string]interface{}{
						"starts_at": window.StartsAt,
						"ends_at":   window.EndsAt,
					},
				).
				Error(); err != nil {

				return err
			}
		}

		//
		// DELETE REMOVED WINDOWS
		//
		for _, oldWindow := range existingWindows {

			if !keepWindows[oldWindow.ID] {

				if err := tx.
					Unscoped().
					Delete(
						&models.TicketAccessWindow{},
						"id = ?",
						oldWindow.ID,
					).
					Error(); err != nil {

					return err
				}
			}
		}
	}

	//
	// DELETE REMOVED CATEGORIES
	//
	for _, oldCategory := range existing {

		if !keepCategories[oldCategory.ID] {

			if err := tx.
				Unscoped().
				Where(
					"ticket_category_id = ?",
					oldCategory.ID,
				).
				Delete(
					&models.TicketAccessWindow{},
				).
				Error(); err != nil {

				return err
			}

			if err := tx.
				Unscoped().
				Delete(
					&models.TicketCategory{},
					"id = ?",
					oldCategory.ID,
				).
				Error(); err != nil {

				return err
			}
		}
	}

	return nil
}
