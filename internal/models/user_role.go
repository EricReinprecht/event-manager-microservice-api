package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models/enum"
)

type UserRole struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_role"`
	User   User      `gorm:"constraint:OnDelete:CASCADE"`

	Role enum.UserRole `gorm:"type:varchar(30);not null;uniqueIndex:idx_user_role;check:user_role_check,role IN ('DEFAULT','PREMIUM','PLATINUM','ARTIST','VERIFIED_ARTIST','PARTNER','MODERATOR')"`

	CreatedAt time.Time
}
