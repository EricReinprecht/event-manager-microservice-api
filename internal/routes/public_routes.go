package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/handlers"
	"github.com/reinp/event-platform/backend/internal/middleware"
	"github.com/reinp/event-platform/backend/internal/routes/constants"
)

func registerPublicRoutes(
	router *gin.Engine,
	api *gin.RouterGroup,
	deps *dependencies.Container,
) {

	authHandler := handlers.NewAuthHandler(
		deps.AuthService,
		deps.RefreshTokenDuration,
		deps.CookieSecure,
	)

	paymentHandler := handlers.NewPaymentHandler(
		deps.PaymentService,
	)

	authLimiter := middleware.NewIPRateLimiter(
		1.0/12,
		5,
	)

	refreshLimiter := middleware.NewIPRateLimiter(
		0.5,
		20,
	)

	router.GET(
		"/health",
		handlers.Health,
	)

	api.POST(
		constants.AuthRegister,
		authLimiter.Middleware(),
		authHandler.Register,
	)

	api.GET(
		constants.AuthVerifyEmail,
		authLimiter.Middleware(),
		authHandler.VerifyEmail,
	)

	api.POST(
		constants.AuthLogin,
		authLimiter.Middleware(),
		authHandler.Login,
	)

	api.POST(
		constants.AuthRefresh,
		refreshLimiter.Middleware(),
		authHandler.Refresh,
	)

	api.POST(
		constants.AuthLogout,
		authLimiter.Middleware(),
		authHandler.Logout,
	)

	api.POST(
		constants.AuthForgotPassword,
		refreshLimiter.Middleware(),
		authHandler.ForgotPassword,
	)

	api.POST(
		constants.AuthResetPassword,
		authLimiter.Middleware(),
		authHandler.ResetPassword,
	)

	api.POST(
		constants.AuthResendVerification,
		authLimiter.Middleware(),
		authHandler.ResendVerificationEmail,
	)

	router.POST(
		constants.PayPalWebhook,
		paymentHandler.Webhook,
	)

}
