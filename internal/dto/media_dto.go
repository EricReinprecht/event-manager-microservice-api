package dto

import "github.com/google/uuid"

type MediaResponse struct {
	ID       uuid.UUID `json:"id"`
	Filename string    `json:"filename"`
	URL      string    `json:"url"`
	MimeType string    `json:"mimeType"`
}
