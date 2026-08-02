package models

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	Filename string `json:"filename"`

	Path string `json:"path"`

	URL string `json:"url"`

	MimeType string `json:"mimeType"`

	Size int64 `json:"size"`

	Width int

	Height int

	CreatedAt time.Time
}
