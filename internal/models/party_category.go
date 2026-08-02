package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PartyCategory struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	Name string `gorm:"unique;not null" json:"name"`

	Parties []Party `gorm:"many2many:party_categories;" json:"-"`

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
