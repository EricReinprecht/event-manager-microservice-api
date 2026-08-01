package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/handlers"
	"github.com/reinp/event-platform/backend/internal/routes/constants"
)

func registerCategoryRoutes(
	protected *gin.RouterGroup,
	deps *dependencies.Container,
) {

	handler := handlers.NewCategoryHandler(
		deps.CategoryService,
	)

	protected.GET(
		constants.CategoryList,
		handler.GetAll,
	)

	protected.GET(
		constants.CategoryListPopular,
		handler.GetPaginatedByPopularity,
	)

	protected.POST(
		constants.CategoryCreate,
		handler.Create,
	)

	protected.GET(
		constants.CategoryByID,
		handler.GetByID,
	)

	protected.PUT(
		constants.CategoryByID,
		handler.Update,
	)

	protected.DELETE(
		constants.CategoryByID,
		handler.Delete,
	)

}
