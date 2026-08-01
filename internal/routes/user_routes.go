package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/handlers"
	"github.com/reinp/event-platform/backend/internal/routes/constants"
)

func registerUserRoutes(
	protected *gin.RouterGroup,
	deps *dependencies.Container,
) {

	userHandler := handlers.NewUserHandler(
		deps.UserService,
	)

	sessionHandler := handlers.NewSessionHandler(
		deps.AuthService,
	)

	protected.GET(
		constants.UserSessions,
		sessionHandler.GetSessions,
	)

	protected.DELETE(
		constants.UserSessionByFamilyID,
		sessionHandler.DeleteSession,
	)

	protected.DELETE(
		constants.UserSessions,
		sessionHandler.LogoutAll,
	)

	protected.GET(
		constants.UserMe,
		userHandler.Me,
	)

	protected.PUT(
		constants.UserCompleteProfile,
		userHandler.CompleteProfile,
	)

	protected.PUT(
		constants.UserPassword,
		userHandler.ChangePassword,
	)

	protected.GET(
		constants.UserParties,
		handlers.NewPartyHandler(
			deps.PartyService,
		).GetMyParties,
	)

}
