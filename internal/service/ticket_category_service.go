package service

import (
	"context"
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/mapper"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	ticketCategoryRepository "github.com/reinp/event-platform/backend/internal/repository/ticket_category"
)

type TicketCategoryService struct {
	ticketCategories *ticketCategoryRepository.Facade
	permission       *PermissionService
}

func NewTicketCategoryService(
	ticketCategories *ticketCategoryRepository.Facade,
	permission *PermissionService,
) *TicketCategoryService {

	return &TicketCategoryService{
		ticketCategories: ticketCategories,
		permission:       permission,
	}
}

func (s *TicketCategoryService) CreateFromRequest(
	ctx context.Context,
	partyID, userID uuid.UUID,
	req dto.CreateTicketCategoryRequest,
) (*models.TicketCategory, error) {
	if err := s.permission.RequirePartyRole(ctx, partyID, userID, enum.PartyRoleOrganizer, enum.PartyRoleAdmin); err != nil {
		return nil, err
	}
	category := &models.TicketCategory{
		Name: req.Name, Price: req.Price, Capacity: req.Capacity,
		RequiresVerification:   req.RequiresVerification,
		RefundRequiresApproval: req.RefundRequiresApproval,
		RefundPolicyID:         req.RefundPolicyID, PartyID: partyID,
		AccessWindows: mapper.AccessWindowsFromRequest(req.AccessWindows),
	}
	if err := s.Create(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *TicketCategoryService) UpdateFromRequest(
	ctx context.Context,
	id, userID uuid.UUID,
	req dto.UpdateTicketCategoryRequest,
) (*models.TicketCategory, error) {
	category, err := s.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.permission.RequirePartyRole(ctx, category.PartyID, userID, enum.PartyRoleOrganizer, enum.PartyRoleAdmin); err != nil {
		return nil, err
	}
	category.Name, category.Price, category.Capacity = req.Name, req.Price, req.Capacity
	category.RequiresVerification = req.RequiresVerification
	category.RefundRequiresApproval, category.RefundPolicyID = req.RefundRequiresApproval, req.RefundPolicyID
	category.AccessWindows = mapper.AccessWindowsFromUpdateRequest(req.AccessWindows)
	if err := s.Update(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *TicketCategoryService) DeleteByID(ctx context.Context, id, userID uuid.UUID) error {
	category, err := s.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.permission.RequirePartyRole(ctx, category.PartyID, userID, enum.PartyRoleOrganizer, enum.PartyRoleAdmin); err != nil {
		return err
	}
	return s.Delete(ctx, category)
}

func (s *TicketCategoryService) Create(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	if err := s.validateCategory(category); err != nil {
		return err
	}

	return s.ticketCategories.Repository.Create(ctx, category)
}

func (s *TicketCategoryService) Update(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	if err := s.validateCategory(category); err != nil {
		return err
	}

	return s.ticketCategories.Repository.Update(ctx, category)
}

func (s *TicketCategoryService) FindByParty(
	ctx context.Context,
	partyID uuid.UUID,
) ([]models.TicketCategory, error) {

	return s.ticketCategories.Repository.FindByParty(
		ctx,
		partyID,
	)
}

func (s *TicketCategoryService) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.TicketCategory, error) {

	return s.ticketCategories.Repository.FindByID(
		ctx,
		id,
	)
}

func (s *TicketCategoryService) Delete(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	return s.ticketCategories.Repository.Delete(
		ctx,
		category,
	)
}

func (s *TicketCategoryService) validateCategory(
	category *models.TicketCategory,
) error {

	if len(category.AccessWindows) == 0 {
		return appErrors.ErrTicketAccessWindowRequired
	}

	for _, window := range category.AccessWindows {

		if window.EndsAt.Before(window.StartsAt) {
			return appErrors.ErrAccessWindowInvalid
		}
	}

	return nil
}
