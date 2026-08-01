package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/helpers"
	"github.com/reinp/event-platform/backend/internal/mapper"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/responses"
	"github.com/reinp/event-platform/backend/internal/service"
)

type TicketCategoryHandler struct {
	service    *service.TicketCategoryService
	permission *service.PermissionService
}

func NewTicketCategoryHandler(
	service *service.TicketCategoryService,
	permission *service.PermissionService,
) *TicketCategoryHandler {

	return &TicketCategoryHandler{
		service:    service,
		permission: permission,
	}
}

func (h *TicketCategoryHandler) Create(c *gin.Context) {

	ctx := c.Request.Context()

	partyID, err := helpers.UUIDParam(c, "id")

	if err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	userID, ok := helpers.MustUserID(c)

	if !ok {
		return
	}

	if err := h.permission.RequirePartyRole(
		ctx,
		partyID,
		userID,
		enum.RoleOrganizer,
		enum.RoleAdmin,
	); err != nil {

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}

	var req dto.CreateTicketCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	category := &models.TicketCategory{
		Name:                 req.Name,
		Price:                req.Price,
		Capacity:             req.Capacity,
		RequiresVerification: req.RequiresVerification,
		PartyID:              partyID,
		AccessWindows: mapper.AccessWindowsFromRequest(
			req.AccessWindows,
		),
	}

	if err := h.service.Create(
		ctx,
		category,
	); err != nil {

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}

	responses.Success(
		c,
		http.StatusCreated,
		mapper.TicketCategoryResponse(category),
	)
}

func (h *TicketCategoryHandler) GetAll(c *gin.Context) {

	partyID, err := helpers.UUIDParam(c, "id")

	if err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	categories, err := h.service.FindByParty(
		c.Request.Context(),
		partyID,
	)

	if err != nil {

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}

	responses.Success(
		c,
		http.StatusOK,
		mapper.TicketCategoryResponses(categories),
	)
}

func (h *TicketCategoryHandler) GetByID(c *gin.Context) {

	id, err := helpers.UUIDParam(c, "id")

	if err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	category, err := h.service.FindByID(
		c.Request.Context(),
		id,
	)

	if err != nil {

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}

	responses.Success(
		c,
		http.StatusOK,
		mapper.TicketCategoryResponse(category),
	)
}

func (h *TicketCategoryHandler) Update(c *gin.Context) {

	ctx := c.Request.Context()

	id, err := helpers.UUIDParam(c, "id")

	if err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	category, err := h.service.FindByID(
		ctx,
		id,
	)

	if err != nil {

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}

	userID, ok := helpers.MustUserID(c)

	if !ok {
		return
	}

	if err := h.permission.RequirePartyRole(
		ctx,
		category.PartyID,
		userID,
		enum.RoleOrganizer,
		enum.RoleAdmin,
	); err != nil {

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}

	var req dto.UpdateTicketCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	category.Name = req.Name

	category.Price = req.Price

	category.Capacity = req.Capacity

	category.RequiresVerification = req.RequiresVerification

	category.AccessWindows =
		mapper.AccessWindowsFromUpdateRequest(
			req.AccessWindows,
		)

	if err := h.service.Update(
		ctx,
		category,
	); err != nil {

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}

	responses.Success(
		c,
		http.StatusOK,
		mapper.TicketCategoryResponse(category),
	)
}

func (h *TicketCategoryHandler) Delete(c *gin.Context) {

	ctx := c.Request.Context()

	id, err := helpers.UUIDParam(c, "id")

	if err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	category, err := h.service.FindByID(
		ctx,
		id,
	)

	if err != nil {

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}

	userID, ok := helpers.MustUserID(c)

	if !ok {
		return
	}

	if err := h.permission.RequirePartyRole(
		ctx,
		category.PartyID,
		userID,
		enum.RoleOrganizer,
		enum.RoleAdmin,
	); err != nil {

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}
	if err := h.service.Delete(
		ctx,
		category,
	); err != nil {

		responses.HandleDomainError(
			c,
			err,
		)

		return
	}

	responses.Success(
		c,
		http.StatusOK,
		gin.H{
			"message": "ticket category deleted",
		},
	)
}
