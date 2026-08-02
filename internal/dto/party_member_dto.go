package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models/enum"
)

type CreatePartyMemberRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`

	Roles []enum.PartyMemberRole `json:"roles" binding:"required"`
}

type UpdatePartyMemberRolesRequest struct {
	Roles []enum.PartyMemberRole `json:"roles" binding:"required"`
}

type PartyMemberResponse struct {
	ID uuid.UUID `json:"id"`

	UserID uuid.UUID `json:"userId"`

	PartyID uuid.UUID `json:"partyId"`

	Roles []PartyMemberRoleResponse `json:"roles"`

	CreatedAt time.Time `json:"createdAt"`
}

type PartyMemberRoleResponse struct {
	ID uuid.UUID `json:"id"`

	Role enum.PartyMemberRole `json:"role"`
}
