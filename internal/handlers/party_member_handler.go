package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/helpers"
	"github.com/reinp/event-platform/backend/internal/mapper"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/responses"
	"github.com/reinp/event-platform/backend/internal/service"
)

type PartyMemberHandler struct {
	partyMemberService *service.PartyMemberService
}

func NewPartyMemberHandler(
	service *service.PartyMemberService,
) *PartyMemberHandler {

	return &PartyMemberHandler{
		partyMemberService: service,
	}
}

func (h *PartyMemberHandler) Create(
	c *gin.Context,
) {

	ctx := c.Request.Context()

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

	userID, ok := helpers.RequireUserID(c)

	if !ok {

		responses.Unauthorized(c)

		return
	}

	if !h.partyMemberService.HasRole(
		ctx,
		partyID,
		userID,
		enum.RoleOrganizer,
		enum.RoleAdmin,
	) {

		responses.Forbidden(c)

		return
	}

	var req dto.CreatePartyMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	member, err := h.partyMemberService.Create(
		ctx,
		partyID,
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
		mapper.PartyMemberResponse(*member),
	)
}

func (h *PartyMemberHandler) GetAll(
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

	members, err := h.partyMemberService.FindByParty(
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
		mapper.PartyMemberResponses(members),
	)
}

func (h *PartyMemberHandler) Delete(
	c *gin.Context,
) {

	memberID, err := helpers.UUIDParam(
		c,
		"memberID",
	)

	if err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	err = h.partyMemberService.Delete(
		c.Request.Context(),
		memberID,
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
			"message": "member removed",
		},
	)
}

func (h *PartyMemberHandler) UpdateRoles(
	c *gin.Context,
) {

	memberID, err := helpers.UUIDParam(
		c,
		"memberID",
	)

	if err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	var req dto.UpdatePartyMemberRolesRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		responses.BadRequest(
			c,
			err,
		)

		return
	}

	err = h.partyMemberService.SyncRoles(
		c.Request.Context(),
		memberID,
		req.Roles,
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
			"message": "roles updated",
		},
	)
}
