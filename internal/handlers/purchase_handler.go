package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/helpers"
	"github.com/reinp/event-platform/backend/internal/responses"
	"github.com/reinp/event-platform/backend/internal/service"
)

type PurchaseHandler struct {
	purchaseService *service.PurchaseService
}

func NewPurchaseHandler(
	purchaseService *service.PurchaseService,
) *PurchaseHandler {

	return &PurchaseHandler{
		purchaseService: purchaseService,
	}
}

func (h *PurchaseHandler) Create(
	c *gin.Context,
) {

	partyID, err := helpers.UUIDParam(
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

	userID, ok := helpers.RequireUserID(
		c,
	)

	if !ok {

		responses.Unauthorized(
			c,
		)

		return
	}

	var req dto.CreatePurchaseRequest

	if err := c.ShouldBindJSON(
		&req,
	); err != nil {

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

	response, err :=
		h.purchaseService.CreatePurchase(
			c.Request.Context(),
			userID,
			partyID,
			req.Items,
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

func (h *PurchaseHandler) GetByID(
	c *gin.Context,
) {

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

	purchase, err :=
		h.purchaseService.GetPurchase(
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
		purchase,
	)
}
