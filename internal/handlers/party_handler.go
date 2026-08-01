package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/helpers"
	"github.com/reinp/event-platform/backend/internal/responses"
	"github.com/reinp/event-platform/backend/internal/service"
)

type PartyHandler struct {
	service *service.PartyService
}

func NewPartyHandler(
	service *service.PartyService,
) *PartyHandler {

	return &PartyHandler{
		service: service,
	}
}

func (h *PartyHandler) Create(c *gin.Context) {

	var req dto.CreatePartyRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	userID, ok := helpers.RequireUserID(c)

	if !ok {

		responses.Unauthorized(c)

		return
	}

	response, err := h.service.Create(
		c.Request.Context(),
		req,
		userID,
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
		response,
	)
}

func (h *PartyHandler) GetAll(c *gin.Context) {

	parties, err := h.service.FindAll(
		c.Request.Context(),
	)

	if err != nil {

		responses.InternalError(
			c,
			err,
		)

		return
	}

	responses.Success(
		c,
		http.StatusOK,
		parties,
	)
}

func (h *PartyHandler) GetByID(c *gin.Context) {

	id, err := helpers.UUIDParam(
		c,
		"id",
	)

	if err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	party, err := h.service.FindByID(
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
		party,
	)
}

func (h *PartyHandler) Update(c *gin.Context) {

	ctx := c.Request.Context()

	id, err := helpers.UUIDParam(
		c,
		"id",
	)

	if err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	var req dto.UpdatePartyRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	response, err := h.service.Update(
		ctx,
		id,
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
		response,
	)
}

func (h *PartyHandler) Delete(c *gin.Context) {

	id, err := helpers.UUIDParam(
		c,
		"id",
	)

	if err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	err = h.service.Delete(
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
		gin.H{
			"message": "party deleted",
		},
	)
}

func (h *PartyHandler) GetMyParties(c *gin.Context) {

	userID, ok := helpers.RequireUserID(c)

	if !ok {

		responses.Unauthorized(c)

		return
	}

	name := c.Query("name")
	startAt := c.Query("startAt")
	endAt := c.Query("endAt")
	locationName := c.Query("locationName")

	page, limit := helpers.QueryPagination(c)

	sorts := c.Query("sorts")

	parties, total, err := h.service.FindOrganizedByUser(
		c.Request.Context(),
		userID,
		name,
		startAt,
		endAt,
		locationName,
		sorts,
		page,
		limit,
	)

	if err != nil {

		responses.InternalError(
			c,
			err,
		)

		return
	}

	responses.Paginated(
		c,
		parties,
		page,
		limit,
		total,
	)
}
