package repository

import (
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type TicketCategoryWriteRepository struct {
}

func NewTicketCategoryWriteRepository() *TicketCategoryWriteRepository {

	return &TicketCategoryWriteRepository{}
}

func (r *TicketCategoryWriteRepository) SyncCategories(
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

	keepCategories := make(
		map[uuid.UUID]bool,
	)

	for _, category := range categories {

		if category.ID == uuid.Nil {

			if err := r.createCategory(
				tx,
				partyID,
				category,
			); err != nil {

				return err
			}

			continue
		}

		keepCategories[category.ID] = true

		if err := r.updateCategory(
			tx,
			category,
		); err != nil {

			return err
		}

		if err := r.syncWindows(
			tx,
			category,
		); err != nil {

			return err
		}
	}

	for _, existingCategory := range existing {

		if keepCategories[existingCategory.ID] {
			continue
		}

		if err := r.deleteCategory(
			tx,
			existingCategory.ID,
		); err != nil {

			return err
		}
	}

	return nil
}

func (r *TicketCategoryWriteRepository) createCategory(
	tx database.DBExecutor,
	partyID uuid.UUID,
	category models.TicketCategory,
) error {

	category.ID = uuid.New()
	category.PartyID = partyID

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

	return nil
}

func (r *TicketCategoryWriteRepository) updateCategory(
	tx database.DBExecutor,
	category models.TicketCategory,
) error {

	return tx.
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
		Error()
}

func (r *TicketCategoryWriteRepository) syncWindows(
	tx database.DBExecutor,
	category models.TicketCategory,
) error {

	var existing []models.TicketAccessWindow

	if err := tx.
		Where(
			"ticket_category_id = ?",
			category.ID,
		).
		Find(
			&existing,
		).
		Error(); err != nil {

		return err
	}

	keep := make(
		map[uuid.UUID]bool,
	)

	for _, window := range category.AccessWindows {

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

		keep[window.ID] = true

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

	for _, old := range existing {

		if keep[old.ID] {
			continue
		}

		if err := tx.
			Unscoped().
			Delete(
				&models.TicketAccessWindow{},
				"id = ?",
				old.ID,
			).
			Error(); err != nil {

			return err
		}
	}

	return nil
}

func (r *TicketCategoryWriteRepository) deleteCategory(
	tx database.DBExecutor,
	categoryID uuid.UUID,
) error {

	if err := tx.
		Unscoped().
		Where(
			"ticket_category_id = ?",
			categoryID,
		).
		Delete(
			&models.TicketAccessWindow{},
		).
		Error(); err != nil {

		return err
	}

	return tx.
		Unscoped().
		Delete(
			&models.TicketCategory{},
			"id = ?",
			categoryID,
		).
		Error()
}
