package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/models"
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

func getUserID(
	c *gin.Context,
) (uuid.UUID, bool) {

	value, exists := c.Get("userID")

	if !exists {
		return uuid.Nil, false
	}

	switch id := value.(type) {

	case uuid.UUID:
		return id, true

	case string:

		userID, err := uuid.Parse(id)

		if err != nil {
			return uuid.Nil, false
		}

		return userID, true
	}

	return uuid.Nil, false
}

func (h *PartyHandler) Create(
	c *gin.Context,
) {

	var req dto.CreatePartyRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		responses.Error(
			c,
			http.StatusBadRequest,
			err,
		)

		return
	}

	if req.EndAt.Before(req.StartAt) {

		responses.Error(
			c,
			http.StatusBadRequest,
			errors.New(
				"end date must be after start date",
			),
		)

		return
	}

	if req.Latitude == 0 || req.Longitude == 0 {

		responses.Error(
			c,
			http.StatusBadRequest,
			errors.New(
				"location coordinates are required",
			),
		)

		return
	}

	if req.Timezone == "" {

		responses.Error(
			c,
			http.StatusBadRequest,
			errors.New(
				"timezone is required",
			),
		)

		return
	}

	userID, ok := getUserID(c)

	if !ok {

		responses.Error(
			c,
			http.StatusUnauthorized,
			errors.New(
				"not authenticated",
			),
		)

		return
	}

	party := &models.Party{

		Title: req.Title,

		Description: req.Description,

		LocationName: req.LocationName,

		Latitude: req.Latitude,

		Longitude: req.Longitude,

		Timezone: req.Timezone,

		StartAt: req.StartAt,

		EndAt: req.EndAt,

		OrganizerID: userID,

		ThumbnailID: req.ThumbnailID,

		Categories: []models.Category{},
	}

	err := h.service.Create(
		c.Request.Context(),
		party,
		req.ImageIDs,
	)

	if err != nil {

		switch {

		case errors.Is(
			err,
			appErrors.ErrCategoryNotFound,
		):

			responses.Error(
				c,
				http.StatusBadRequest,
				err,
			)

			return

		case errors.Is(
			err,
			appErrors.ErrMediaNotFound,
		):

			responses.Error(
				c,
				http.StatusBadRequest,
				err,
			)

			return

		default:

			responses.Error(
				c,
				http.StatusInternalServerError,
				err,
			)

			return
		}
	}

	party, err = h.service.FindByID(
		c.Request.Context(),
		party.ID,
	)

	if err != nil {

		responses.Error(
			c,
			http.StatusInternalServerError,
			err,
		)

		return
	}

	responses.Success(
		c,
		http.StatusCreated,
		party,
	)
}

func (h *PartyHandler) GetAll(
	c *gin.Context,
) {

	parties, err := h.service.FindAll(
		c.Request.Context(),
	)

	if err != nil {

		responses.Error(
			c,
			http.StatusInternalServerError,
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

func (h *PartyHandler) GetByID(
	c *gin.Context,
) {

	id, err := uuid.Parse(
		c.Param("id"),
	)

	if err != nil {

		responses.Error(
			c,
			http.StatusBadRequest,
			errors.New("invalid id"),
		)

		return
	}

	party, err := h.service.FindByID(
		c.Request.Context(),
		id,
	)

	if err != nil {

		responses.Error(
			c,
			http.StatusNotFound,
			errors.New("party not found"),
		)

		return
	}

	responses.Success(
		c,
		http.StatusOK,
		party,
	)
}

func (h *PartyHandler) Update(
	c *gin.Context,
) {

	id, err := uuid.Parse(
		c.Param("id"),
	)

	if err != nil {

		responses.Error(
			c,
			http.StatusBadRequest,
			errors.New("invalid id"),
		)

		return
	}

	userID, ok := getUserID(c)

	if !ok {

		responses.Error(
			c,
			http.StatusUnauthorized,
			errors.New("invalid user"),
		)

		return
	}

	party, err := h.service.FindOwnedParty(
		c.Request.Context(),
		id,
		userID,
	)

	if err != nil {

		responses.Error(
			c,
			http.StatusForbidden,
			errors.New("not allowed"),
		)

		return
	}

	var req dto.UpdatePartyRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		responses.Error(
			c,
			http.StatusBadRequest,
			err,
		)

		return
	}

	if req.EndAt.Before(req.StartAt) {

		responses.Error(
			c,
			http.StatusBadRequest,
			errors.New(
				"end date must be after start date",
			),
		)

		return
	}

	if req.Latitude == 0 || req.Longitude == 0 {

		responses.Error(
			c,
			http.StatusBadRequest,
			errors.New(
				"location coordinates are required",
			),
		)

		return
	}

	if req.Timezone == "" {

		responses.Error(
			c,
			http.StatusBadRequest,
			errors.New(
				"timezone is required",
			),
		)

		return
	}

	party.Title = req.Title

	party.Description = req.Description

	party.LocationName = req.LocationName

	party.Latitude = req.Latitude

	party.Longitude = req.Longitude

	party.Timezone = req.Timezone

	party.ThumbnailID = req.ThumbnailID

	party.StartAt = req.StartAt

	party.EndAt = req.EndAt

	if err := h.service.Update(
		c.Request.Context(),
		party,
	); err != nil {

		responses.Error(
			c,
			http.StatusInternalServerError,
			err,
		)

		return
	}

	if err := h.service.UpdateImages(
		c.Request.Context(),
		party.ID,
		req.ImageIDs,
	); err != nil {

		responses.Error(
			c,
			http.StatusInternalServerError,
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

func (h *PartyHandler) Delete(
	c *gin.Context,
) {

	id, err := uuid.Parse(
		c.Param("id"),
	)

	if err != nil {

		responses.Error(
			c,
			http.StatusBadRequest,
			errors.New("invalid id"),
		)

		return
	}

	userID, ok := getUserID(c)

	if !ok {

		responses.Error(
			c,
			http.StatusUnauthorized,
			errors.New("invalid user"),
		)

		return
	}

	party, err := h.service.FindOwnedParty(
		c.Request.Context(),
		id,
		userID,
	)

	if err != nil {

		responses.Error(
			c,
			http.StatusForbidden,
			errors.New("not allowed"),
		)

		return
	}

	if err := h.service.Delete(
		c.Request.Context(),
		party,
	); err != nil {

		responses.Error(
			c,
			http.StatusInternalServerError,
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

func (h *PartyHandler) GetMyParties(
	c *gin.Context,
) {

	userID, ok := getUserID(c)

	if !ok {

		responses.Error(
			c,
			http.StatusUnauthorized,
			errors.New("not authenticated"),
		)

		return
	}

	name := c.Query("name")
	startAt := c.Query("startAt")
	endAt := c.Query("endAt")

	page, _ := strconv.Atoi(
		c.DefaultQuery(
			"page",
			"1",
		),
	)

	limit, _ := strconv.Atoi(
		c.DefaultQuery(
			"limit",
			"10",
		),
	)

	sorts := c.Query("sorts")

	parties, total, err := h.service.FindOrganizedByUser(
		c.Request.Context(),
		userID,
		name,
		startAt,
		endAt,
		sorts,
		page,
		limit,
	)

	if err != nil {

		responses.Error(
			c,
			http.StatusInternalServerError,
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
