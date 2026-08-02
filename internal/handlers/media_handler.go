package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/responses"
	"github.com/reinp/event-platform/backend/internal/service"
)

type MediaHandler struct {
	service *service.MediaService
}

func NewMediaHandler(
	service *service.MediaService,
) *MediaHandler {

	return &MediaHandler{
		service: service,
	}
}

func (h *MediaHandler) Upload(
	c *gin.Context,
) {

	file, err := c.FormFile("file")

	if err != nil {

		responses.BadRequest(c, err)

		return
	}

	filename := filepath.Base(
		file.Filename,
	)

	path := "uploads/" + filename

	err = c.SaveUploadedFile(
		file,
		path,
	)

	if err != nil {

		responses.HandleDomainError(c, err)

		return
	}

	media, err := h.service.Create(
		c.Request.Context(),
		filename,
		path,
		file.Header.Get("Content-Type"),
		file.Size,
	)

	if err != nil {

		responses.HandleDomainError(c, err)

		return
	}

	responses.Success(c, http.StatusCreated, media)
}
