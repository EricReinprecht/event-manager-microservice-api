package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/models"
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

		c.JSON(400, gin.H{
			"error": "file missing",
		})

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

		c.JSON(500, gin.H{
			"error": err.Error(),
		})

		return
	}

	media := &models.Media{

		Filename: filename,

		Path: path,

		URL: "/" + path,

		MimeType: file.Header.Get("Content-Type"),

		Size: file.Size,
	}

	err = h.service.Create(
		c.Request.Context(),
		media,
	)

	if err != nil {

		c.JSON(500, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(
		http.StatusCreated,
		media,
	)
}
