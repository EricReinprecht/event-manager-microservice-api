package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/helpers"
	"github.com/reinp/event-platform/backend/internal/mapper"
	"github.com/reinp/event-platform/backend/internal/responses"
	"github.com/reinp/event-platform/backend/internal/service"
)

type TicketCategoryHandler struct {
	service *service.TicketCategoryService
}

func NewTicketCategoryHandler(
	service *service.TicketCategoryService,
) *TicketCategoryHandler {

	return &TicketCategoryHandler{
		service: service,
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

	var req dto.CreateTicketCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	category, err := h.service.CreateFromRequest(
		ctx,
		partyID,
		userID,
		req,
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

	userID, ok := helpers.MustUserID(c)

	if !ok {
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

	category, err := h.service.UpdateFromRequest(
		ctx,
		id,
		userID,
		req,
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

	userID, ok := helpers.MustUserID(c)

	if !ok {
		return
	}

	if err := h.service.DeleteByID(
		ctx,
		id,
		userID,
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
