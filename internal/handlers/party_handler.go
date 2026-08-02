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
	partyService *service.PartyService
}

func NewPartyHandler(
	partyService *service.PartyService,
) *PartyHandler {

	return &PartyHandler{
		partyService: partyService,
	}
}

func (h *PartyHandler) Create(c *gin.Context) {

	var req dto.CreatePartyRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		validationError :=
			helpers.BindingValidationErrors(
				err,
				req,
			)

		if validationError != nil {

			responses.HandleDomainError(
				c,
				validationError,
			)

			return
		}

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

	response, err := h.partyService.Create(
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

	parties, err := h.partyService.FindAll(
		c.Request.Context(),
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

	party, err := h.partyService.FindByID(
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

		validationError :=
			helpers.BindingValidationErrors(
				err,
				req,
			)

		if validationError != nil {

			responses.HandleDomainError(
				c,
				validationError,
			)

			return
		}

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

	response, err := h.partyService.Update(
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

	err = h.partyService.Delete(
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

func (h *PartyHandler) Publish(c *gin.Context) {
	id, err := helpers.UUIDParam(c, "id")
	if err != nil {
		responses.BadRequest(c, err)
		return
	}

	userID, ok := helpers.RequireUserID(c)
	if !ok {
		responses.Unauthorized(c)
		return
	}

	response, err := h.partyService.Publish(c.Request.Context(), id, userID)
	if err != nil {
		responses.HandleDomainError(c, err)
		return
	}

	responses.Success(c, http.StatusOK, response)
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

	parties, total, err := h.partyService.FindOrganizedByUser(
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

		responses.HandleDomainError(
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
