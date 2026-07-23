package models

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Filename string

	Path string

	URL string

	MimeType string

	Size int64

	Width int

	Height int

	CreatedAt time.Time
}
