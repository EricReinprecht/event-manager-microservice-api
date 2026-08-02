package ticket_category_repository

import (
	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type TicketCategoryWriteRepository struct{}

func NewTicketCategoryWriteRepository() *TicketCategoryWriteRepository {

	return &TicketCategoryWriteRepository{}
}

func (r *TicketCategoryWriteRepository) SyncCategories(
	tx database.DBExecutor,
	partyID uuid.UUID,
	categories []models.TicketCategory,
) error {

	var existingCategories []models.TicketCategory

	if err := tx.
		Where(
			"party_id = ?",
			partyID,
		).
		Find(
			&existingCategories,
		).
		Error(); err != nil {

		return err
	}

	existingCategoryIDs := make(
		map[uuid.UUID]struct{},
		len(existingCategories),
	)

	for _, category := range existingCategories {

		existingCategoryIDs[category.ID] =
			struct{}{}
	}

	keepCategoryIDs := make(
		map[uuid.UUID]struct{},
		len(categories),
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

		if _, exists :=
			existingCategoryIDs[category.ID]; !exists {

			return appErrors.ErrTicketCategoryNotFound
		}

		keepCategoryIDs[category.ID] =
			struct{}{}

		if err := r.updateCategory(
			tx,
			partyID,
			category,
		); err != nil {

			return err
		}

		if err := r.syncAccessWindows(
			tx,
			category,
		); err != nil {

			return err
		}
	}

	for _, existingCategory := range existingCategories {

		if _, keep :=
			keepCategoryIDs[existingCategory.ID]; keep {

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

	accessWindows :=
		category.AccessWindows

	category.AccessWindows = nil

	if err := tx.
		Create(
			&category,
		).
		Error(); err != nil {

		return err
	}

	for _, accessWindow := range accessWindows {

		accessWindow.ID = uuid.New()

		accessWindow.TicketCategoryID =
			category.ID

		if err := tx.
			Create(
				&accessWindow,
			).
			Error(); err != nil {

			return err
		}
	}

	return nil
}

func (r *TicketCategoryWriteRepository) updateCategory(
	tx database.DBExecutor,
	partyID uuid.UUID,
	category models.TicketCategory,
) error {

	result := tx.
		Model(
			&models.TicketCategory{},
		).
		Where(
			"id = ? AND party_id = ?",
			category.ID,
			partyID,
		).
		Updates(
			map[string]any{
				"name": category.Name,

				"price": category.Price,

				"capacity": category.Capacity,

				"requires_verification": category.RequiresVerification,

				"refund_requires_approval": category.RefundRequiresApproval,

				"refund_policy_id": category.RefundPolicyID,
			},
		)

	if err := result.Error(); err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return appErrors.ErrTicketCategoryNotFound
	}

	return nil
}

func (r *TicketCategoryWriteRepository) syncAccessWindows(
	tx database.DBExecutor,
	category models.TicketCategory,
) error {

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

	existingWindowIDs := make(
		map[uuid.UUID]struct{},
		len(existingWindows),
	)

	for _, window := range existingWindows {

		existingWindowIDs[window.ID] =
			struct{}{}
	}

	keepWindowIDs := make(
		map[uuid.UUID]struct{},
		len(category.AccessWindows),
	)

	for _, window := range category.AccessWindows {

		if window.ID == uuid.Nil {

			window.ID = uuid.New()

			window.TicketCategoryID =
				category.ID

			if err := tx.
				Create(
					&window,
				).
				Error(); err != nil {

				return err
			}

			continue
		}

		if _, exists :=
			existingWindowIDs[window.ID]; !exists {

			return appErrors.ErrInvalidAccessWindow
		}

		keepWindowIDs[window.ID] =
			struct{}{}

		result := tx.
			Model(
				&models.TicketAccessWindow{},
			).
			Where(
				"id = ? AND ticket_category_id = ?",
				window.ID,
				category.ID,
			).
			Updates(
				map[string]any{
					"starts_at": window.StartsAt,

					"ends_at": window.EndsAt,
				},
			)

		if err := result.Error(); err != nil {
			return err
		}

		if result.RowsAffected() == 0 {
			return appErrors.ErrInvalidAccessWindow
		}
	}

	for _, existingWindow := range existingWindows {

		if _, keep :=
			keepWindowIDs[existingWindow.ID]; keep {

			continue
		}

		if err := tx.
			Unscoped().
			Delete(
				&models.TicketAccessWindow{},
				"id = ? AND ticket_category_id = ?",
				existingWindow.ID,
				category.ID,
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
