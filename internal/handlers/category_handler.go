package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/helpers"
	"github.com/reinp/event-platform/backend/internal/responses"
	"github.com/reinp/event-platform/backend/internal/service"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(
	service *service.CategoryService,
) *CategoryHandler {

	return &CategoryHandler{
		service: service,
	}
}

func (h *CategoryHandler) GetAll(c *gin.Context) {

	categories, err := h.service.FindAll(
		c.Request.Context(),
	)

	if err != nil {
		responses.HandleDomainError(c, err)
		return
	}

	responses.Success(c, http.StatusOK, categories)
}

func (h *CategoryHandler) Create(c *gin.Context) {

	var req dto.CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		responses.BadRequest(c, err)
		return
	}

	category, err := h.service.Create(
		c.Request.Context(),
		req,
	)

	if err != nil {

		responses.HandleDomainError(c, err)
		return
	}

	responses.Success(c, http.StatusCreated, category)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {

	id, err := helpers.UUIDParam(c, "id")

	if err != nil {
		responses.BadRequest(c, err)
		return
	}

	category, err := h.service.FindByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		responses.HandleDomainError(c, err)
		return
	}

	responses.Success(c, http.StatusOK, category)
}

func (h *CategoryHandler) Update(c *gin.Context) {

	id, err := helpers.UUIDParam(c, "id")

	if err != nil {
		responses.BadRequest(c, err)
		return
	}

	var req dto.UpdateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		responses.BadRequest(c, err)
		return
	}

	category, err := h.service.Update(
		c.Request.Context(),
		id,
		req,
	)

	if err != nil {

		responses.HandleDomainError(c, err)
		return
	}

	responses.Success(c, http.StatusOK, category)
}

func (h *CategoryHandler) Delete(c *gin.Context) {

	id, err := helpers.UUIDParam(c, "id")

	if err != nil {
		responses.BadRequest(c, err)
		return
	}

	err = h.service.Delete(
		c.Request.Context(),
		id,
	)

	if err != nil {

		responses.HandleDomainError(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "category deleted",
		},
	)
}

func (h *CategoryHandler) GetPaginatedByPopularity(
	c *gin.Context,
) {

	limit, _ := strconv.Atoi(
		c.DefaultQuery(
			"limit",
			"10",
		),
	)

	categories, err := h.service.FindPaginatedByPopularity(
		c.Request.Context(),
		limit,
	)

	if err != nil {

		responses.Error(
			c,
			http.StatusInternalServerError,
			err,
		)

		return
	}

	responses.Success(
		c,
		http.StatusOK,
		categories,
	)
}
