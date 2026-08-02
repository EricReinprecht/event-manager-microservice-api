package handlers

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

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
	if file.Size > 10*1024*1024 {
		responses.BadRequest(c, errors.New("file exceeds 10 MB limit"))
		return
	}
	if contentType := file.Header.Get("Content-Type"); len(contentType) < 6 || contentType[:6] != "image/" {
		responses.BadRequest(c, errors.New("only image uploads are supported"))
		return
	}

	path := "uploads/" + uuid.NewString() + filepath.Ext(filename)

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
